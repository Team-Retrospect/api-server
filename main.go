package main

import (
  /* debug */
  "fmt"
  "log"

  /* config */
  "github.com/ilyakaznacheev/cleanenv"

  /*  */
  // "encoding/json"

  /* webserver */
  "net/http"
  // "io/ioutil"
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

func insert_span(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func insert_event(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

/* orchestrate */
func main() {
  load_cfg()

  output("(not connecting to DB)")
  /* connect to cassandra here */

  output("Declaring router...")
  r := mux.NewRouter()
  // r.HandleFunc("/",       root_path)
  r.HandleFunc("/spans",  get_all_spans   ).Methods("GET")
  r.HandleFunc("/events", get_all_events  ).Methods("GET")
  r.HandleFunc("/span",   insert_span     ).Methods("POST")
  r.HandleFunc("/event",  insert_event    ).Methods("POST")
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
