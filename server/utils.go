package server

import (
	"net/http"
	"path"

	"github.com/flosch/pongo2"
)

const (
	sessionName = "__session"

	// Bootstrap alert
	infoAlert    = "info"
	warningAlert = "warning"
	dangerAlert  = "danger"
)

type flash struct {
	Type string
	Body string
}

func (s *Server) flash(w http.ResponseWriter, req *http.Request, flashType, body string) {
	session, _ := s.sessions.Get(req, sessionName)
	defer session.Save(req, w)
	session.AddFlash(&flash{Type: flashType, Body: body})
}

func (s *Server) flashes(w http.ResponseWriter, req *http.Request) interface{} {
	session, _ := s.sessions.Get(req, sessionName)
	defer session.Save(req, w)
	return session.Flashes()
}

func (s *Server) render(w http.ResponseWriter, req *http.Request, templateName string, ctx pongo2.Context) {
	tmpl, err := pongo2.FromFile(path.Join(s.basePath, "templates", templateName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx["flashes"] = s.flashes(w, req)
	b, err := tmpl.ExecuteBytes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (s *Server) admin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// TODO: check credentials

		if req.Method == http.MethodPost {
			if err := req.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		next(w, req)
	}
}
