package rc

import (
	"encoding/json"
	"io/ioutil"

	"go.uber.org/zap"
)

type Client struct {
	url string

	connected bool
	anon      bool

	debug    bool
	realtime bool
	c        *restClient
	d        *ddpClient

	cred *Credential
	log  *zap.SugaredLogger
}

type ClientOption func(*Client)

func Debug(d bool) ClientOption {
	return func(c *Client) {
		c.debug = d
	}
}

func Realtime(b bool) ClientOption {
	return func(c *Client) {
		c.realtime = b
	}
}

func ServerURL(u string) ClientOption {
	return func(c *Client) {
		c.url = u
	}
}

func Credentials(user, pass string) ClientOption {
	return func(c *Client) {
		c.cred.Username = user
		c.cred.Password = pass
	}
}

func AccessToken(id, token string) ClientOption {
	return func(c *Client) {
		c.cred.ID = id
		c.cred.Token = token
	}
}

func CredFromJson(path string) ClientOption {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		panic("unable to read credential file: " + err.Error())
	}

	cred := &Credential{}
	err = json.Unmarshal(f, cred)
	if err != nil {
		panic("unable to read credential file: " + err.Error())
	}

	return func(c *Client) {
		c.cred = cred
	}
}

func Anonymous(b bool) ClientOption {
	return func(c *Client) {
		c.anon = true
	}
}

func New(options ...ClientOption) *Client {
	c := &Client{
		url:  "http://localhost:3000",
		cred: &Credential{},
	}

	c.cred.fromEnv()

	for _, opt := range options {
		opt(c)
	}

	logger := buildLogger(c.debug)
	c.log = logger.Sugar()

	c.c = newRESTClient(c.url, c.debug, c.log.Named("rest"))

	return c
}

func (c *Client) Connect() error {
	if c.connected {
		return nil
	}

	if c.realtime {
		ddp, err := newDDPClient(c.url, c.debug, c.log.Named("ddp"))
		if err != nil {
			return err
		}
		c.d = ddp
	}

	if c.anon {
		c.connected = true
		return nil
	}

	if c.cred.tokenReady() {
		c.c.setAuthHeader(c.cred.ID, c.cred.Token)
		if c.realtime {
			if err := c.ResumeRT(); err != nil {
				return err
			}
		}

		c.connected = true
		return nil
	}

	if c.cred.hasUP() {
		if err := c.Login(); err != nil {
			return err
		}
	}

	return nil
}

func buildLogger(debug bool) *zap.Logger {
	var z *zap.Logger
	var err error
	if debug {
		z, err = zap.NewDevelopment()
		if err != nil {
			panic("unable to initalize logger: " + err.Error())
		}
	} else {
		z, err = zap.NewProduction()
		if err != nil {
			panic("unable to initalize logger: " + err.Error())
		}
	}

	return z
}
