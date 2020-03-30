package main

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"time"
)

func demoFunc(i int) {
	fmt.Println("Hello World!",i)
}

func main() {
	p, _ := ants.NewPool(10000)
	p.Submit(func() {
		demoFunc(1)
	})
	time.Sleep(5*time.Second)
}