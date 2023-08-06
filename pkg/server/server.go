package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
	"github.com/ustclug/podzol/pkg/docker"
)

type Server struct {
	dockerClient *docker.Client
	mux          *http.ServeMux
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
		dockerClient: dockerClient,
		mux:          http.NewServeMux(),
	}, nil
}

func HandleDefault(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: "not found"})
}

type CreateResponse struct {
	ID string `json:"id"`
}

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
	id, err := s.dockerClient.Create(ctx, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s := fmt.Sprintf("failed to create container: %v", err)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: s})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(CreateResponse{ID: id})
}

func (s *Server) HandleRemove(w http.ResponseWriter, r *http.Request) {
	var opts docker.ContainerOptions
	if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	ctx := r.Context()
	if err := s.dockerClient.Remove(ctx, opts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s := fmt.Sprintf("failed to remove container: %v", err)
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: s})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (s *Server) HandleList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	optsJSON := r.URL.Query().Get("opts")
	var opts docker.ContainerOptions
	if optsJSON != "" {
		if err := json.Unmarshal([]byte(optsJSON), &opts); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
	}

	ctx := r.Context()
	containers, err := s.dockerClient.List(ctx, opts)
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
	defer r.Body.Close()
	s.mux.ServeHTTP(w, r)
}

func (s *Server) Run() error {
	s.mux.HandleFunc("/", HandleDefault)
	s.mux.HandleFunc("/create", s.HandleCreate)
	return http.ListenAndServe(":8080", s)
}
