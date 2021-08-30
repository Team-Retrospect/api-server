package webserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/Team-Retrospect/api-server/src/structs"
)

// initialze a database session variable
// --> this is used for DB queries
var session *gocql.Session

func DeclareRouter(cfg structs.Config, dbSession *gocql.Session) {
	session = dbSession

	fmt.Println("Declaring router...")
	r := mux.NewRouter()
	r.Path("/spans").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans)
	r.Path("/spans").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_spans)
	r.Path("/spans_by_trace/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_trace)
	r.Path("/spans_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_chapter)
	r.Path("/spans_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_spans_by_session)
	r.Path("/events").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events)
	r.Path("/events").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_events)
	r.Path("/events_by_chapter/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events_by_chapter)
	r.Path("/events_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_events_by_session)
	r.Path("/trigger_routes").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_trigger_routes)
	r.Path("/trace_ids_by_trigger").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(get_all_trace_ids_by_trigger)
	r.Path("/chapter_ids_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_session)
	r.Path("/chapter_ids_by_trigger").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(get_all_chapter_ids_by_trigger)
	r.Path("/events/snapshots").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_snapshots)
	r.Path("/events/snapshots_by_session/{id}").Methods(http.MethodGet, http.MethodOptions).HandlerFunc(get_all_snapshots_by_session)
	r.Path("/events/snapshots").Methods(http.MethodPost, http.MethodOptions).HandlerFunc(insert_snapshots)

	http.Handle("/", r)

	header := handlers.AllowedHeaders([]string{"Accept", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token", "User-Id", "Session-Id", "Chapter-Id", "X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	fmt.Println("Now listening on", cfg.Port)
	if cfg.UseHTTPS {
		if err := http.ListenAndServeTLS(cfg.Port, cfg.FullCert, cfg.PrivateKey, handlers.CORS(header, methods, origins)(r)); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := http.ListenAndServe(cfg.Port, handlers.CORS(header, methods, origins)(r)); err != nil {
			log.Fatal(err)
		}
	}
}
