package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polkaApiKey    string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	if secret == "" {
		log.Fatal("SECRET must be set")
	}

	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not connect to the database: %s", err)
		return
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         secret,
		polkaApiKey:    polkaKey,
	}

	// ==============================
	// Server Handlers
	// ==============================
	mux := http.NewServeMux()

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))

	mux.Handle("GET /admin/metrics", apiCfg.writeNumberOfRequests())
	mux.Handle("POST /admin/reset", apiCfg.resetHitCount())

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleWebhookEvent)

	mux.Handle("POST /api/chirps", apiCfg.handleChirpCreation())
	mux.HandleFunc("GET /api/chirps", apiCfg.handleChirpRetrieval)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleIndividualChirpRetrieval)
	mux.Handle("POST /api/users", apiCfg.createUser())
	mux.HandleFunc("POST /api/login", apiCfg.loginUser)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUpdateCredentials)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleChirpDeletion)

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Fatal("An error occurred while writing the body")
		}
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) writeNumberOfRequests() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")

		w.Write([]byte(fmt.Sprintf(`
		<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())))
	})
}

func (cfg *apiConfig) resetHitCount() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.platform == "dev" {
			cfg.db.DeleteAllUsers(r.Context())
		} else {
			respondWithError(
				w,
				http.StatusForbidden,
				"Cannot delete users while in non development mode!!!",
				nil)
		}
		cfg.fileserverHits.Store(0)
	})
}
