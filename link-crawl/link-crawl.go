package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/jackdanger/collectlinks"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("please specify parameter")
		os.Exit(1)
	}
	resp, err := http.Get(args[0])
	if err != nil {
		return
	}
	defer resp.Body.Close()
	links := collectlinks.All(resp.Body)
	for _, link := range links {
		fmt.Println(link)
	}
}
