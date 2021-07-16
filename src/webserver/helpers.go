package webserver

import (
  "strconv"
  "fmt"
  "strings"
  "net/http"
  "encoding/json"
  "io"

  "github.com/Team-Textrix/cassandra-connector/src/structs"
)

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
  // convert them into cassandra-compatible structs
  cevents := make([]*structs.CassandraEvent, 1)

    cevents = append(cevents, &structs.CassandraEvent{
      User_id:            r.Header.Get("user-id"),
      Session_id:         r.Header.Get("session-id"),
      Chapter_id:         r.Header.Get("chapter-id"),
      Body: string(blob),
    })
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
