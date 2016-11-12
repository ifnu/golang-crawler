package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/jackdanger/collectlinks"
)

var visited = make(map[string]bool)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("please specify arguments")
		os.Exit(1)
	}

	queue := make(chan string)

	go func() {
		queue <- args[0]
	}()

	for uri := range queue {
		avoidErrorUri := uri
		enqueue(avoidErrorUri, queue)
	}

}

func enqueue(url string, queue chan string) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	start := makeTimeStamp()
	resp, err := client.Get(url)
	end := makeTimeStamp()
	fmt.Println(url, end-start, "ms")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		absolute := fixURL(link, url)
		if absolute != "" {
			if !visited[absolute] {
				go func() {
					queue <- absolute
				}()
			}
		}
	}

}

func fixURL(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseURL.ResolveReference(uri)
	return uri.String()
}

func makeTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
