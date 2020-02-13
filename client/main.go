package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	total  = flag.Int("n", 10, "Total request")
	conc   = flag.Int("c", 2, "Max concurrent")
	target = flag.String("t", "", "Url for testing")
	method = flag.String("m", "GET", "Http method")
	cType  = flag.String("ct", "text/plain", "MIME Content type")
	file   = flag.String("f", "", "Path to file with content")
	data   = flag.String("d", "", "Content string")
)

func init() {
	flag.Parse()
}

func main() {
	sema := make(chan struct{}, *conc)

	var content []byte
	if *file != "" {
		content = readContentFromFile(*file)
	} else {
		content = []byte(*data)
	}

	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < *total; i++ {
		wg.Add(1)
		sema <- struct{}{}
		go func() {
			defer func() {
				wg.Done()
				<-sema
			}()
			body := bytes.NewBuffer(content)
			req, err := http.NewRequest(*method, *target, body)
			checkAndFailFast(err)
			if *cType != "" {
				req.Header.Add("Content-Type", *cType)
			}
			rsp, err := http.DefaultClient.Do(req)
			checkAndFailFast(err)
			defer rsp.Body.Close()
			_, err = io.Copy(ioutil.Discard, rsp.Body)
			checkAndFailFast(err)
		}()
	}
	wg.Wait()
	dur := time.Since(start)
	fmt.Printf("Total time: %s\n", dur)
	fmt.Printf("Average time per request: %s\n", dur/time.Duration(*total))
	fmt.Printf("Average rps: %.2f\n", float64(*total)/(dur.Seconds()))
}

func readContentFromFile(path string) []byte {
	f, err := os.Open(path)
	checkAndFailFast(err)
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	checkAndFailFast(err)
	return b
}

func checkAndFailFast(err error) {
	if err != nil {
		log.Fatalf("Fatal: %v", err)
	}
}
