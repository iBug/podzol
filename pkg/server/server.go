package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/docker"
)

type Server struct {
	docker *docker.Client
	mux    *http.ServeMux

	listenAddr string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewServer(v *viper.Viper) (*Server, error) {
	dockerClient, err := docker.NewClient(v)
	if err != nil {
		return nil, err
	}

	return &Server{
		docker: dockerClient,
		mux:    http.NewServeMux(),

		listenAddr: v.GetString("listen-addr"),
	}, nil
}

func HandleDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "not found"})
}

// Create a container.
func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var opts docker.ContainerOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	ctx := r.Context()
	info, err := s.docker.Create(ctx, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s := fmt.Sprintf("failed to create container: %v", err)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: s})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(info)
}

// Remove a container.
func (s *Server) HandleRemove(w http.ResponseWriter, r *http.Request) {
	var opts docker.ContainerOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	ctx := r.Context()
	if err := s.docker.Remove(ctx, opts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s := fmt.Sprintf("failed to remove container: %v", err)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: s})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

// List containers.
// Filters of type docker.ContainerOptions may be passed as either the "opts" query parameter or as request body. In either case, the filters are JSON-encoded.
func (s *Server) HandleList(w http.ResponseWriter, r *http.Request) {
	var opts docker.ContainerOptions
	switch r.Method {
	case http.MethodGet:
		optsJSON := r.URL.Query().Get("opts")
		if optsJSON != "" {
			if err := json.Unmarshal([]byte(optsJSON), &opts); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
				return
			}
		}
	case http.MethodPost:
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	containers, err := s.docker.List(ctx, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s := fmt.Sprintf("failed to list containers: %v", err)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: s})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(containers)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	s.mux.HandleFunc("/", HandleDefault)
	s.mux.HandleFunc("/create", s.HandleCreate)
	s.mux.HandleFunc("/remove", s.HandleRemove)
	s.mux.HandleFunc("/list", s.HandleList)
	return http.ListenAndServe(s.listenAddr, s)
}
