package main

import (
	"cat"
)

func main(){
	cats := cat.NewCats()
	cats.Get("/",do)
	cats.Run(":8080")
}

func do(r *cat.Request,w cat.Response){
	w.Result("200","hello world")
}