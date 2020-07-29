package main

import (
	"time"

	"github.com/pi-pi-miao/cat"
)

func main(){
	cats := cat.NewCats()
	cats.Get("/",do)
	cats.Post("/",do)
	cats.Delete("/",do)
	cats.Patch("/",do)
	cats.Put("/",do)
	go cats.Run(":8080")
	time.Sleep(5*time.Second)
	cats.ShutDown()   //server closed gracefully
}

func do(r *cat.Request,w cat.Response){
	w.Result("200","hello world")
}
