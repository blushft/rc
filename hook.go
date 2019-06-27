package rc

import (
	"gopkg.in/resty.v1"
)

type WebHook struct {
	url string

	c *resty.Client
}

func NewWebHook(url string) *WebHook {
	r := resty.New()
	r.HostURL = url
	h := &WebHook{
		url: url,
		c:   r,
	}

	return h
}

func (h *WebHook) Send(msg Message) error {
	_, err := h.c.R().
		SetBody(msg).
		Post("")

	if err != nil {
		return err
	}
	return nil
}
