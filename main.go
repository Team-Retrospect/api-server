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
type SpanStruct struct {
  trace_id      string `json:"traceId"`
  span_id       string `json:"id"`
  session_id    string `json:"tags.frontendSession"`
  user_id       string `json:"tags.frontendUser"`

  // trigger_route string `json:spanTags["http.method"] + " " + spanTags["http.route"]"`,
  trigger_route string `json:"http.route"`

  time_sent     string `json:"timestamp"`
  status_code   string `json:"tags.http.status_code`
  span_data     string `json:"updatedSpanTags"`
}




/* connect to db */

func db_init() {}

func post_to_db() {}

/*
  var collection *mongo.Collection
  var collections = make(map[string]*mongo.Collection)
  var ctx = context.TODO()

  func db_init(uri string) {
    clientOptions := options.Client().ApplyURI(uri)

    client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
      log.Fatal(err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
      log.Fatal(err)
    }

    collection = client.Database("tracing").Collection("logs")
    for _, source := range(cfg.Sources) {
      collections[source] = client.Database("tracing").Collection(source)
    }
  }

  func post_to_db(collection *mongo.Collection, td TraceData) {
    if _, err := collection.InsertOne(context.TODO(), td); err != nil {
      output("Error writing!")
    } else {
      output("Write success, into:", td.Collection)
    }
  }
//*/



/* web server */
// func get_source_from_headers(headers map[string][]string) string {
//   header := headers[cfg.SourceHeader]
//   if headers == nil { return cfg.DefaultSource }
//   if len(header) == 0 { return cfg.DefaultSource }

//   source := header[0]
//   _, ok := collections[source]
//   if !ok { return cfg.DefaultSource }

//   return source
// }

func form_handler(w http.ResponseWriter, r *http.Request) {
  // if err := r.ParseForm(); err != nil { return }

  // enableCors(&w)

  // headers := headers_to_map(r.Header)

  // collectionName := get_source_from_headers(headers)
  // collection := collections[collectionName]

  // buf, _ := ioutil.ReadAll(r.Body)
  // content := fmt.Sprintf("%q", buf[:])

  // td := TraceData{headers, collectionName, content}

  // post_to_db(collection, td)
}

// func enableCors(w *http.ResponseWriter) {
//   allowedHeaders := "*" // "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token, Tracer-Source"
//   (*w).Header().Set("Access-Control-Allow-Origin", "*")
//   (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//   (*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
//   (*w).Header().Set("Access-Control-Expose-Headers", "Authorization")
// }

// func headers_to_map(headers http.Header) map[string][]string {
//   m := make(map[string][]string)
//   for h, va := range(headers) {
//     for _, v := range(va) { m[h] = append(m[h], v) }
//   }
//   return m
// }

func router(w http.ResponseWriter, r *http.Request) {
}

func root_path() string {
  return "hello"
}

func get_all_spans(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "this is a test");
}

func get_all_events(w http.ResponseWriter, r *http.Request) {

}

func insert_span(w http.ResponseWriter, r *http.Request) {

}

func insert_event(w http.ResponseWriter, r *http.Request) {

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
