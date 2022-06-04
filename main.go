package main

import (
	"encoding/json"
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

type requestHandler struct {
	guests map[string]string
}

type checkResponse struct {
	Valid bool   `json:"valid"`
	Form  string `json:"form"`
}

func getRequestHandler() requestHandler {
	guests := make(map[string]string)
	guests["luke"] = "https://tripadvisor.com"
	guests["nick"] = "https://twitter.com"
	return requestHandler{
		guests: guests,
	}
}

func (rh requestHandler) CheckGuest(w http.ResponseWriter, r *http.Request) {
	res := checkResponse{}
	name := r.URL.Query().Get("name")
	if val, ok := rh.guests[name]; ok {
		res.Valid = true
		res.Form = val
	}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(js)

}

func main() {
	rh := getRequestHandler()
	if server {
		http.HandleFunc("/token", GetToken)
		http.HandleFunc("/nickandluke/guest", rh.CheckGuest)
		port := getPort()
		fmt.Printf("Listening on %s", port)
		http.ListenAndServe(port, nil)
	}
	fmt.Printf(token.GetToken())
	fmt.Println()

	fmt.Printf("%v", rh)
}

// GetToken API to get a single token
func GetToken(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, token.GetToken())
}
