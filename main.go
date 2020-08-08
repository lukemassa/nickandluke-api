package main

import (
    "os"
    "fmt"
    "net/http"
)

func getPort() string {
  p := os.Getenv("PORT")
  if p != "" {
    return ":" + p
  }
  return ":3000"
}

func main() {
    http.HandleFunc("/", HelloServer)
    port := getPort()
    fmt.Printf("Listening on %s", port)
    http.ListenAndServe(":3000", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
