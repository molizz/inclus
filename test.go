package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {
	httpUrl := "https://github.com/molizz/inclus"

	httpURI, _ := url.Parse(httpUrl)

	path := strings.Trim(httpURI.Path, "/")
	n := strings.Split(path, "/")
	fmt.Println(n)

}
