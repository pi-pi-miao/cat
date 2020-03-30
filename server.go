package cat

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"syscall"
)

type (
	context struct {
		request  *Request
		response Response
		f  Fn
	}
	Handlers interface {
		Miao(request *Request, response Response)
	}

	Response struct {
		Header map[string][]string
		Body   interface{}
		Uri    string
		proto  string
		conn   net.Conn
		server *server
	}
	server struct {
		*Cats
		addr string
		handlers Handlers
		max  int
		fd    int
		lock  *sync.RWMutex
		m  map[int]net.Conn
		network string
		work *work
	}
	work struct {
		max  int
		once *sync.Once
	}
	Fn func(*Request,Response)
	middleFn func(*Request, *Response)
)

func newServer(c *Cats,addr string)(*server,error){
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &server{
		Cats:     c,
		addr:     addr,
		handlers: c,
		fd:fd,
		lock:     &sync.RWMutex{},
		m:        make(map[int]net.Conn,1000),
		network:"tcp",
	},nil
}


func (s *server) handler() {
	p,_ := ants.NewPool(1000)
	conns := make([]net.Conn,0,1000)
	buf := make([]byte,1024)
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		err := s.wait(&conns)
		if err != nil || len(conns) == 0{
			continue
		}
		for k,_ := range conns {
			conn := conns[k]
			b, err := conn.Read(buf)
			if err != nil {
				if err != nil {
					s.remove(conns[k])
					continue
				}
			}
			req := NewRequst()
			resp := NewResonse()
			resp.server = s
			req.AllHeader = strings.Split(string(buf[:b]), "\n\r")[0]
			req.readHeader()
			req.Body = strings.Split(string(buf[:b]), "\n\r")[1]
			req.parseRequestLine()
			resp.conn = conn
			resp.proto = req.proto
			p.Submit(func() {
				s.handlers.Miao(req,resp)
			})
		}
	}
}

func (s *server) run() {
	listen, err := net.Listen(s.network, s.addr)
	if err != nil {
		fmt.Println("listen err", err)
		return
	}
	fmt.Println("cat web framework is running at ",s.addr)
	go func() {
		s.handler()
	}()
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		conn, err := listen.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}
			log.Printf("accept err: %v", err)
			return
		}
		if err := s.add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}


func (e *server) add(conn net.Conn) error {
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	e.lock.Lock()
	e.m[fd] = conn
	e.lock.Unlock()
	return nil
}

func (e *server) remove(conn net.Conn) error {
	fd := socketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.lock.Lock()
	delete(e.m,fd)
	e.lock.Unlock()
	return nil
}
func (e *server) wait(conns *[]net.Conn)error{
	*conns = (*conns)[0:0]
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		e.lock.RLock()
		conn,ok := e.m[int(events[i].Fd)]
		if ok {
			*conns = append(*conns,conn.(net.Conn))
			e.lock.RUnlock()
			continue
		}
		e.lock.RUnlock()
	}
	return nil
}
func socketFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

