package main

import (
	"log"
	"net/http"
	"runtime"

	"github.com/codegangsta/martini"
	"marti/control"
	"marti/models"
)

func init() {
	models.InitDB()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	m := martini.Classic()

	defer models.Close()

	m.Get("/", control.HomeGet)
	m.Get("/add", control.AddFuncGet)
	m.Post("/add", control.AddFuncPost)
	m.Get("/view", control.View)
	m.Get("/autocomplete", control.Autocomplete)

	log.Println("listen on Port:9090")
	http.ListenAndServe(":9090", m)
}
