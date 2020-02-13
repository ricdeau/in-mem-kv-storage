package main

import (
	"flag"
	"fmt"
	"github.com/ricdeau/in-mem-kv-storage/logger"
	"github.com/ricdeau/in-mem-kv-storage/middleware"
	"github.com/ricdeau/in-mem-kv-storage/service"
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	route = "/api/storage/"
)

var (
	verbosity    = flag.Bool("v", false, "Print logs to stdout")
	port         = flag.Int("p", 5339, "Server listening port")
	maxKeySize   = flag.Int("mks", 0, "Max key size in bytes. If <=0 - unlimited")
	maxValueSize = flag.Int("mvs", 0, "Max value size in bytes. If <=0 - unlimited")
)

func init() {
	flag.Parse()
}

func main() {
	if !*verbosity {
		log.SetOutput(ioutil.Discard)
	}
	srv := service.New(route, storage.New())
	srv = middleware.LimitsMiddleware(srv, route, *maxKeySize, *maxValueSize)
	srv = middleware.RequestIDMiddleware(srv)
	http.Handle(route, srv)

	logger.Infof("Server start listening on port %d", *port)
	logger.Infof("Operations url: %s", route)
	logger.Infof("Max key size: %s", sizeToString(*maxKeySize))
	logger.Infof("Max value size: %s", sizeToString(*maxValueSize))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		logger.Fatalf("Server start error: %v", err)
	}
}

func sizeToString(size int) string {
	if size > 0 {
		return fmt.Sprintf("%d bytes", size)
	}
	return "unlimited"
}
