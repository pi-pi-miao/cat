package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"net/http"
	"sync"
)

var (
	addr = "http://127.0.0.1:8080"
	ch   = make(chan func(),number)
	number = 1000
	count = 40000
)

func main(){

	//httpGet()
	bench()
}


func httpGet(){
	resp,err := http.Get(addr)
	if err != nil {
		fmt.Println("[httpGet] response err :",err)
		return
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[httpGet] ioutil readAll err :", err)
		return
	}
	fmt.Println("[httpGet] body :",string(body))
}

func bench(){
	wg := &sync.WaitGroup{}
	wg.Add(count)
	start := time.Now()
	for i:=0 ; i < count; i ++ {
		ch <- func() {
			resp,err := http.Get(addr)
			if err != nil {
				fmt.Println("[httpGet] response err :",err)
				return
			}
			body,err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("[httpGet] ioutil readAll err :", err)
				return
			}
			fmt.Println("[httpGet] body :",string(body))
			resp.Body.Close()
			wg.Done()
		}
	}
	wg.Wait()
	fmt.Println("[getHttpGetTime] time is :",time.Since(start))
}

func init(){
	for i:= 0; i < number; i ++ {
		go func() {
			for v := range ch {
				v()
			}
		}()
	}
}