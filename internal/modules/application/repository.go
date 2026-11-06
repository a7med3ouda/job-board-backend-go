package application

import (
	"database/sql"
	"job-board/internal/database/gen/job-board/public/model"
	"job-board/internal/database/gen/job-board/public/table"

	"github.com/go-jet/jet/v2/postgres"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Apply(app model.Applications) (model.Applications, error) {
	stmt := table.Applications.INSERT(table.Applications.JobID, table.Applications.ApplicantName, table.Applications.Email, table.Applications.ResumeURL).
		MODEL(app).
		RETURNING(table.Applications.AllColumns)

	var inserted model.Applications
	err := stmt.Query(r.db, &inserted)
	return inserted, err
}

func (r *Repository) GetByJobID(jobID int32) ([]model.Applications, error) {
	stmt := postgres.SELECT(table.Applications.AllColumns).
		FROM(table.Applications).
		WHERE(table.Applications.JobID.EQ(postgres.Int(int64(jobID)))).
		ORDER_BY(table.Applications.CreatedAt.DESC())

	var apps []model.Applications
	err := stmt.Query(r.db, &apps)
	return apps, err
}