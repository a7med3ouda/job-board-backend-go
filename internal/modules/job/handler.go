package job

import (
	"encoding/json"
	"job-board/internal/database/gen/job-board/public/model"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /jobs", h.handleListJobs)
	mux.HandleFunc("POST /jobs", h.handleCreateJob)
}

func (h *Handler) handleListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.service.ListJobs()
	if err != nil {
		http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *Handler) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req model.Jobs
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	createdJob, err := h.service.CreateJob(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdJob)
}