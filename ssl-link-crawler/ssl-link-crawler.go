package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/jackdanger/collectlinks"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 0 {
		fmt.Println("Please specify parameter")
		os.Exit(1)
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client := http.Client{
		Transport: transport,
	}

	resp, err := client.Get(args[0])
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)

	for _, link := range links {
		fmt.Println(link)
	}

}
