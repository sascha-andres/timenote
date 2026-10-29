package main

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"go.livingit.de/timenote"
	"go.livingit.de/timenote/remote"
)

const maxBodyBytes = 1 << 20

type server struct {
	// mu serializes backend calls; the toggl backend and its file cache are not safe for concurrent use
	mu      sync.Mutex
	backend timenote.Backend
	auth    []byte
}

// newHandler exposes backend as the /v1 REST API. Each request must send
// "Authorization: Bearer <token>".
func newHandler(backend timenote.Backend, token string) http.Handler {
	s := &server{backend: backend, auth: []byte("Bearer " + token)}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/timeentries", s.newTimeEntry)
	mux.HandleFunc("GET /v1/timeentries", s.listForDay)
	mux.HandleFunc("GET /v1/timeentries/current", s.current)
	mux.HandleFunc("PATCH /v1/timeentries/current", s.updateCurrent)
	mux.HandleFunc("PUT /v1/timeentries/current/description", s.setDescription)
	mux.HandleFunc("PATCH /v1/timeentries/current/description", s.appendDescription)
	mux.HandleFunc("PATCH /v1/timeentries/current/tags", s.toggleTag)
	mux.HandleFunc("PUT /v1/timeentries/current/project", s.setProject)
	mux.HandleFunc("GET /v1/projects", s.projects)
	mux.HandleFunc("POST /v1/projects", s.createProject)
	mux.HandleFunc("DELETE /v1/projects/{nameOrID}", s.deleteProject)
	mux.HandleFunc("GET /v1/clients", s.clients)
	mux.HandleFunc("POST /v1/clients", s.createClient)
	return s.authenticate(mux)
}

func (s *server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), s.auth) != 1 {
			writeProblem(w, http.StatusUnauthorized, "Unauthorized", "missing or invalid bearer token", "UNAUTHORIZED")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeProblem(w http.ResponseWriter, status int, title, detail, code string) {
	w.Header().Set("Content-Type", remote.ProblemContentType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(remote.Problem{Type: "about:blank", Title: title, Status: status, Detail: detail, Code: code})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// decode reads the JSON body into v. On failure it writes a problem and returns false.
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeProblem(w, http.StatusBadRequest, "Invalid request body", err.Error(), "INVALID_BODY")
		return false
	}
	return true
}

// requireName writes a validation problem and returns false if value is empty.
func requireName(w http.ResponseWriter, field, value string) bool {
	if value != "" {
		return true
	}
	w.Header().Set("Content-Type", remote.ProblemContentType)
	w.WriteHeader(http.StatusUnprocessableEntity)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"type":       "about:blank",
		"title":      "Validation failed",
		"status":     http.StatusUnprocessableEntity,
		"detail":     "The request has 1 invalid field",
		"code":       "VALIDATION_FAILED",
		"validation": map[string]any{field: map[string]any{"value": value, "error": "must not be empty"}},
	})
	return false
}

// call runs fn under the backend lock and writes the result: a problem on
// error, 204 if out is nil, else out as JSON.
func (s *server) call(w http.ResponseWriter, fn func() (any, error)) {
	s.mu.Lock()
	out, err := fn()
	s.mu.Unlock()
	switch {
	case errors.Is(err, timenote.ErrNoCurrentEntry):
		writeProblem(w, http.StatusNotFound, "No current time entry", err.Error(), "NO_CURRENT_ENTRY")
		return
	case errors.Is(err, timenote.ErrProjectNotFound):
		writeProblem(w, http.StatusNotFound, "Project not found", err.Error(), "PROJECT_NOT_FOUND")
		return
	case err != nil:
		slog.Error("backend call failed", "error", err)
		writeProblem(w, http.StatusBadGateway, "Backend call failed", err.Error(), "BACKEND_ERROR")
		return
	}
	if out == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	writeJSON(w, out)
}

// noResult adapts a backend method without result to call.
func noResult(fn func() error) func() (any, error) {
	return func() (any, error) { return nil, fn() }
}

func (s *server) newTimeEntry(w http.ResponseWriter, r *http.Request) {
	var req remote.NewTimeEntryRequest
	if !decode(w, r, &req) {
		return
	}
	if req.FromPrevious {
		s.call(w, noResult(s.backend.StartPrevious))
		return
	}
	s.call(w, noResult(s.backend.New))
}

func (s *server) listForDay(w http.ResponseWriter, _ *http.Request) {
	s.call(w, func() (any, error) {
		entries, err := s.backend.ListForDay()
		if entries == nil {
			entries = []timenote.TimeEntry{}
		}
		return entries, err
	})
}

func (s *server) current(w http.ResponseWriter, _ *http.Request) {
	s.call(w, func() (any, error) {
		te, err := s.backend.Current()
		if err == nil && te == nil {
			err = timenote.ErrNoCurrentEntry
		}
		return te, err
	})
}

func (s *server) updateCurrent(w http.ResponseWriter, r *http.Request) {
	var req remote.UpdateCurrentRequest
	if !decode(w, r, &req) {
		return
	}
	if !req.Stop {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	s.call(w, noResult(s.backend.Done))
}

func (s *server) setDescription(w http.ResponseWriter, r *http.Request) {
	var req remote.SetDescriptionRequest
	if !decode(w, r, &req) {
		return
	}
	s.call(w, noResult(func() error { return s.backend.SetDescription(req.Description) }))
}

func (s *server) appendDescription(w http.ResponseWriter, r *http.Request) {
	var req remote.AppendDescriptionRequest
	if !decode(w, r, &req) {
		return
	}
	s.call(w, noResult(func() error { return s.backend.Append(req.Line, req.Separator) }))
}

func (s *server) toggleTag(w http.ResponseWriter, r *http.Request) {
	var req remote.ToggleTagRequest
	if !decode(w, r, &req) || !requireName(w, "toggle", req.Toggle) {
		return
	}
	s.call(w, noResult(func() error { return s.backend.Tag(req.Toggle) }))
}

func (s *server) setProject(w http.ResponseWriter, r *http.Request) {
	var req remote.SetProjectRequest
	if !decode(w, r, &req) {
		return
	}
	s.call(w, noResult(func() error { return s.backend.SetProjectForCurrentTimestamp(req.Name, req.AutoCreateProject) }))
}

func (s *server) projects(w http.ResponseWriter, _ *http.Request) {
	s.call(w, func() (any, error) {
		projects, err := s.backend.Projects()
		if projects == nil {
			projects = []timenote.Project{}
		}
		return projects, err
	})
}

func (s *server) createProject(w http.ResponseWriter, r *http.Request) {
	var req remote.NameRequest
	if !decode(w, r, &req) || !requireName(w, "name", req.Name) {
		return
	}
	s.call(w, noResult(func() error { return s.backend.CreateProject(req.Name) }))
}

func (s *server) deleteProject(w http.ResponseWriter, r *http.Request) {
	s.call(w, noResult(func() error { return s.backend.DeleteProject(r.PathValue("nameOrID")) }))
}

func (s *server) clients(w http.ResponseWriter, _ *http.Request) {
	s.call(w, func() (any, error) {
		clients, err := s.backend.Clients()
		if clients == nil {
			clients = []timenote.Client{}
		}
		return clients, err
	})
}

func (s *server) createClient(w http.ResponseWriter, r *http.Request) {
	var req remote.NameRequest
	if !decode(w, r, &req) || !requireName(w, "name", req.Name) {
		return
	}
	s.call(w, noResult(func() error { return s.backend.CreateClient(req.Name) }))
}
