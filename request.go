package cat

import (
	"errors"
	"strings"
)

type (
	Request struct {
		Method    string
		Header    map[string][]string
		Body      interface{}
		Uri       string
		proto     string
		AllHeader string
	}
)

func NewRequst() *Request {
	return &Request{
		Header: make(map[string][]string, 10),
	}
}

func (r *Request) GetCookie() ([]string, error) {
	cookie, ok := r.Header["cookie"]
	if !ok {
		return nil, errors.New("get cookie err")
	}
	return cookie, nil
}

func (r *Request) GetUserAgent() ([]string, error) {
	userAgent, ok := r.Header["User-Agent"]
	if !ok {
		return nil, errors.New("get userAgent err")
	}
	return userAgent, nil
}

func (r *Request) SetUserAgent(v []string) {
	r.Header["User-Agent"] = v
}

func (r *Request) GetMethod() string {
	return r.Method
}

func (r *Request) GetHost() ([]string, error) {
	host, ok := r.Header["Host"]
	if !ok {
		return nil, errors.New("get host err")
	}
	return host, nil
}

func (r *Request) parseRequestLine() {
	line := strings.Split(strings.Split(r.AllHeader, "\r")[0], " ")
	r.Method = line[0]
	r.Uri = line[1]
	r.proto = line[2]
}

func (r *Request) readHeader() {
	requestHeaders := strings.Split(r.AllHeader, "\r\n")
	for _, requestHeader := range requestHeaders[1:] {
		v := strings.Split(requestHeader, ":")
		r.Header[v[0]] = []string{v[1]}
	}
}
