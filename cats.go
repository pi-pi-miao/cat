package cat

import (
	"log"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	stopCh = make(chan interface{})
	once   = &sync.Once{}
)

type (
	GroupRoute struct {
		*Cats
		F        Fn
		GroupRoute   *strings.Builder
	}
	Cats struct {
		*GroupRoute
		GroupMap     map[string]Fn      
		Middleware   map[string][]middleFn  
		MaxGoroutine int
		route        string 
	}
)

func NewCats() *Cats {
	c := &Cats{
		GroupRoute:&GroupRoute{
			GroupRoute:&strings.Builder{},
		},
		GroupMap:   make(map[string]Fn, 10),
		Middleware: make(map[string][]middleFn,10),
	}
	c.GroupRoute.Cats = c
	return c
}

func (c *Cats) Run(addr string)error {
	if s,err := newServer(c,addr);err == nil {
		s.run()
	}else {
		return err
	}
	return nil
}

func (c *Cats) Miao(r *Request, w Response) {
	fn, getRoute := c.GroupMap[r.Method+r.Uri]
	switch {
	case !getRoute:
		goto loop
	default:
		fn(r, w)
		return
	}
loop:
	c.NotMethod(w)
	return
}

func (c *Cats) Group(route string) *GroupRoute {
	cat := &GroupRoute{
		Cats:c,
		GroupRoute:&strings.Builder{},
	}
	cat.GroupRoute.WriteString(route)
	return cat
}

func (c *GroupRoute) Get(route string, f Fn)*GroupRoute{
	c.addMethod("GET", route, f)
	return c
}

func (c *GroupRoute) Post(route string, f Fn)*GroupRoute {
	c.addMethod("POST", route, f)
	return c
}

func (c *GroupRoute) Delete(route string, f Fn)*GroupRoute{
	c.addMethod("DELETE", route, f)
	return c
}

func (c *GroupRoute) Patch(route string, f Fn)*GroupRoute {
	c.addMethod("PATCH", route, f)
	return c
}

func (c *GroupRoute) Put(route string, f Fn)*GroupRoute {
	c.addMethod("PUT", route, f)
	return c
}

func (c *Cats) NotMethod(w Response) {
	w.Result("404", "the method not found")
}

func (c *GroupRoute)addMethod(method, route string, f Fn) {
	switch {
	case c.GroupRoute.Len() == 0:
		c.GroupMap[method+route] = f
		c.route = method+route

	case c.GroupRoute.Len() >= 1:
		groupRoute := c.GroupRoute.String()
		route := join(groupRoute, route)
		c.GroupMap[method+route] = f
		c.route = method+route
	}
}

func join(group, route string) string {
	var buf strings.Builder
	buf.Reset()
	switch {
	case group[len(group)-1] == '/' && route[len(route)-1] != '/':
		buf.WriteString(group)
		buf.WriteString(route)
		return buf.String()
	default:
		buf.WriteString(group)
		buf.WriteString(route)
		return buf.String()
	}
}

func (c *Cats)Metric() {
	for {
		select {
		case <-stopCh:
			return
		default:
			log.Println(runtime.NumGoroutine())
			time.Sleep(5 * time.Second)
		}
	}
}

func (s *server) ShutDown() {
	s.Cats.ShutDown()
}

func (c *Cats) ShutDown() {
	for k, _ := range c.GroupMap {
		delete(c.GroupMap, k)
	}
	once.Do(func() {
		close(stopCh)
		log.Println("the server is closed")
	})
}
