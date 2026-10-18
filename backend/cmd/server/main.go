package main

import (
	"log"
	"net/http"

	"cryptography/internal/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	// RSA
	mux.HandleFunc("POST /api/rsa/generate", handlers.RSAGenerate)
	mux.HandleFunc("POST /api/rsa/encrypt", handlers.RSAEncrypt)
	mux.HandleFunc("POST /api/rsa/decrypt", handlers.RSADecrypt)
	mux.HandleFunc("POST /api/rsa/attack", handlers.RSAAttack)

	// Healthcheck
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	addr := ":8080"
	log.Printf("Crypto Lab backend listening on %s", addr)
	if err := http.ListenAndServe(addr, corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
