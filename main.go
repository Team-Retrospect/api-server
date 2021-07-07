package main

import (
  /* debug */
  "fmt"
  "log"

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



// zipkin data format
type SpanStructInput struct {
  Trace_id      string `json:"traceId"`
  Span_id       string `json:"id"`

  Time_sent     int64  `json:"timestamp"`
  Duration      int64  `json:"duration"`
  Data          map[string]string `json:"tags"`

  Status_code   string
  Session_id    string
  User_id       string
  Trigger_route string
}



type SpanStructFromDB struct {
  trace_id      string `json:"traceId"`
  span_id       string `json:"id"`
  session_id    string `json:"tags.frontendSession"`
  user_id       string `json:"tags.frontendUser"`

  // trigger_route string `json:spanTags["http.method"] + " " + spanTags["http.route"]"`,
  trigger_route string `json:"http.route"`

  time_sent     string `json:"timestamp"`
  status_code   int16  `json:"tags.http.status_code`
  span_data     string `json:"updatedSpanTags"`
}



/* connect to db */

var cluster *gocql.ClusterConfig

func db_init() {}

func post_to_db() {}



/* web server */

func root_path() string {
  return "hello"
}

func get_all_spans(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w,
  `[{"trace_id": "1", "span_id": "2", "session_id": "3", "user_id": "4", "time_sent": "1299038700000", "trigger_route": "/some/place", "status_code": 200, "data": {} }]`)
}

func get_all_events(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  fmt.Fprintf(w,
  `[{"id": "1", "session_id": "2", "user_id": "3", "time_sent": "1299038700000", "data": {} }]`)
}

func insert_spans(w http.ResponseWriter, r *http.Request) {
  body, _ := io.ReadAll(r.Body)

  var spans []*SpanStructInput

  json.Unmarshal(body, &spans)

  for _, span := range(spans) {
    span.Status_code = span.Data["http.status_code"]
    span.Session_id = span.Data["frontendSession"]
    span.User_id = span.Data["frontendUser"]
    span.Trigger_route = span.Data["http.method"] + " " + span.Data["http.route"]
  }

  for _, span := range(spans) {
    // query here
    fmt.Println(span.Trace_id)
  }

  w.WriteHeader(http.StatusOK)
}

func insert_events(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
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
  session, err := cluster.CreateSession()
        if err != nil {
    log.Println(err)
    return
  }
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
