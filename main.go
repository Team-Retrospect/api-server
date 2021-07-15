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

	"github.com/gorilla/handlers"
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

// cluster: "cassandra.domain.com"
  Cluster       string  `yaml:"cluster"`

// port: ":443"
  Port          string  `yaml:"port"`

// use_https: true
  UseHTTPS      bool    `yaml:"use_https"`
// fullcert: "/etc/letsencrypt/live/api.domain.com/fullchain.pem"
  FullCert      string  `yaml:"fullcert"`
// privatekey: "/etc/letsencrypt/live/api.domain.com/privkey.pem"
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

// initialze a database session variable
// --> this is used for DB queries
var session *gocql.Session

// the cluster holds nodes (like layers) which hold tables
// a cluster holds nodes which would be like, two different localhosts
// (in this case we only have one)
// the nodes hold keyspaces (like a table)

// zipkin data format
type SpanStructInput struct {
  Trace_id          string              `json:"traceId"`
  Span_id           string              `json:"id"`
  Time_sent         int                 `json:"timestamp"`
  Duration          int                 `json:"duration"`

  Tags              map[string]string   `json:"tags"`
}

// span insertion struct
type CassandraSpan struct {
  // top-level span data
  Trace_id          string        `json:"trace_id"`
  Span_id           string        `json:"span_id"`
  Time_sent         int           `json:"time_sent"`
  Duration          string        `json:"time_duration"`
  Data              string        `json:"data"`

  // derived from tags
  Trigger_route     string        `json:"trigger_route"`
  User_id           string        `json:"user_id"`
  Session_id        string        `json:"session_id"`
  Chapter_id        string        `json:"chapter_id"`
  Status_code       int16         `json:"status_code"`
  Request_data      string        `json:"request_data"`
}

// event insertion struct
type CassandraEvent struct {
  User_id           string        `json:"user_id"`
  Session_id        string        `json:"session_id"`
  Chapter_id        string        `json:"chapter_id"`
  Body              string        `json:"data"`
}

type TriggerRoutePayload struct {
  Route string `json:"trigger_route"`
}

/* web server */

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

// r.Path("/spans_by_trace/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_trace)
func get_all_spans_by_trace(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  trace_id, ok := vars["id"]

  if !ok {
    fmt.Println("trace_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans WHERE trace_id='%s';", trace_id);
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
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.spans WHERE chapter_id='%s';", chapter_id);
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

// r.Path("/events_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events_by_chapter)
func get_all_events_by_chapter(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  chapter_id, ok := vars["id"]

  if !ok {
    fmt.Println("chapter_id is missing in parameters")
  }

  query := fmt.Sprintf("SELECT JSON * FROM project.events WHERE chapter_id='%s';", chapter_id);
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
  query := "SELECT JSON trigger_route, data FROM project.spans;"
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

//   r.Path("/trace_ids_by_trigger").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trace_ids_by_trigger)
func get_all_trace_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  trigger_route := string(body)

  query := fmt.Sprintf("SELECT trace_id FROM project.spans WHERE trigger_route='%s' ALLOW FILTERING;", trigger_route);
  scanner := session.Query(query).Iter().Scanner()

  var j []string
  for scanner.Next() {
    var s string
    scanner.Scan(&s)
    j = append(j, s)
  }

  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// r.Path("/chapters_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_session)
func get_all_chapter_ids_by_session(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r);
  session_id, ok := vars["id"]

  if !ok {
    fmt.Println("session_id is missing in parameters")
  }

  fmt.Println(`session id=`, session_id)

  query := fmt.Sprintf("SELECT JSON chapter_id FROM project.spans WHERE session_id='%s';", session_id);
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

// r.Path("/chapter_ids_by_trigger").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_trigger)
func get_all_chapter_ids_by_trigger(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)
  target := string(body)

  query := fmt.Sprintf("SELECT chapter_id FROM project.spans WHERE trigger_route='%v' ALLOW FILTERING;", target);
  scanner := session.Query(query).Iter().Scanner()

  var j []string
  for scanner.Next() {
    var s string
    scanner.Scan(&s)
    j = append(j, s)
  }

  js := fmt.Sprintf("[\"%s\"]", strings.Join(j, "\", \""))

  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w, js)
}

// r.Path("/span_search").Queries("trace_id", "{[a-zA-Z0-9]*?}").Queries("status_code", "{[0-9]*?}").HandlerFunc(span_search_handler)
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
  fmt.Println("query", query)
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
  fmt.Println("query", query)
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
  // initializing an array of SpanStructInput objects
  var jspans []*SpanStructInput
  json.Unmarshal(blob, &jspans)

  // convert them into cassandra-compatible structs
  cspans := make([]*CassandraSpan, len(jspans))

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

    cspans = append(cspans, &CassandraSpan{
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

func format_events(blob []byte, r *http.Request) []*CassandraEvent {
  // convert them into cassandra-compatible structs\
  cevents := make([]*CassandraEvent, 1)

    cevents = append(cevents, &CassandraEvent{
      User_id:            r.Header.Get("user-id"),
      Session_id:         r.Header.Get("session-id"),
      Chapter_id:         r.Header.Get("chapter-id"),
      Body: string(blob),
    })
  // }
  return cevents
}

/* orchestrate */
func main() {
  load_cfg()

  // connect to the Cassandra cluster
  output("Connecting to Cassandra...")
  cluster := gocql.NewCluster(cfg.Cluster)
  cluster.Consistency = gocql.Quorum
  cluster.ProtoVersion = 4
  cluster.ConnectTimeout = time.Second * 10
  // cluster.Authenticator = gocql.PasswordAuthenticator{
  //   Username: "Username",
  //   Password: "Password",
  // } //replace the username and password fields with their real settings.
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
  r.Path("/span_search").Queries("trace_id", "{trace_id:[\\w\\-]*?}", "user_id", "{user_id:[\\w\\-]*?}", "session_id", "{session_id:[\\w\\-]*?}", "chapter_id", "{chapter_id:[\\w\\-]*?}", "status_code", "{status_code:[0-9]*?}").HandlerFunc(span_search_handler)
  r.Path("/event_search").Queries("user_id", "{user_id:[\\w\\-]*?}", "session_id", "{session_id:[\\w\\-]*?}", "chapter_id", "{chapter_id:[\\w\\-]*?}").HandlerFunc(event_search_handler)
  r.Path("/spans").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_spans)
  r.Path("/events").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_events)
  r.Path("/trigger_routes").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trigger_routes)
  r.Path("/trace_ids_by_trigger").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trace_ids_by_trigger)
  r.Path("/chapter_ids_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_session)
  r.Path("/chapter_ids_by_trigger").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_trigger)
  http.Handle("/", r)

  header := handlers.AllowedHeaders([]string{"Accept", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token", "User-Id", "Session-Id", "Chapter-Id", "X-Requested-With", "Content-Type", "Authorization"})
  methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
  origins := handlers.AllowedOrigins([]string{"*"})

  output("Now listening on", cfg.Port)
  if (cfg.UseHTTPS) {
    if err := http.ListenAndServeTLS(cfg.Port, cfg.FullCert, cfg.PrivateKey, handlers.CORS(header, methods, origins)(r)); err != nil {
      log.Fatal(err)
    }
  } else {
    if err := http.ListenAndServe(cfg.Port, handlers.CORS(header, methods, origins)(r)); err != nil {
      log.Fatal(err)
    }
  }
}
