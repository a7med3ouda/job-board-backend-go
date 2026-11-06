package job

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

func (r *Repository) GetAll() ([]model.Jobs, error) {
	stmt := postgres.SELECT(table.Jobs.AllColumns).FROM(table.Jobs).ORDER_BY(table.Jobs.CreatedAt.DESC())

	var jobs []model.Jobs
	err := stmt.Query(r.db, &jobs)
	return jobs, err
}

func (r *Repository) Create(job model.Jobs) (model.Jobs, error) {
	stmt := table.Jobs.INSERT(table.Jobs.Title, table.Jobs.CompanyName, table.Jobs.Location, table.Jobs.SalaryRange, table.Jobs.Description).
		MODEL(job).
		RETURNING(table.Jobs.AllColumns)

	var insertedJob model.Jobs
	err := stmt.Query(r.db, &insertedJob)
	return insertedJob, err
}