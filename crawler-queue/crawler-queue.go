package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackdanger/collectlinks"
)

func main() {

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}

	queue := make(chan string)

	go func() {
		queue <- args[0]
	}()

	for uri := range queue {
		enqueue(uri, queue)
	}
}

func enqueue(uri string, queue chan string) {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{
		Transport: transport,
	}
	start := makeTimeStamp()
	resp, err := client.Get(uri)
	end := makeTimeStamp()
	fmt.Println("fetching", uri, end-start, "ms")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		go func() { queue <- link }()
	}

}

func makeTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
