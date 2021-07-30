package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// --> GET /spans
func get_all_spans(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM retrospect.spans;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /spans_by_trace/{id}
func get_all_spans_by_trace(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  trace_id, ok := vars["id"]

  if !ok {
    send_incorrect_params(w)
    return
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.spans WHERE trace_id='%s';", trace_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /spans_by_chapter/{id}
func get_all_spans_by_chapter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    send_incorrect_params(w)
    return
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.spans WHERE chapter_id='%s';", chapter_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /spans_by_session/{id}
func get_all_spans_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    send_incorrect_params(w)
    return
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.spans WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /events
func get_all_events(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM retrospect.events;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /events_by_chapter/{id}
func get_all_events_by_chapter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    send_incorrect_params(w)
    return
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.events WHERE chapter_id='%s';", chapter_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /events_by_session/{id}
func get_all_events_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    send_incorrect_params(w)
    return
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.events WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> POST /spans
func insert_spans(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Inserting a Span")
  // r.Body is type *http.bodyblob
  // io.ReadAll returns an array of bytes
  body, err := io.ReadAll(r.Body)
  if err != nil {
    send_incorrect_params(w)
    return
  } else if len(body) == 0 {
    send_missing_body(w)
    return
  }

  // format_spans takes an array of bytes
  // and returns an array of structs.CassandraSpan objects
  cspans := format_spans(body)

  for _, span := range(cspans) {
    if span == nil { continue }
    // json.Marshal returns the json encoding of the variable passed into it
    j, _ := json.Marshal(span)

    table := "spans"
    if span.Session_id == "" { table = "db_span_buffer" }

    // each json-ified span is stringified and inserted into the database as a json object
    query := "INSERT INTO retrospect." + (table) + " JSON '" + string(j) + "';"
    fmt.Println(query)
    session.Query(query).Exec()
  }

  send_creation_success(w)
}

// --> POST /events
func insert_events(w http.ResponseWriter, r *http.Request) {
  body, err := io.ReadAll(r.Body)

  if err != nil {
    send_incorrect_params(w)
    return
  } else if len(body) == 0 {
    send_missing_body(w)
    return
  }

  cevent := format_event(body, r)
  if cevent == nil { return }

  j, _ := json.Marshal(cevent)
  query := "INSERT INTO retrospect.events JSON '" + string(j) + "';"
  session.Query(query).Exec()

  send_creation_success(w)
}

// --> GET /events/snapshots
func get_snapshots(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM retrospect.snapshots;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /events/snapshots_by_session/{id}
func get_all_snapshots_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    send_incorrect_params(w)
    return
  }

  query := fmt.Sprintf("SELECT JSON * FROM retrospect.snapshots WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> POST /events/snapshots
func insert_snapshots(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Inserting a snapshot")
  body, _ := io.ReadAll(r.Body)
  snapshot := format_snapshot(body, r)
  if snapshot == nil { return }

  j, _ := json.Marshal(snapshot)
  query := "INSERT INTO retrospect.snapshots JSON '" + string(j) + "';"
  session.Query(query).Exec()

  send_creation_success(w)
}

// --> GET /trigger_routes
func get_all_trigger_routes(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON trigger_route, data FROM retrospect.spans;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /trace_ids_by_trigger
func get_all_trace_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, err := io.ReadAll(r.Body)

  if err != nil {
    send_incorrect_params(w)
    return
  } else if len(body) == 0 {
    send_missing_body(w)
    return
  }

  trigger_route := string(body)

  query := fmt.Sprintf("SELECT trace_id FROM retrospect.spans WHERE trigger_route='%s' ALLOW FILTERING;", trigger_route);

  j := enumerate_query(query)
  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  send_json_response(w, js)
}

// --> GET /chapters_by_session/{id}
func get_all_chapter_ids_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    send_incorrect_params(w)
    return
  }

  query := fmt.Sprintf("SELECT JSON chapter_id FROM retrospect.spans WHERE session_id='%s';", session_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  send_json_response(w, js)
}

// --> GET /chapter_ids_by_trigger
func get_all_chapter_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, err := io.ReadAll(r.Body)

  if err != nil {
    send_incorrect_params(w)
    return
  } else if len(body) == 0 {
    send_missing_body(w)
    return
  }

  target := string(body)

  query := fmt.Sprintf("SELECT chapter_id FROM retrospect.spans WHERE trigger_route='%v' ALLOW FILTERING;", target);

  j := enumerate_query(query)
  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  send_json_response(w, js)
}
