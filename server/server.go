package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ympons/flasher/db"
)

type Server struct {
	router   *mux.Router
	sessions *sessions.CookieStore
	db       *db.DB
	basePath string
}

func New(basePath, secretKey string, db *db.DB) *Server {
	s := &Server{
		basePath: basePath,
		db:       db,
		sessions: sessions.NewCookieStore([]byte(secretKey)),
		router:   mux.NewRouter(),
	}

	s.router.HandleFunc("/general", s.admin(s.Cards)).Methods("GET")
	s.router.HandleFunc("/cards", s.admin(s.Cards)).Methods("GET")
	s.router.HandleFunc("/filter_cards/{name}", s.admin(s.FilterCards)).Methods("GET")
	s.router.HandleFunc("/add", s.admin(s.AddCard)).Methods("POST")
	s.router.HandleFunc("/edit/{id:[0-9]+}", s.admin(s.EditCard)).Methods("GET")

	s.router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./web/")))

	return s
}

func (s Server) Run(host string) {
	log.Printf("Listening on %s", host)
	log.Fatal(http.ListenAndServe(host, s.router))
}

func (s *Server) Close() error {
	return s.db.Close()
}
