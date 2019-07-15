package rc

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
)

// Client struct contains all methods for interacting with the Rocket.Chat api.
// Client is opinionated about realtime vs syncronous api calls and will use the
// appropriate method (rest vs ddp)
type Client struct {
	url string

	connected bool
	anon      bool

	debug    bool
	realtime bool
	c        *restClient
	d        *ddpClient

	strOpts []StreamOption

	cred *Credential
	log  *zap.SugaredLogger
}

// ClientOption is a functional argument that sets optional values on Client
type ClientOption func(*Client)

// Debug sets Client in debug logging mode
func Debug(d bool) ClientOption {
	return func(c *Client) {
		c.debug = d
	}
}

// Realtime specifies Client should use ddp for all interaction
func StreamOptions(opts ...StreamOption) ClientOption {
	return func(c *Client) {
		c.realtime = true
		c.strOpts = opts
	}
}

// ServerURL sets the server URL for Client
// If unset, Client will use the value in $RC_SERVER_URL and fall back to
// http://localhost:3000
func ServerURL(u string) ClientOption {
	return func(c *Client) {
		c.url = u
	}
}

// Credentials sets Client username and password.
// If unset, Client will use values from $RC_USERNAME and $RC_PASSWORD.
// Client will prefer a token over credentials if available.
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
		url:  getURL(),
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
		ddp, err := newDDPClient(c.url, c.debug, c.log.Named("ddp"), c.strOpts...)
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
			if err := c.Resume(); err != nil {
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

func getURL() string {
	u := os.Getenv(rcURL)
	if u == "" {
		return "http://localhost:3000"
	}
	return u
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
