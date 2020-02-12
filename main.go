package main

import (
	"github.com/ricdeau/in-mem-kv-storage/service"
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"log"
	"net/http"
)

const (
	route = "/api/storage/"
)

func main() {
	srv := service.New(route, storage.New())
	http.Handle(route, srv)

	if err := http.ListenAndServe(":5339", nil); err != nil {
		log.Fatal(err)
	}
}
