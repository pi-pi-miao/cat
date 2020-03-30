package cat

import (
	"bytes"
	"fmt"
	"strings"
)

func NewResonse() Response {
	return Response{
		Header: make(map[string][]string, 10),
	}
}

func (w  Response)Write(p []byte) (n int, err error){
	w.conn.Write(p)
	return n,nil
}

func (w *Response) SetCookie(v []string) {
	w.Header["Set-Cookie"] = v
}

func (w *Response) getCookie() bool {
	_, ok := w.Header["Set-Cookie"]
	if !ok {
		return false
	}
	return true
}

func (w *Response) Result(status string, v interface{}) {
	var buff bytes.Buffer
	var sbuff strings.Builder
	buff.WriteString(w.proto)
	buff.WriteString(" ")
	buff.WriteString(status)
	buff.WriteString("\n")
	for k, v := range w.Header {
		buff.WriteString(k)
		for _, v := range v {
			sbuff.WriteString(v)
		}
		buff.WriteString(":")
		buff.WriteString(sbuff.String())
		buff.WriteString("\r")
		sbuff.Reset()
	}
	buff.WriteString("\n")
	buff.WriteString(fmt.Sprintf("%v", v))
	result := buff.Bytes()
	buff.Reset()
	w.conn.Write(result)
	w.conn.Close()
	w.server.remove(w.conn)
}
