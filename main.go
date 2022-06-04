package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	nickandluke "github.com/lukemassa/nickandluke-api/internal/nickandluke"
)

var server bool

func init() {
	flag.BoolVar(&server, "server", false, "Whether or not to run the server")
	flag.Parse()
}

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":3000"
}

func main() {
	rh := nickandluke.RequestHandler()
	if server {
		http.HandleFunc("/guest", rh.CheckGuest)
		port := getPort()
		fmt.Printf("Listening on %s", port)
		http.ListenAndServe(port, nil)
	}
	fmt.Println()

	fmt.Printf("%v", rh)
}
