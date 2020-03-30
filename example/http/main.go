package main

import "net/http"

var (
	addr = "127.0.0.1:8080"
)

func main(){
	http.HandleFunc("/",handler)
	http.ListenAndServe(addr,nil)
}

func handler(w http.ResponseWriter,r *http.Request) {
	w.Write([]byte("hello world"))
}
