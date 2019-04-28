package cat

import (
	"net"
	"sync"
	"bufio"
	"io"
	"fmt"
	"strings"
	"log"
)

var TaskCh chan []interface{}

type (
	Handlers interface {
		Miao(request *Request, response Response)
	}

	Response struct {
		Header map[string][]string
		Body   interface{}
		Uri    string
		proto  string
		conn   net.Conn
	}
	Server struct {
		*Cats
		Addr string
		W    Handlers
		Max  int
	}
	Work struct {
		Max  int
		once *sync.Once
	}
	Fn func(*Request,Response)
	middleFn func(*Request, *Response)
)



func (s *Server) handler(conn net.Conn) {
	defer conn.(net.Conn).Close()
	req := NewRequst()
	resp := NewResonse()
	r := bufio.NewReader(conn.(net.Conn))
	buff := make([]byte, 2048)
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		b, err := r.Read(buff)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read byte err", err)
				return
			}
		}
		req.AllHeader = strings.Split(string(buff[:b]), "\n\r")[0]
		req.readHeader()
		req.Body = strings.Split(string(buff[:b]), "\n\r")[1]
		req.parseRequestLine()
		resp.conn = conn.(net.Conn)
		resp.proto = req.proto
		s.W.Miao(req, resp)
		return
	}
}

func (s *Server) Run() {
	listen, err := net.Listen("tcp", s.Addr)
	work := Pool(s.Max).NewPool()
	if err != nil {
		fmt.Println("listen err", err)
		goto loop
	}
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept err")
			goto loop
		}
		connContainer := make([]interface{}, 0, 2)
		connContainer = append(connContainer, conn)
		connContainer = append(connContainer, s.handler)
		work.Do(connContainer)
	}
loop:
	return
}

func Pool(max int) *Work {
	work := &Work{
		Max:  max,
		once: &sync.Once{},
	}
	TaskCh = make(chan []interface{}, max)
	return work
}

func (w *Work) NewPool() *Work {
	for i := 0; i < w.Max; i++ {
		go func() {
			defer func(){
				if err := recover();err != nil {
					log.Println("this goroutine err",err)
				}
			}()
			for task := range TaskCh {
				conn := task[0].(net.Conn)
				f := task[1].(func(net.Conn))
				f(conn)
			}
		}()
	}
	return w
}

func (w *Work) Do(f []interface{}) {
	TaskCh <- f
}

func (w *Work) Close() {
	w.once.Do(func() {
		close(TaskCh)
	})
}


