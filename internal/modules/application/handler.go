package application

import (
	"encoding/json"
	"job-board/internal/database/gen/job-board/public/model"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /jobs/{id}/apply", h.handleApply)
	mux.HandleFunc("GET /jobs/{id}/applications", h.handleGetApplications)
}

func (h *Handler) handleApply(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	jobID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Job ID", http.StatusBadRequest)
		return
	}

	var app model.Applications
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	app.JobID = int32(jobID)

	inserted, err := h.service.ApplyForJob(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inserted)
}

func (h *Handler) handleGetApplications(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	jobID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Job ID", http.StatusBadRequest)
		return
	}

	apps, err := h.service.GetApplicationsForJob(int32(jobID))
	if err != nil {
		http.Error(w, "Failed to fetch applications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}