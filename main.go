package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
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
		Action string `long:"action" required:"true" default:"run" choice:"server" choice:"validate" choice:"download" choice:"download-and-server"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	})

	if err != nil {
		panic(err)
	}
	rh := nickandluke.RequestHandler()
	dh := nickandluke.DataHandler(sess)
	if opts.Action == "download" || opts.Action == "download-and-server" {
		err := dh.Download()
		if err != nil {
			panic(err)
		}
	}
	if opts.Action == "server" || opts.Action == "download-and-server" {
		http.HandleFunc("/guest", rh.CheckGuest)
		port := getPort()
		fmt.Printf("Listening on %s", port)
		panic(http.ListenAndServe(port, nil))

	}
	if opts.Action == "validate" {
		fmt.Println()
		fmt.Printf("%v", rh)
	}

}
