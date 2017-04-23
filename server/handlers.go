package server

import (
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

func (s *Server) Index(w http.ResponseWriter, req *http.Request) {
	s.General(w, req)
}

func (s *Server) Code(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.memorize(w, req, "code", vars["id"])
}

func (s *Server) General(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	s.memorize(w, req, "general", vars["id"])
}

func (s *Server) Cards(w http.ResponseWriter, req *http.Request) {
	cards := []Card{}
	err := s.db.Select(&cards, "SELECT id, type, front, back, known FROM cards ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.render(w, req, "cards.html", pongo2.Context{
		"cards":       cards,
		"filter_name": "all",
	})
}

func (s *Server) FilterCards(w http.ResponseWriter, req *http.Request) {
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
		return
	}

	cards := []Card{}
	err := s.db.Select(&cards, "SELECT id, type, front, back, known FROM cards "+query+" ORDER BY id DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.render(w, req, "cards.html", pongo2.Context{
		"cards":       cards,
		"filter_name": filterName,
	})
}

func (s *Server) CreateCard(w http.ResponseWriter, req *http.Request) {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.flash(w, req, infoAlert, "New card was successfully added.")
	http.Redirect(w, req, "/cards", http.StatusFound)
}

func (s *Server) EditCard(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var card Card
	err := s.db.Get(&card, "SELECT id, type, front, back, known FROM cards WHERE id=$1", vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.render(w, req, "edit.html", pongo2.Context{
		"card": card,
	})
}

func (s *Server) UpdateCard(w http.ResponseWriter, req *http.Request) {
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.flash(w, req, infoAlert, "Card was successfully updated.")
	http.Redirect(w, req, "/cards", http.StatusFound)
}

func (s *Server) DeleteCard(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	if _, err := s.db.Exec("DELETE FROM cards WHERE id = $1", vars["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.flash(w, req, infoAlert, "Card was deleted.")
	http.Redirect(w, req, "/cards", http.StatusFound)
}

func (s *Server) memorize(w http.ResponseWriter, req *http.Request, cardType, cardId string) {
	var typ int
	switch cardType {
	case "general":
		typ = 1
	case "code":
		typ = 2
	default:
		http.Redirect(w, req, "/cards", http.StatusFound)
		return
	}

	var card Card
	if cardId != "" {
		s.db.Get(&card, "SELECT id, type, front, back, known FROM cards WHERE id = ? and type = ? LIMIT 1", cardId, typ)
	} else {
		s.db.Get(&card, "SELECT id, type, front, back, known FROM cards WHERE type = ? and known = 0 ORDER BY RANDOM() LIMIT 1", typ)
	}
	if card.ID == 0 {
		s.flash(w, req, infoAlert, "You've learned all the "+cardType+" cards.")
		http.Redirect(w, req, "/cards", http.StatusFound)
		return
	}

	shortAnswer := false
	if len(card.Back) < 75 {
		shortAnswer = true
	}

	s.render(w, req, "memorize.html", pongo2.Context{
		"card":         card,
		"card_type":    cardType,
		"short_answer": shortAnswer,
	})
}

func (s *Server) MarkKnown(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	if _, err := s.db.Exec("UPDATE cards SET known = 1 WHERE id = ? and type = ?",
		vars["id"],
		vars["type"],
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.flash(w, req, infoAlert, "Card marked as known.")
	s.memorize(w, req, vars["type"], vars["id"])
}

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		req.ParseForm()
		var (
			username = req.Form.Get("username")
			password = req.Form.Get("password")
		)
		if s.credentials != username+":"+password {
			s.flash(w, req, dangerAlert, "Invalid username/password")
		} else {
			s.auth(w, req, true)
			http.Redirect(w, req, "/cards", http.StatusFound)
			return
		}
	}
	s.render(w, req, "login.html", pongo2.Context{})
}

func (s *Server) Logout(w http.ResponseWriter, req *http.Request) {
	s.auth(w, req, false)
	s.flash(w, req, dangerAlert, "You've logged out")
	http.Redirect(w, req, "/", http.StatusFound)
}
