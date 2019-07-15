package rc

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	SortAscending  = 1
	SortDescending = -1
)

type Query struct {
	offset int
	count  int
	sort   map[string]int
}

func NewQuery() *Query {
	return &Query{
		sort: make(map[string]int),
	}
}

func (q *Query) Offset(o int) *Query {
	q.offset = o
	return q
}

func (q *Query) Count(c int) *Query {
	q.count = c
	return q
}

func (q *Query) Sort(field string, sort int) *Query {
	if sort == 1 || sort == -1 {
		q.sort[field] = sort
	}
	return q
}

func (q *Query) URLValues() url.Values {
	qry := url.Values{}
	if q.count > 0 {
		qry.Add("count", strconv.Itoa(q.count))
	}
	if q.offset > 0 {
		qry.Add("offset", strconv.Itoa(q.offset))
	}

	srt := []string{}
	for k, v := range q.sort {
		srt = append(srt, fmt.Sprintf("\"%s\":%d", k, v))
	}
	if len(srt) > 0 {
		qry.Add("sort", fmt.Sprintf("{%s}", strings.Join(srt, ",")))
	}
	return qry
}
