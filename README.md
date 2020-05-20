# CAT



service webframework



## use

```
go get github.com/pi-pi-miao/cat
```



## How to use

### first

```GO
package main

import (
	"cat"
)

func main(){
    cats := cat.NewCats() 
    cats.Get("/",do)
    cats.Post("/",do)
    cats.Delete("/",do)
    cats.Patch("/",do)
    cats.Put("/",do)
    go cats.Run(":8080")
    time.sleep(5*time.Second)
    cats.ShutDown()   //server closed gracefully
}

func do(r *cat.Request,w cat.Response){
	w.Result("200","hello world")
}
```

### second

```go
package main

import (
	"cat"
)

func main(){
         cats := cat.NewCats()
	 g1 := cats.Group("/")
	 g1.Get("/v1", Do)
	 g1.Post("/v1", Do)
	 g1.Delete("/v1", Do)
	 g1.Patch("/v1", Do)
	 g1.Put("/v1", Do)
	 cat2 := cats.Group("/app")
	 cat2.Get("/v2",Do)
	 cats.Run(":8080")
}

func do(r *cat.Request,w cat.Response){
	w.Result("200","hello world")
}
```
## 注意此项目后续会不断维护，如果使用中有什么问题欢迎主动联系
