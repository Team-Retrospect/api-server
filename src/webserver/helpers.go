package webserver

import (
  "strconv"
  "fmt"
  // "time"
  // "strings"
  "net/http"
  "encoding/json"

  "github.com/Team-Textrix/cassandra-connector/src/structs"
)

func format_spans(blob []byte) []*structs.CassandraSpan {
  // initializing an array of SpanStructInput objects
  fmt.Println(string(blob))
  var jspans []*structs.SpanStructInput
  json.Unmarshal(blob, &jspans)
  cspans := make([]*structs.CassandraSpan, len(jspans))
  fmt.Println("cspans len:", len(cspans))

  for _, e := range jspans {
    if e == nil { continue }
    if e.Tags["http.method"] == "OPTIONS" { continue }

    sc, _ := strconv.ParseInt(e.Tags["http.status_code"], 10, 64)
    rd := []byte(e.Tags["requestData"])

    // if it's a db span, add frontend session info
    _, is_db := e.Tags["db.system"]
    // if is_db {
      // get trace id
      // tId := e.Trace_id

      // oneSpan := get_span_by_trace(tId)
      // if oneSpan == nil { continue }

      // e.Tags["frontendChapter"] = oneSpan.Chapter_id
      // e.Tags["frontendSession"] = oneSpan.Session_id
      // e.Tags["frontendUser"] = oneSpan.User_id
      // e.Tags["triggerRoute"] = oneSpan.Trigger_route
      // fmt.Println(e.Tags)
    // }

    delete(e.Tags, "requestData")

    tags := "{}"
    if len(e.Tags) > 0 {
      tags = "{"
      for k, v := range e.Tags {
        tags += fmt.Sprintf(`"%s": "%s", `, k, v)
      }
      tags = tags[0:len(tags)-2] + "}"
    }
    blob, _ := json.Marshal(tags)

    if is_db {
      cspans = append(cspans, &structs.CassandraSpan{
        Trace_id:       e.Trace_id,
        Span_id:        e.Span_id,
        Time_sent:      e.Time_sent,
        Duration:       strconv.Itoa(e.Duration) + "us",
        Session_id:     "",
        User_id:        "",
        Chapter_id:     "",
        Trigger_route:  "",
        Request_data:   rd,
        Status_code:    int16(sc),
        Data:           blob,
      })
    } else {
      cspans = append(cspans, &structs.CassandraSpan{
        Trace_id:       e.Trace_id,
        Span_id:        e.Span_id,
        Time_sent:      e.Time_sent,
        Duration:       strconv.Itoa(e.Duration) + "us",
        Session_id:     e.Tags["frontendSession"],
        User_id:        e.Tags["frontendUser"],
        Chapter_id:     e.Tags["frontendChapter"],
        Trigger_route:  e.Tags["triggerRoute"],
        Request_data:   rd,
        Status_code:    int16(sc),
        Data:           blob,
      })
    }

  }
  return cspans
}

// internal method used to fix frontend session values in db spans
func get_span_by_trace(traceId string) *structs.CassandraSpan {
  query := fmt.Sprintf("SELECT JSON chapter_id, user_id, session_id, trigger_route FROM project.spans WHERE trace_id='%s' LIMIT 1;", traceId)

  j := enumerate_query(query)

  var cspan *structs.CassandraSpan
  if len(j) > 0 { json.Unmarshal([]byte(j[0]), &cspan) }

  return cspan
}

func format_event(blob []byte, r *http.Request) *structs.CassandraEvent {
  // convert them into cassandra-compatible structs

  cevent := &structs.CassandraEvent{
    User_id:      r.Header.Get("user-id"),
    Session_id:   r.Header.Get("session-id"),
    Chapter_id:   r.Header.Get("chapter-id"),
    // Body:         strings.Replace(fmt.Sprint(string(blob)), "'", "\\'", -1),
    Body:         blob,
  }
  return cevent
}

func format_snapshot(blob []byte, r *http.Request) *structs.Snapshot {
  snapshot := &structs.Snapshot{
    Session_id:   r.Header.Get("session-id"),
    Body: blob,
  }
  return snapshot
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
