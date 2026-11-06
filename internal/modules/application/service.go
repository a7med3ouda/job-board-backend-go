package application

import (
	"errors"
	"job-board/internal/database/gen/job-board/public/model"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ApplyForJob(app model.Applications) (model.Applications, error) {
	if app.ApplicantName == "" || app.Email == "" || app.ResumeURL == "" {
		return model.Applications{}, errors.New("applicantName, email, and resumeUrl are required")
	}

	return s.repo.Apply(app)
}

func (s *Service) GetApplicationsForJob(jobID int32) ([]model.Applications, error) {
	return s.repo.GetByJobID(jobID)
}