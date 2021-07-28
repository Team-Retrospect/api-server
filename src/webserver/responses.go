package webserver

import (
  "fmt"
  "net/http"
)

func send_json_response(w http.ResponseWriter, js string) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, js)
}
