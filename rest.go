package rc

import (
	"encoding/json"
	"net/url"
	"strings"

	"go.uber.org/zap"

	"gopkg.in/resty.v1"
)

const (
	RESTAPIPath  = "/api"
	RESTV1Path   = "/v1"
	RESTInfoPath = "/info"

	TimeFormat = "2006-01-02T15:04:05.999Z"
)

type restClient struct {
	*resty.Client
	server string
	rest   string
	info   string

	debug bool
	log   *zap.SugaredLogger
}

func newRESTClient(server string, debug bool, logger *zap.SugaredLogger) *restClient {
	server = stripTrailingSlash(server)

	return &restClient{
		Client: resty.New(),
		server: server,
		rest:   server + RESTAPIPath + RESTV1Path,
		info:   server + RESTAPIPath + RESTInfoPath,
		debug:  debug,
		log:    logger,
	}
}

func stripTrailingSlash(u string) string {
	return strings.TrimRight(u, "/")
}

type Result interface {
	Body() []byte
	Error() error
	JSON(v interface{}) error
	String() string
	StatusCode() int
}

type restReturn struct {
	code int
	body []byte
	err  error
}

func (rr *restReturn) Body() []byte {
	return rr.body
}

func (rr *restReturn) Error() error {
	return rr.err
}

func (rr *restReturn) JSON(v interface{}) error {
	if rr.err != nil {
		return rr.err
	}
	return json.Unmarshal(rr.body, v)
}

func (rr *restReturn) String() string {
	return string(rr.Body())
}

func (rr *restReturn) StatusCode() int {
	return rr.code
}

func (r *restClient) setAuthHeader(id, token string) {
	r.SetHeaders(map[string]string{
		"X-Auth-Token": token,
		"X-User-Id":    id,
	})
}

func (r *restClient) postForm(path string, vals url.Values) Result {
	call, err := r.R().
		SetMultiValueFormData(vals).
		Post(r.rest + path)

	return &restReturn{
		code: call.StatusCode(),
		body: call.Body(),
		err:  err,
	}
}

func (r *restClient) postJSON(path string, v interface{}) Result {
	call, err := r.R().
		SetBody(v).
		Post(r.rest + path)

	return &restReturn{
		code: call.StatusCode(),
		body: call.Body(),
		err:  err,
	}
}

func (r *restClient) get(path string, vals url.Values) Result {
	call, err := r.R().
		SetMultiValueQueryParams(vals).
		Get(r.rest + path)

	if r.debug {
		r.log.Debugw("rest_get", "path", path, "status", call.StatusCode(), "body", string(call.Body()))
	}
	return &restReturn{
		code: call.StatusCode(),
		body: call.Body(),
		err:  err,
	}
}

type urlQ struct {
	vals url.Values
}

func query(k, v string) *urlQ {
	return &urlQ{
		vals: url.Values{k: []string{v}},
	}
}
func (q *urlQ) V(k, v string) *urlQ {
	q.vals[k] = []string{v}
	return q
}

func (q *urlQ) Q() url.Values {
	return q.vals
}

func (r *restClient) getInfo() Result {
	call, err := r.R().
		Get(r.info)

	return &restReturn{
		code: call.StatusCode(),
		body: call.Body(),
		err:  err,
	}
}
