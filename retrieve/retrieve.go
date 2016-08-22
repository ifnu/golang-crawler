package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("http://6brand.com.com")
	fmt.Println("http transport error is ", err)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("read error is ", err)
	fmt.Println(string(body))
}
