package webserver

import (
  "fmt"
  "strings"
  "strconv"
  "net/http"
  "github.com/gorilla/mux"
  "encoding/json"
  "io"

  "github.com/Team-Textrix/cassandra-connector/structs"
)

// --> GET /spans
func get_all_spans(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM project.spans;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /spans_by_trace/{id}
func get_all_spans_by_trace(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  trace_id, ok := vars["id"]

  if !ok {
    fmt.Println("trace_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans WHERE trace_id='%s';", trace_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /spans_by_chapter/{id}
func get_all_spans_by_chapter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans WHERE chapter_id='%s';", chapter_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /events
func get_all_events(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM project.events;"

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
  fmt.Println("Retrieved events", js)
}

// --> GET /events_by_chapter/{id}
func get_all_events_by_chapter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.events WHERE chapter_id='%s';", chapter_id);

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> POST /spans
func insert_spans(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Inserting a Span")
  // r.Body is type *http.bodyblob
  // io.ReadAll returns an array of bytes
  body, _ := io.ReadAll(r.Body)

  // format_spans takes an array of bytes
  // and returns an array of structs.CassandraSpan objects
  cspans := format_spans(body)

  for _, span := range(cspans) {
    if span == nil { continue }
    // json.Marshal returns the json encoding of the variable passed into it
    j, _ := json.Marshal(span)

    // each json-ified span is stringified and inserted into the database as a json object
    query := "INSERT INTO project.spans JSON '" + string(j) + "';"
    session.Query(query).Exec()
  }

  w.WriteHeader(http.StatusOK)
}

// --> GET /trigger_routes
func get_all_trigger_routes(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON trigger_route, data FROM project.spans;"
  fmt.Println(query)

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /trace_ids_by_trigger
func get_all_trace_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  trigger_route := string(body)

  query := fmt.Sprintf("SELECT trace_id FROM project.spans WHERE trigger_route='%s' ALLOW FILTERING;", trigger_route);
  fmt.Println(query)

  j := enumerate_query(query)
  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /chapters_by_session/{id}
func get_all_chapter_ids_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    fmt.Println("session_id is missing in parameters")
    // TODO: return this as error
  }

  query := fmt.Sprintf("SELECT JSON chapter_id FROM project.spans WHERE session_id='%s';", session_id);
  fmt.Println(query)

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET /chapter_ids_by_trigger
func get_all_chapter_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  target := string(body)

  query := fmt.Sprintf("SELECT chapter_id FROM project.spans WHERE trigger_route='%v' ALLOW FILTERING;", target);
  fmt.Println(query)

  j := enumerate_query(query)
  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET? /span_search
func span_search_handler(w http.ResponseWriter, r *http.Request) {
  acceptedParams := []string {
    "trace_id",
    "user_id",
    "session_id",
    "chapter_id",
    "status_code",
  }

  var dynamicQuery []string

  for _, p := range acceptedParams {
    val := r.FormValue(p)
    if val != "" {
      if p != "status_code" {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v='%v'", p, val))
      } else {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v=%v", p, val))
      }
    }
  }

  dynamicQueryString := strings.Join(dynamicQuery," AND ")
  fmt.Println(dynamicQueryString)

  if len(dynamicQueryString) != 0 {
    dynamicQueryString = "WHERE " + dynamicQueryString + " ALLOW FILTERING"
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans " + dynamicQueryString + ";")
  fmt.Println(query)

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// --> GET? /event_search
func event_search_handler(w http.ResponseWriter, r *http.Request) {
  acceptedParams := []string {
    "user_id",
    "session_id",
    "chapter_id",
  }

  var dynamicQuery []string

  for _, p := range acceptedParams {
    val := r.FormValue(p)
    if val != "" {
      if p != "status_code" {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v='%v'", p, val))
      } else {
        dynamicQuery = append(dynamicQuery, fmt.Sprintf("%v=%v", p, val))
      }
    }
  }

  dynamicQueryString := strings.Join(dynamicQuery," AND ")
  fmt.Println(dynamicQueryString)

  if len(dynamicQueryString) != 0 {
    dynamicQueryString = "WHERE " + dynamicQueryString + " ALLOW FILTERING"
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.events " + dynamicQueryString + ";")

  j := enumerate_query(query)
  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

func format_spans(blob []byte) []*structs.CassandraSpan {
  // initializing an array of structs.SpanStructInput objects
  var jspans []*structs.SpanStructInput
  json.Unmarshal(blob, &jspans)

  // convert them into cassandra-compatible structs
  cspans := make([]*structs.CassandraSpan, len(jspans))

  for _, e := range(jspans) {
    if e == nil { continue }
    sc, _ := strconv.ParseInt(e.Tags["http.status_code"], 10, 64)
    request_data := e.Tags["requestData"]
    delete(e.Tags, "requestData")

    // encase in JSON syntax
    tags := "{";
    for k, v := range(e.Tags) { tags += fmt.Sprintf(`"%s": "%s", `, k, v) }
    tags = tags[0:len(tags)-2] + "}"
    // escape where necessary to safeguard against injections
    tags = strings.Replace(fmt.Sprint(tags), "'", "\\'", -1) //

    cspans = append(cspans, &structs.CassandraSpan{
      Trace_id:       e.Trace_id,
      Span_id:        e.Span_id,
      Time_sent:      e.Time_sent,
      Duration:       strconv.Itoa(e.Duration) + "us",
      Session_id:     e.Tags["frontendSession"],
      User_id:        e.Tags["frontendUser"],
      Chapter_id:     e.Tags["frontendChapter"],
      Trigger_route:  e.Tags["triggerRoute"],
      Request_data:   request_data,
      Status_code:    int16(sc),
      Data:           tags,
    })
  }
  return cspans
}

func insert_events(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  cevents := format_events(body, r)

  for _, event := range(cevents) {
    if event == nil { continue }
    j, _ := json.Marshal(event)

    query := "INSERT INTO project.events JSON '" + string(j) + "';"
    session.Query(query).Exec()
  }

  w.WriteHeader(http.StatusOK)
}

func format_events(blob []byte, r *http.Request) []*structs.CassandraEvent {
  // convert them into cassandra-compatible structs\
  cevents := make([]*structs.CassandraEvent, 1)

    cevents = append(cevents, &structs.CassandraEvent{
      User_id:            r.Header.Get("user-id"),
      Session_id:         r.Header.Get("session-id"),
      Chapter_id:         r.Header.Get("chapter-id"),
      Body: string(blob),
    })
  // }
  return cevents
}

func enumerate_query(query string) (sa []string) {
  scan := session.Query(query).Iter().Scanner()

  for scan.Next() {
    var s string
    scan.Scan(&s)
    sa = append(sa, s)
  }

  return
}
