package main

import (
	/* debug */
	"fmt"
	"log"
	"strconv"
	"strings"

	/* config */
	// "Clean and minimalistic environment configuration reader for Golang"
	// "reads and parses configuration structure from the file
	// reads and overwrites configuration structure from environment variables
	// writes a detailed variable list to help output"
	//   ReadConfig method takes a string representing the name of a config file
	//    and a pointer to a struct (ex. ConfigStruct)

	// may have to use this when dockerizing
	"github.com/ilyakaznacheev/cleanenv"

	/* database */
	"encoding/json"
	"time"

	// gocql implements a fast and robust Cassandra client for the Go programming language
	"github.com/gocql/gocql"

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
// property names match the data in the config.yml


// may have to add database credentials to this struct
// when we dockerize the application
// (read from environment variables? err := cleanenv.ReadEnv(&cfg))
//     Password string `env:"PASSWORD"`


type ConfigStruct struct {
// debug: false
  Debug         bool    `yaml:"debug"`

// cluster: "cassandra.xadi.io"
  Cluster       string  `yaml:"cluster"`

// port: ":443"
  Port          string  `yaml:"port"`

// use_https: true
  UseHTTPS      bool    `yaml:"use_https"`
// fullcert: "/etc/letsencrypt/live/api.xadi.io/fullchain.pem"
  FullCert      string  `yaml:"fullcert"`
// privatekey: "/etc/letsencrypt/live/api.xadi.io/privkey.pem"
  PrivateKey    string  `yaml:"privatekey"`
}

var debug bool = false;
var cfg ConfigStruct
// taking the information from the .yml file and putting it into a Struct
func load_cfg() {
  cleanenv.ReadConfig("config.yml", &cfg)
  debug = cfg.Debug
}

/* connect to db */

// initializing a variable of type gocql.ClusterConfig
// we will eventually set it equal to gocql.NewCluster(cfg.Cluster)
// Q for Nicole: Why are these initialized outlide of main function?
var cluster *gocql.ClusterConfig

// initialize a variable of type gocql.Session
// will eventually be set equal to cluster.CreateSession()
// (see above to see what cluster is equal to)
var session *gocql.Session

// the cluster holds nodes (like layers) which hold tables
// a cluster holds nodes which would be like, two different localhosts
// (in this case we only have one)
// the nodes hold keyspaces (like a table)


// zipkin data format
type SpanStructInput struct {
  Trace_id          string        `json:"traceId"`
  Span_id           string        `json:"id"`
  Time_sent         int           `json:"timestamp"`
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
  Time_sent         int           `json:"time_sent"`
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
  if (cfg.UseHTTPS) { enableCors(&w) }

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

// r.Path("/spans_by_trace/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_trace)
func get_all_spans_by_trace(w http.ResponseWriter, r *http.Request) {
  if (cfg.UseHTTPS) { enableCors(&w) }

  vars := mux.Vars(r);
  trace_id, ok := vars["id"]

  if !ok {
    fmt.Println("trace_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans WHERE trace_id='%s' ALLOW FILTERING;", trace_id);
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

// r.Path("/spans_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_chapter)
func get_all_spans_by_chapter(w http.ResponseWriter, r *http.Request) {
  if (cfg.UseHTTPS) { enableCors(&w) }

  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans WHERE chapter_id='%s' ALLOW FILTERING;", chapter_id);
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
  if (cfg.UseHTTPS) { enableCors(&w) }

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

// r.Path("/events_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events_by_chapter)
func get_all_events_by_chapter(w http.ResponseWriter, r *http.Request) {
  if (cfg.UseHTTPS) { enableCors(&w) }

  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.events WHERE chapter_id='%s' ALLOW FILTERING;", chapter_id);
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

func insert_spans(w http.ResponseWriter, r *http.Request) {
  output("Inserting a Span")
  if (cfg.UseHTTPS) { enableCors(&w) }

  // r.Body is type *http.bodyblob
  // io.ReadAll returns an array of bytes
  body, _ := io.ReadAll(r.Body)

  // format_spans takes an array of bytes
  // and returns an array of CassandraSpan objects
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

func get_all_trigger_routes(w http.ResponseWriter, r *http.Request) {
  if (cfg.UseHTTPS) { enableCors(&w) }

  query := "SELECT JSON trigger_route FROM project.spans;"
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

// r.Path("/trace_ids/{trigger_route}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trace_ids_by_trigger_route)
func get_all_trace_ids_by_trigger_route(w http.ResponseWriter, r *http.Request) {
  if (cfg.UseHTTPS) { enableCors(&w) }

  vars := mux.Vars(r);
  trigger_route, ok := vars["trigger_route"]

  if !ok {
    fmt.Println("trigger_route is missing in parameters")
  }

  tre := strings.Fields(trigger_route)
  trigger_route = tre[0] + " " + tre[1] + "//" + tre[2] + "/" + tre[3]
  fmt.Println(`trigger_route=`, trigger_route)

  query := fmt.Sprintf("SELECT JSON trace_id FROM project.spans WHERE trigger_route='%s' ALLOW FILTERING;", trigger_route);
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

// r.Path("/chapters_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_session)
func get_all_chapter_ids_by_session(w http.ResponseWriter, r *http.Request) {
  if (cfg.UseHTTPS) { enableCors(&w) }

  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    fmt.Println("session_id is missing in parameters")
  }

  fmt.Println(`session id=`, session_id)

  query := fmt.Sprintf("SELECT JSON chapter_id FROM project.spans WHERE session_id='%s' ALLOW FILTERING;", session_id);
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

func format_spans(blob []byte) []*CassandraSpan {
  // unmarshal the json blob
  
  // initializing an array of SpanStructInput objects
  var jspans []*SpanStructInput
  json.Unmarshal(blob, &jspans)

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
      // note for nicole: what is the point of the string replacement?
      // We can't see any single quotes in the tags, nor any double
      // slashes in the resulting Data
      Data:           strings.Replace(fmt.Sprint(tags), "'", "\\'", -1),
    })
  }
  return cspans
}



func insert_events(w http.ResponseWriter, r *http.Request) {
  if (cfg.UseHTTPS) { enableCors(&w) }

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
    fmt.Println(r.Header)

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

func enableCors(w *http.ResponseWriter) {
  allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token, Tracer-Source, User-Id, Session-Id, Chapter-Id"
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  (*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
  (*w).Header().Set("Access-Control-Expose-Headers", "Authorization")
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
  r.Path("/spans").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans)
  r.Path("/spans_by_trace/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_trace)
  r.Path("/spans_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_chapter)
  r.Path("/events").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events)
  r.Path("/events_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events_by_chapter)
  r.Path("/spans").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_spans)
  r.Path("/events").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_events)
  r.Path("/trigger_routes").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trigger_routes)
  r.Path("/trace_ids/{trigger_route}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trace_ids_by_trigger_route)
  r.Path("/chapters_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_session)
  http.Handle("/", r)

  output("Now listening on", cfg.Port)
  // if (cfg.UseHTTPS) {
  //   if err := http.ListenAndServeTLS(cfg.Port, cfg.FullCert, cfg.PrivateKey, nil); err != nil {
  //     log.Fatal(err)
  //   }
  // } else {
    if err := http.ListenAndServe(cfg.Port, nil); err != nil {
      log.Fatal(err)
    }
  // }
}
