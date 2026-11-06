# Job Board API

A small REST API for posting job openings and receiving applications, written in Go with a clean
handler → service → repository layering and type-safe SQL via [go-jet](https://github.com/go-jet/jet).

## Features

- List and create job postings
- Apply to a job posting with name, email and resume URL
- List all applications for a given job
- Duplicate protection: one application per email per job (DB unique constraint)
- Health check endpoint
- Environment-based configuration (`.env`)
- Type-safe SQL generation with `go-jet` (generated models + table builders)

## Tech Stack

| Concern     | Choice                                  |
| ----------- | --------------------------------------- |
| Language    | Go (see `go.mod`, `go 1.27.2`)          |
| Router      | `net/http` `ServeMux` (method patterns) |
| Database    | PostgreSQL                              |
| DB driver   | `github.com/lib/pq`                     |
| SQL builder | `github.com/go-jet/jet/v2`              |
| Config      | `github.com/joho/godotenv`              |

## Project Structure

```text
.
├── cmd/
│   └── main.go # entrypoint: config, DB, routes, server
├── internal/
│   ├── config/
│   │   └── environment.go # loads .env, validates required vars
│   ├── database/
│   │   ├── database.go # Postgres connection + ping
│   │   ├── schema.sql # tables, constraints, indexes
│   │   └── gen/job-board/public/
│   │       ├── model/ # go-jet generated structs
│   │       └── table/ # go-jet generated table builders
│   └── modules/
│       ├── job/
│       │   ├── handler.go # GET /jobs, POST /jobs
│       │   ├── service.go # validation / business rules
│       │   └── repository.go # SQL via go-jet
│       └── application/
│           ├── handler.go # POST /jobs/{id}/apply, GET /jobs/{id}/applications
│           ├── service.go # validation
│           └── repository.go # SQL via go-jet
├── example.env
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.27+ (as declared in `go.mod`)
- A running PostgreSQL instance

## Getting Started

### 1. Clone and install dependencies

```bash
git clone <repo-url>
cd job-board
go mod download
```

### 2. Configure environment

```bash
cp example.env .env
```

Then fill in `.env`:

```env
PORT=8080
DB_CONNECTION_STRING=postgres://user:password@localhost:5432/job-board?sslmode=disable
```

| Variable               | Required | Default | Description                    |
| ---------------------- | -------- | ------- | ------------------------------ |
| `PORT`                 | No       | `8080`  | HTTP port the server binds to  |
| `DB_CONNECTION_STRING` | **Yes**  | –       | Postgres DSN (`lib/pq` format) |

### 3. Create the database schema

```bash
psql "$DB_CONNECTION_STRING" -f internal/database/schema.sql
```

This creates the `jobs` and `applications` tables, a `UNIQUE (job_id, email)` constraint,
and indexes on `jobs(location)` and `applications(job_id)`.

### 4. Run the server

```bash
go run ./cmd/main.go
```

```text
Server is running on port 8080
```

## API Reference

Base URL: `http://localhost:8080`

### `GET /health`

Health probe.

```bash
curl http://localhost:8080/health
```

```json
{ "status": "ok" }
```

---

### `GET /jobs`

Returns all job postings, newest first.

```bash
curl http://localhost:8080/jobs
```

```json
[
  {
    "ID": 1,
    "Title": "Backend Engineer",
    "CompanyName": "Acme",
    "Location": "Cairo",
    "SalaryRange": "$3,000 - $5,000",
    "Description": "Build and maintain Go services.",
    "CreatedAt": "2025-01-01T10:00:00Z"
  }
]
```

---

### `POST /jobs`

Creates a new job posting.

Required fields: `title`, `companyName`, `location`, `description`.

```bash
curl -X POST http://localhost:8080/jobs \
 -H "Content-Type: application/json" \
 -d '{
"title": "Backend Engineer",
"companyName": "Acme",
"location": "Cairo",
"salaryRange": "$3,000 - $5,000",
"description": "Build and maintain Go services."
}'
```

**201 Created**

```json
{
  "ID": 1,
  "Title": "Backend Engineer",
  "CompanyName": "Acme",
  "Location": "Cairo",
  "SalaryRange": "$3,000 - $5,000",
  "Description": "Build and maintain Go services.",
  "CreatedAt": "2025-01-01T10:00:00Z"
}
```

**400 Bad Request** — invalid JSON, or missing `title`, `companyName` or `location`.

---

### `POST /jobs/{id}/apply`

Submits an application for the job with the given id.
The `jobID` from the path overrides any body value.

Required body fields: `applicantName`, `email`, `resumeUrl`.

```bash
curl -X POST http://localhost:8080/jobs/1/apply \
 -H "Content-Type: application/json" \
 -d '{
"applicantName": "Jane Doe",
"email": "jane@example.com",
"resumeUrl": "https://example.com/cv/jane.pdf"
}'
```

**201 Created**

```json
{
  "ID": 1,
  "JobID": 1,
  "ApplicantName": "Jane Doe",
  "Email": "jane@example.com",
  "ResumeURL": "https://example.com/cv/jane.pdf",
  "CreatedAt": "2025-01-01T10:05:00Z"
}
```

**400 Bad Request** — invalid job id, invalid JSON, missing required fields, or a duplicate
application for the same `(job_id, email)` pair.

---

### `GET /jobs/{id}/applications`

Lists applications for a job, newest first.

```bash
curl http://localhost:8080/jobs/1/applications
```

```json
[
  {
    "ID": 1,
    "JobID": 1,
    "ApplicantName": "Jane Doe",
    "Email": "jane@example.com",
    "ResumeURL": "https://example.com/cv/jane.pdf",
    "CreatedAt": "2025-01-01T10:05:00Z"
  }
]
```

## Error Responses

Errors are returned as plain text with a matching HTTP status code:

| Status | Meaning                                               |
| ------ | ----------------------------------------------------- |
| 400    | Invalid path id, malformed JSON, or failed validation |
| 500    | Database or unexpected server failure                 |

Example: `applicantName, email, and resumeUrl are required`

## Data Model

```text
jobs
├── id SERIAL PRIMARY KEY
├── title VARCHAR(255) NOT NULL
├── company_name VARCHAR(255) NOT NULL
├── location VARCHAR(100) NOT NULL
├── salary_range VARCHAR(100)
├── description TEXT NOT NULL
└── created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP

applications
├── id SERIAL PRIMARY KEY
├── job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE
├── applicant_name VARCHAR(255) NOT NULL
├── email VARCHAR(255) NOT NULL
├── resume_url TEXT NOT NULL
├── created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
└── UNIQUE (job_id, email)
```

## Regenerating go-jet Code

The files under `internal/database/gen/` are generated — do not edit them by hand.
After changing `schema.sql`, regenerate with:

```bash
go install github.com/go-jet/jet/v2/cmd/jet@latest

jet -dsn="$DB_CONNECTION_STRING" \
 -schema=public \
 -path=./internal/database/gen
```

## Notes & Possible Improvements

- JSON request bodies are decoded into the generated `model.*` structs, so camelCase keys
  (`companyName`, `applicantName`, `resumeUrl`) are expected — snake_case keys are ignored.
- There is no authentication/authorization yet; `POST /jobs` is open to anyone.
- No pagination or filtering on `GET /jobs` (the `location` index is ready for it).
- Add middleware for request logging, panic recovery, and CORS.
- Return structured JSON error bodies instead of `http.Error` plain text.
