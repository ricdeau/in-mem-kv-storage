package main

import (
	"flag"
	"fmt"
	"github.com/ricdeau/in-mem-kv-storage/logger"
	"github.com/ricdeau/in-mem-kv-storage/middleware"
	"github.com/ricdeau/in-mem-kv-storage/service"
	"github.com/ricdeau/in-mem-kv-storage/storage"
	"net/http"
)

const (
	route = "/api/storage/"
)

var (
	port         = flag.Int("port", 5339, "Server listening port")
	maxKeySize   = flag.Int("max-key-size", 0, "Max key size in bytes. If <=0 - unlimited")
	maxValueSize = flag.Int("max-value-size", 0, "Max value size in bytes. If <=0 - unlimited")
)

func init() {
	flag.Parse()
}

func main() {
	srv := service.New(route, storage.New())
	srv = middleware.LimitsMiddleware(srv, route, *maxKeySize, *maxValueSize)
	srv = middleware.LoggingMiddleware(srv)
	http.Handle(route, srv)

	logger.Infof("Server start listening at port %d", *port)
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
