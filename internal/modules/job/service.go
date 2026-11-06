package job

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

func (s *Service) ListJobs() ([]model.Jobs, error) {
	return s.repo.GetAll()
}

func (s *Service) CreateJob(job model.Jobs) (model.Jobs, error) {
	if job.Title == "" || job.CompanyName == "" || job.Location == "" {
		return model.Jobs{}, errors.New("title, companyName, and location are required")
	}
	return s.repo.Create(job)
}