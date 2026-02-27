package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/joliverstrom-cmd/goServe/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't open database: %v", err)
	}
	dbQueries := database.New(db)

	serveMux := http.NewServeMux()
	var root http.Dir

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	jwtSecret := os.Getenv("JWTSECRET")
	if jwtSecret == "" {
		log.Fatal("JWTSECRET must be set")
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	serveMux.HandleFunc("GET /api/healthz", readinessHandler)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.createPost)
	serveMux.HandleFunc("POST /api/users", apiCfg.createUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.login)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refreshCheck)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revokeRefreshToken)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.getPosts)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getPost)

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(root))))

	serveMux.HandleFunc("GET /admin/metrics", apiCfg.countHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
