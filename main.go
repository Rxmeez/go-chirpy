package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rxmeez/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.DB
	jwtSecret      string
	polkaSecret    string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		w.Header().Add("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func main() {
	const port string = "8080"
	const filepathRoot string = "./app"
	godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	polkaSecret := os.Getenv("POLKA_SECRET")
	if polkaSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := database.NewDB("./database.json")
	if err != nil {
		log.Fatal(err)
	}
	apiCfg := apiConfig{fileserverHits: 0, db: db, jwtSecret: jwtSecret, polkaSecret: polkaSecret}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("/api/reset", apiCfg.handlerReset)

	mux.HandleFunc("POST /api/login", apiCfg.handlerUsersLogin)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps/", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerChirpRetrieveId)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.handlerChirpDeleteId)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUsersUpgrade)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
