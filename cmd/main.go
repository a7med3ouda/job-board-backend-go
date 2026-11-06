package main

import (
	"encoding/json"
	"job-board/internal/config"
	"job-board/internal/database"
	"job-board/internal/modules/application"
	"job-board/internal/modules/job"
	"log"
	"net/http"
)

func main() {
	// load environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration failed: %v", err)
	}
	
	// connect database
	db, err := database.NewPostgresConnection(cfg.DbConnectionString)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	// close db connection at the end of the function
	defer db.Close()

	// initialize server
	mux := http.NewServeMux()

	// register health route
	mux.HandleFunc("GET /health", handleHealthCheck)

	// register job routes
	jobRepo := job.NewRepository(db)
	jobService := job.NewService(jobRepo)
	jobHandler := job.NewHandler(jobService)
	jobHandler.RegisterRoutes(mux)

	// register application routes
	appRepo := application.NewRepository(db)
	appService := application.NewService(appRepo)
	appHandler := application.NewHandler(appService)
	appHandler.RegisterRoutes(mux)

	// run server
	log.Printf("Server is running on port %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatal(err)
	}
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}