package rc

import (
	"net/url"

	"github.com/gopackage/ddp"
	"go.uber.org/zap"
)

type ddpClient struct {
	ddp *ddp.Client

	server string
	debug  bool
	log    *zap.SugaredLogger
}

func newDDPClient(server string, debug bool, logger *zap.SugaredLogger) (*ddpClient, error) {
	urlVals, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	if urlVals.Scheme == "https" {
		urlVals.Scheme = "wss"
	} else {
		urlVals.Scheme = "ws"
	}

	urlVals.Path = "websocket"

	u := urlVals.String()

	d := ddp.NewClient(u, server)

	if debug {
		d.SetSocketLogActive(true)
	}

	if err := d.Connect(); err != nil {
		return nil, err
	}

	client := &ddpClient{
		ddp:    d,
		log:    logger,
		server: server,
		debug:  debug,
	}

	return client, nil
}

func (d *ddpClient) Reconnect() {
	d.ddp.Reconnect()
}

func (d *ddpClient) Close() {
	d.ddp.Close()
}

func (d *ddpClient) call(method string, args ...interface{}) (interface{}, error) {
	return d.ddp.Call(method, args...)
}

func (c *Client) NewUpdateListener(fn func(ddp.Update) (interface{}, error)) (ddp.UpdateListener, *SubChannel) {
	u := make(chan interface{}, 10)
	e := make(chan error, 10)

	ml := &updateListener{
		updates: u,
		errors:  e,
		process: fn,
	}

	return ml, &SubChannel{Updates: u, Errors: e}
}

type SubChannel struct {
	Updates <-chan interface{}
	Errors  <-chan error
}

type updateListener struct {
	updates chan interface{}
	errors  chan error
	process func(ddp.Update) (interface{}, error)
}

func (ul updateListener) CollectionUpdate(coll, op, id string, update ddp.Update) {
	u, err := ul.process(update)
	if err != nil {
		ul.errors <- err
	}
	ul.updates <- u
}
