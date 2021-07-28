package webserver

import (
  "fmt"
  "net/http"
)

const (
  // 2xx
  CreationSuccess = "Creation was successful"

  // 3xx

  // 4xx
  IncorrectParams = "Required parameters are missing or incorrect"
  MissingBody = "No content supplied in the request body"

  // 5xx
)

// 200 OK
func send_ok(w http.ResponseWriter, c ...string) {
  var b string
  if len(c) > 0 { b = c[0] }

  w.WriteHeader(http.StatusOK)
  fmt.Fprintf(w, b)
}

// 200 OK
func send_json_response(w http.ResponseWriter, js string) {
  w.Header().Set("Content-Type", "application/json")
  send_ok(w, js)
}

// 201 SuccessfulCreation
func send_creation_success(w http.ResponseWriter) {
  w.WriteHeader(http.StatusCreated)
  fmt.Fprintf(w, CreationSuccess)
}

// 400 Bad Request
func send_incorrect_params(w http.ResponseWriter) {
  w.WriteHeader(http.StatusBadRequest)
  fmt.Fprintf(w, IncorrectParams)
}

// 400 Bad Request
func send_missing_body(w http.ResponseWriter) {
  w.WriteHeader(http.StatusBadRequest)
  fmt.Fprintf(w, MissingBody)
}
