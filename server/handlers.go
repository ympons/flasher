package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
)

type Card struct {
	ID    int    `json:"id"`
	Type  uint8  `json:"type" db:"type"`
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

func (s Server) CreateCard(w http.ResponseWriter, req *http.Request) {
	var (
		typ, _ = strconv.Atoi(req.Form.Get("type"))
		front  = req.Form.Get("front")
		back   = req.Form.Get("back")
	)

	if _, err := s.db.Exec("INSERT INTO cards (type, front, back) VALUES (?, ?, ?)",
		typ,
		front,
		back,
	); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.flash(w, req, infoAlert, "New card was successfully added.")
	http.Redirect(w, req, "/cards", http.StatusFound)
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

func (s Server) UpdateCard(w http.ResponseWriter, req *http.Request) {
	var known = false
	if req.Form.Get("known") == "1" {
		known = true
	}

	if _, err := s.db.Exec("UPDATE cards SET type = ?, front = ?, back = ?, known = ? WHERE id = ?",
		req.Form.Get("type"),
		req.Form.Get("front"),
		req.Form.Get("back"),
		known,
		req.Form.Get("card_id"),
	); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.flash(w, req, infoAlert, "Card was successfully updated.")
	http.Redirect(w, req, "/cards", http.StatusFound)
}
