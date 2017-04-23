package server

import (
	"encoding/gob"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/ympons/flasher/db"
)

type Server struct {
	credentials string
	router      *mux.Router
	sessions    *sessions.CookieStore
	db          *db.DB
	basePath    string
}

func New(basePath, secretKey, credentials string, db *db.DB) *Server {
	s := &Server{
		credentials: credentials,
		basePath:    basePath,
		db:          db,
		sessions:    sessions.NewCookieStore([]byte(secretKey)),
		router:      mux.NewRouter(),
	}

	s.router.HandleFunc("/login", s.Login).Methods("GET", "POST")
	s.router.HandleFunc("/logout", s.Logout).Methods("GET")

	s.router.HandleFunc("/", s.admin(s.Index)).Methods("GET")
	s.router.HandleFunc("/general", s.admin(s.General)).Methods("GET")
	s.router.HandleFunc("/general/{id:[0-9]+}", s.admin(s.General)).Methods("GET")
	s.router.HandleFunc("/code", s.admin(s.Code)).Methods("GET")
	s.router.HandleFunc("/code/{id:[0-9]+}", s.admin(s.Code)).Methods("GET")
	s.router.HandleFunc("/cards", s.admin(s.Cards)).Methods("GET")
	s.router.HandleFunc("/filter_cards/{name}", s.admin(s.FilterCards)).Methods("GET")
	s.router.HandleFunc("/create", s.admin(s.CreateCard)).Methods("POST")
	s.router.HandleFunc("/edit/{id:[0-9]+}", s.admin(s.EditCard)).Methods("GET")
	s.router.HandleFunc("/update", s.admin(s.UpdateCard)).Methods("POST")
	s.router.HandleFunc("/delete/{id:[0-9]+}", s.admin(s.DeleteCard)).Methods("GET")
	s.router.HandleFunc("/mark_known/{id:[0-9]+}/{type:general|code}", s.admin(s.MarkKnown)).Methods("GET")

	s.router.PathPrefix("/static/").Handler(http.FileServer(http.Dir(basePath)))
	gob.Register(&flash{})

	return s
}

func (s Server) Run(host string) {
	log.Printf("Listening on %s", host)
	log.Fatal(http.ListenAndServe(host, s.router))
}

func (s *Server) Close() error {
	return s.db.Close()
}
