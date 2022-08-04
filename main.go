package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	nickandluke "github.com/lukemassa/nickandluke-api/internal/nickandluke"
)

func getPort() string {
	p := os.Getenv("PORT")
	if p != "" {
		return ":" + p
	}
	return ":3000"
}

func main() {
	var opts struct {
		Action string `long:"action" required:"true" default:"run" choice:"server" choice:"validate" choice:"upload"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	rh := nickandluke.RequestHandler()
	if opts.Action == "server" {
		http.HandleFunc("/guest", rh.CheckGuest)
		port := getPort()
		fmt.Printf("Listening on %s", port)
		panic(http.ListenAndServe(port, nil))

	}
	fmt.Println()
	fmt.Printf("%v", rh)
}
