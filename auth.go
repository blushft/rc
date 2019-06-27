package rc

import (
	"os"
	"reflect"
	"time"
)

const (
	rcUserEnv  = "RC_USERNAME"
	rcPassEnv  = "RC_PASSWORD"
	rcTokenEnv = "RC_TOKEN"
	rcIDEnv    = "RC_USERID"
)

type Credential struct {
	Username string    `json:"username,omitempty"`
	Password string    `json:"password,omitempty"`
	Email    string    `json:"email,omitempty"`
	Token    string    `json:"token,omitempty"`
	ID       string    `json:"id,omitempty"`
	Exp      time.Time `json:"exp,omitempty"`
}

func (cred *Credential) isEmpty() bool {
	return reflect.DeepEqual(&Credential{}, cred)
}

func (cred *Credential) tokenReady() bool {
	return cred.Token != "" && cred.ID != ""
}

func (cred *Credential) hasUP() bool {
	return cred.Username != "" && cred.Password != ""
}

func (cred *Credential) fromEnv() {
	cred.Username = os.Getenv(rcUserEnv)
	cred.Password = os.Getenv(rcPassEnv)
	cred.Token = os.Getenv(rcTokenEnv)
	cred.ID = os.Getenv(rcIDEnv)
}

func NewCredential(user, pass, email string) *Credential {
	return &Credential{
		Username: user,
		Password: pass,
		Email:    email,
	}
}

type StandardLogin struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status string    `json:"status,omitempty"`
	Data   LoginData `json:"data,omitempty"`
}

type LoginData struct {
	Token  string `json:"authToken,omitempty"`
	UserID string `json:"userID,omitempty"`
}

type ResumeLogin struct {
	Token string `json:"resume"`
}

func (c *Client) Login() error {
	resp := &LoginResponse{}
	if err := c.c.postJSON("/login", StandardLogin{c.cred.Username, c.cred.Password}).JSON(&resp); err != nil {
		return err
	}
	data := resp.Data

	c.cred.ID = data.UserID
	c.cred.Token = data.Token
	c.c.setAuthHeader(data.UserID, data.Token)

	if c.realtime {
		if err := c.ResumeRT(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ResumeRT() error {
	ld := ResumeLogin{
		Token: c.cred.Token,
	}
	if _, err := c.d.call("login", ld); err != nil {
		return err
	}

	return nil
}