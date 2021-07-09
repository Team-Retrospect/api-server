package main

import (
  /* debug */
  "fmt"
  "log"
  "strconv"
  "strings"

  /* config */
  "github.com/ilyakaznacheev/cleanenv"

  /* database */
  "github.com/gocql/gocql"
  "time"
  "encoding/json"


  /* webserver */
  "net/http"
  // "io/ioutil"
  "io"
  "github.com/gorilla/mux"
)

func output(contents ...string) {
  if (debug) { fmt.Println(contents) }
}

type TraceData struct {
  Headers map[string][]string
  Collection string
  Content string
}

/* load configs from config.yml */
type ConfigStruct struct {
  Debug         bool    `yaml:"debug"`

  Cluster       string  `yaml:"cluster"`

  Port          string  `yaml:"port"`

  UseHTTPS      bool    `yaml:"use_https"`
  FullCert      string  `yaml:"fullcert"`
  PrivateKey    string  `yaml:"privatekey"`
}
var debug bool = false;
var cfg ConfigStruct
func load_cfg() {
  cleanenv.ReadConfig("config.yml", &cfg)
  debug = cfg.Debug
}



/* connect to db */

var cluster *gocql.ClusterConfig
var session *gocql.Session




// zipkin data format
type SpanStructInput struct {
  Trace_id          string        `json:"traceId"`
  Span_id           string        `json:"id"`
  Time_sent         string        `json:"timestamp"`
  Duration          int           `json:"duration"`

  Trigger_route     string        // `json:"trigger_route"`
  Session_id        string        // `json:"session_id"`
  User_id           string        // `json:"user_id"`
  Status_code       string        // `json:"status_code"`
  Chapter_id        string       // `json:"chapter_id"`
  Tags              map[string]string        `json:"tags"`
}

type CassandraSpan struct {
  Trace_id          string        `json:"trace_id"`
  Span_id           string        `json:"span_id"`
  Time_sent         string        `json:"time_sent"`
  Duration          string        `json:"time_duration"`
  Data              string        `json:"data"`
  // Data              map[string]string `json:"data"`

  Trigger_route     string        `json:"trigger_route"`
  User_id           string        `json:"user_id"`
  Session_id        string        `json:"session_id"`
  Chapter_id        string        `json:"chapter_id"`
  Status_code       int16         `json:"status_code"`
  Request_data      string        `json:"request_data"`
}

type EventStructInput struct {
  // body data only
  // Type              string        `json:"type"`
  // Data              string        `json:"data"`
  // Event_time        string        `json:"timestamp"`
}

type CassandraEvent struct {
  // header data
  User_id           string        `json:"user_id"`
  Session_id        string        `json:"session_id"`
  Chapter_id        string        `json:"chapter_id"`

  // body data
  // Type              string        `json:"type"`
  // Data              string        `json:"data"`
  // Time_sent         string        `json:"time_sent"`
  // Test string `json:"test"`
  Body  string `json:"data"`
}


/* web server */

func root_path() string {
  return "hello"
}

func get_all_spans(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM project.spans;"
  scanner := session.Query(query).Iter().Scanner()

  var j []string
  for scanner.Next() {
    var s string
    scanner.Scan(&s)
    j = append(j, s)
  }

  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

func get_all_events(w http.ResponseWriter, r *http.Request) {
  query := "SELECT JSON * FROM project.events;"
  scanner := session.Query(query).Iter().Scanner()

  var j []string
  for scanner.Next() {
    var s string
    scanner.Scan(&s)
    j = append(j, s)
  }

  js := fmt.Sprintf("[%s]", strings.Join(j, ", "))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
  output("Retrieved events", js)
}

func insert_spans(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  cspans := format_spans(body)

  for _, span := range(cspans) {
    if span == nil { continue }
    j, _ := json.Marshal(span)

    query := "INSERT INTO project.spans JSON '" + string(j) + "';"
    session.Query(query).Exec()
  }

  w.WriteHeader(http.StatusOK)
}

func format_spans(blob []byte) []*CassandraSpan {
  // unmarshal the json blob
  var jspans []*SpanStructInput
  json.Unmarshal(blob, &jspans)
  fmt.Println(jspans)

  // for i, v := range(jspans) { fmt.Println("jspan", i, v) }

  // convert them into cassandra-compatible structs
  cspans := make([]*CassandraSpan, len(jspans))

  for _, e := range(jspans) {
    if e == nil { continue }
    sc, _ := strconv.ParseInt(e.Tags["http.status_code"], 10, 64)
    rd := e.Tags["requestData"]
    delete(e.Tags, "requestData")

    tags := "{";
    for k, v := range(e.Tags) { tags += fmt.Sprintf(`"%s": "%s", `, k, v) }
    tags = tags[0:len(tags)-2] + "}"

    cspans = append(cspans, &CassandraSpan{
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
      Data:           strings.Replace(fmt.Sprint(tags), "'", "\\'", -1),
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
    // fmt.Println(query)
  }

  w.WriteHeader(http.StatusOK)
}

// todo
func format_events(blob []byte, r *http.Request) []*CassandraEvent {
  // unmarshal the json blob
  // var jevents []*EventStructInput
  // json.Unmarshal(blob, &jevents)
  // fmt.Println("blob", string(blob))
  // fmt.Println("jevents", jevents)
  // fmt.Println("len jevents", len(jevents))

  // convert them into cassandra-compatible structs
  // cevents := make([]*CassandraEvent, len(jevents))
  cevents := make([]*CassandraEvent, 1)

  // for i, e := range(jevents) {
    // fmt.Println(i, e)
    // if e == nil { continue }

    // tags := "{";
    // for k, v := range(e.Tags) { tags += fmt.Sprintf(`"%s": "%s", `, k, v) }
    // tags = tags[0:len(tags)-2] + "}"

    cevents = append(cevents, &CassandraEvent{
      User_id:            r.Header.Get("user-id"),
      Session_id:         r.Header.Get("session-id"),
      Chapter_id:         r.Header.Get("chapter-id"),

      // Type:               e.Type,
      // Data:               e.Data,
      // Time_sent:          e.Event_time,
      // Test: e.Test,
      Body: string(blob),
    })
  // }
  return cevents
}



/* orchestrate */
func main() {
  load_cfg()

  /* connect to cassandra here */
  output("Connecting to Cassandra...")
  // connect to the cluster
  cluster = gocql.NewCluster(cfg.Cluster)
  cluster.Consistency = gocql.Quorum
  cluster.ProtoVersion = 4
  cluster.ConnectTimeout = time.Second * 10
  // cluster.Authenticator = gocql.PasswordAuthenticator{Username: "Username", Password: "Password"} //replace the username and password fields with their real settings.
  s, err := cluster.CreateSession()
  if err != nil {
    log.Println(err)
    return
  }
  session = s
  defer session.Close()



  output("Declaring router...")
  r := mux.NewRouter()
  // r.HandleFunc("/",       root_path)
  r.HandleFunc("/spans",  get_all_spans   ).Methods("GET")
  r.HandleFunc("/events", get_all_events  ).Methods("GET")
  r.HandleFunc("/spans",  insert_spans    ).Methods("POST")
  r.HandleFunc("/events", insert_events   ).Methods("POST")
  http.Handle("/", r)

  output("Now listening on", cfg.Port)
  if (cfg.UseHTTPS) {
    if err := http.ListenAndServeTLS(cfg.Port, cfg.FullCert, cfg.PrivateKey, nil); err != nil {
      log.Fatal(err)
    }
  } else {
    if err := http.ListenAndServe(cfg.Port, nil); err != nil {
      log.Fatal(err)
    }
  }
}
