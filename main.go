package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	token "github.com/lukemassa/jclubtakeaways-api/internal/token"
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
	if server {
		http.HandleFunc("/token", GetToken)
		port := getPort()
		fmt.Printf("Listening on %s", port)
		http.ListenAndServe(port, nil)
	}
	fmt.Printf(token.GetToken())
	fmt.Println()
}

// GetToken API to get a single token
func GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, token.GetToken())
}
