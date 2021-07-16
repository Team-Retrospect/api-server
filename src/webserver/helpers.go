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
  // initializing an array of SpanStructInput objects
  var jspans []*SpanStructInput
  json.Unmarshal(blob, &jspans)
  cspans := make([]*CassandraSpan, len(jspans))

  for _, e := range jspans {
    if e == nil { continue }
    if e.Tags["http.method"] == "OPTIONS" { continue }

    sc, _ := strconv.ParseInt(e.Tags["http.status_code"], 10, 64)
    rd := e.Tags["requestData"]

    // if it's a db span, add frontend session info
    _, ok := e.Tags["db.system"]
    if ok {
      // get trace id
      tId := e.Trace_id

      oneSpan := get_span_by_trace(tId)

      e.Tags["frontendChapter"] = oneSpan.Chapter_id
      e.Tags["frontendSession"] = oneSpan.Session_id
      e.Tags["frontendUser"] = oneSpan.User_id
      e.Tags["triggerRoute"] = oneSpan.Trigger_route
    }

    delete(e.Tags, "requestData")

    tags := "{"
    for k, v := range e.Tags {
      tags += fmt.Sprintf(`"%s": "%s", `, k, v)
    }
    tags = tags[0:len(tags)-2] + "}"

    cspans = append(cspans, &CassandraSpan{
      Trace_id:      e.Trace_id,
      Span_id:       e.Span_id,
      Time_sent:     e.Time_sent,
      Duration:      strconv.Itoa(e.Duration) + "us",
      Session_id:    e.Tags["frontendSession"],
      User_id:       e.Tags["frontendUser"],
      Chapter_id:    e.Tags["frontendChapter"],
      Trigger_route: e.Tags["triggerRoute"],
      Request_data:  rd,
      Status_code:   int16(sc),
      Data: strings.Replace(fmt.Sprint(tags), "'", "\\'", -1),
    })
  }
  return cspans
}

// internal method used to fix frontend session values in db spans
func get_span_by_trace(traceId string) *CassandraSpan {
  query := fmt.Sprintf("SELECT JSON chapter_id, user_id, session_id, trigger_route FROM project.spans WHERE trace_id='%s' LIMIT 1;", traceId)

  j := enumerate_query(query)

  var cspan *CassandraSpan
  json.Unmarshal([]byte(j[0]), &cspan)

  return cspan
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