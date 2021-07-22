package webserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/Team-Textrix/cassandra-connector/src/structs"
)

// initialze a database session variable
// --> this is used for DB queries
var session *gocql.Session

func DeclareRouter(cfg structs.Config, dbSession *gocql.Session) {
  session = dbSession

  fmt.Println("Declaring router...")
  r := mux.NewRouter()
  r.Path("/spans").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans)
  r.Path("/spans_by_trace/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_trace)
  r.Path("/spans_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_chapter)
  r.Path("/spans_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_session)
  r.Path("/events").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events)
  r.Path("/events_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events_by_chapter)
  r.Path("/span_search").Queries("trace_id", "{trace_id:[\\w\\-]*?}", "user_id", "{user_id:[\\w\\-]*?}", "session_id", "{session_id:[\\w\\-]*?}", "chapter_id", "{chapter_id:[\\w\\-]*?}", "status_code", "{status_code:[0-9]*?}").HandlerFunc(span_search_handler)
  r.Path("/event_search").Queries("user_id", "{user_id:[\\w\\-]*?}", "session_id", "{session_id:[\\w\\-]*?}", "chapter_id", "{chapter_id:[\\w\\-]*?}").HandlerFunc(event_search_handler)
  r.Path("/spans").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_spans)
  r.Path("/events").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_events)
  r.Path("/trigger_routes").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trigger_routes)
  r.Path("/trace_ids_by_trigger").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(get_all_trace_ids_by_trigger)
  r.Path("/chapter_ids_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_session)
  r.Path("/chapter_ids_by_trigger").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_trigger)
  r.Path("/events/snapshots").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_snapshots)
  r.Path("/events/snapshots").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_snapshots)

  http.Handle("/", r)

  header := handlers.AllowedHeaders([]string{"Accept", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token", "User-Id", "Session-Id", "Chapter-Id", "X-Requested-With", "Content-Type", "Authorization"})
  methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
  origins := handlers.AllowedOrigins([]string{"*"})

  fmt.Println("Now listening on", cfg.Port)
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
