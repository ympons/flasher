package server

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

type Card struct {
	ID    int    `json:"id"`
	Typ   int    `json:"type" db:"type"`
	Front string `json:"front"`
	Back  string `json:"back"`
	Known bool   `json:"known"`
}

func (s Server) Index(w http.ResponseWriter, req *http.Request) {
	// TODO
	fmt.Fprintf(w, "%s", "Hello")
}

func (s *Server) Cards(w http.ResponseWriter, req *http.Request) {
	cards := []Card{}
	err := s.db.Select(&cards, "SELECT id, type, front, back, known FROM cards ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.render(w, req, "cards.html", pongo2.Context{
		"cards":       cards,
		"filter_name": "all",
	})
}

func (s Server) FilterCards(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	filterName := vars["name"]

	filters := map[string]string{
		"all":     "where 1 = 1",
		"general": "where type = 1",
		"code":    "where type = 2",
		"known":   "where known = 1",
		"unknown": "where known = 0",
	}
	query := filters[filterName]

	if query == "" {
		http.Redirect(w, req, "/cards", http.StatusFound)
	}

	cards := []Card{}
	err := s.db.Select(&cards, "SELECT id, type, front, back, known FROM cards "+query+" ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.render(w, req, "cards.html", pongo2.Context{
		"cards":       cards,
		"filter_name": filterName,
	})
}

func (s Server) AddCard(w http.ResponseWriter, req *http.Request) {
	// TODO
	fmt.Fprintf(w, "%s", "Hello")
}

func (s Server) EditCard(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var card Card
	err := s.db.Get(&card, "SELECT id, type, front, back, known FROM cards WHERE id=$1", vars["id"])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.render(w, req, "edit.html", pongo2.Context{
		"card": card,
	})
}
