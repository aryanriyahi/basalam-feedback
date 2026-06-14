package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"basalam-feedback/internal/config"
	"basalam-feedback/internal/repository"
)

type Server struct {
	cfg  config.Config
	repo *repository.FeedbackRepository
	mux  *http.ServeMux
}

func New(cfg config.Config, db *sql.DB) *Server {
	s := &Server{
		cfg:  cfg,
		repo: repository.NewFeedbackRepository(db),
		mux:  http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/assets"))))
	s.mux.HandleFunc("/", s.serveIndex)
	s.mux.HandleFunc("/admin", s.serveAdmin)
	s.mux.HandleFunc("/api/feedbacks", s.feedbacks)
	s.mux.HandleFunc("/api/feedbacks/", s.feedbackStatus)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "ui/index.html")
}

func (s *Server) serveAdmin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/admin" {
		http.NotFound(w, r)
		return
	}

	s.requireBasicAuth(w, r, func() {
		http.ServeFile(w, r, "ui/dashboard.html")
	})
}
func (s *Server) feedbacks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.requireBasicAuth(w, r, func() {
			ctx, cancel := context.WithTimeout(r.Context(), 5_000_000_000)
			defer cancel()
			feedbacks, err := s.repo.List(ctx)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to load feedbacks")
				return
			}
			writeJSON(w, http.StatusOK, feedbacks)
		})
	case http.MethodPost:
		var payload struct {
			Title   string `json:"title"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON payload")
			return
		}
		if strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.Message) == "" {
			writeError(w, http.StatusBadRequest, "title and message are required")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5_000_000_000)
		defer cancel()
		if err := s.repo.Create(ctx, strings.TrimSpace(payload.Title), strings.TrimSpace(payload.Message)); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save feedback")
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) feedbackStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s.requireBasicAuth(w, r, func() {
		if !strings.HasSuffix(r.URL.Path, "/status") {
			http.NotFound(w, r)
			return
		}
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/feedbacks/"), "/status")
		if id == "" || strings.Contains(id, "/") {
			writeError(w, http.StatusBadRequest, "invalid feedback id")
			return
		}
		var payload struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON payload")
			return
		}
		if err := repository.ValidateStatus(payload.Status); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5_000_000_000)
		defer cancel()
		updated, err := s.repo.UpdateStatus(ctx, id, payload.Status)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to update status")
			return
		}
		if !updated {
			writeError(w, http.StatusNotFound, "feedback not found")
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func (s *Server) requireBasicAuth(w http.ResponseWriter, r *http.Request, next func()) {
	username, password, ok := r.BasicAuth()
	if !ok || username != s.cfg.BasicAuthUser || password != s.cfg.BasicAuthPassword {
		w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	next()
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
