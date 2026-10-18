package handlers

import (
	"encoding/json"
	"net/http"

	"cryptography/internal/rsa"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type generateRequest struct {
	Bits int `json:"bits"`
}

type encryptRequest struct {
	Message string `json:"message"`
	N       string `json:"n"`
	E       string `json:"e"`
}

type decryptRequest struct {
	Ciphertext string `json:"ciphertext"`
	N          string `json:"n"`
	D          string `json:"d"`
}

type attackRequest struct {
	N          string `json:"n"`
	E          string `json:"e"`
	Ciphertext string `json:"ciphertext"`
	Method     string `json:"method"` // factor | fermat | small_e | wiener
}

// RSAGenerate — POST /api/rsa/generate
func RSAGenerate(w http.ResponseWriter, r *http.Request) {
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// body может быть пустым — тогда bits по умолчанию
		req.Bits = 64
	}
	if req.Bits == 0 {
		req.Bits = 64
	}

	result, err := rsa.Generate(req.Bits)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// RSAEncrypt — POST /api/rsa/encrypt
func RSAEncrypt(w http.ResponseWriter, r *http.Request) {
	var req encryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный JSON")
		return
	}
	if req.Message == "" || req.N == "" || req.E == "" {
		writeError(w, http.StatusBadRequest, "нужны поля message, n, e")
		return
	}

	result, err := rsa.Encrypt(req.Message, req.N, req.E)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// RSADecrypt — POST /api/rsa/decrypt
func RSADecrypt(w http.ResponseWriter, r *http.Request) {
	var req decryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный JSON")
		return
	}
	if req.Ciphertext == "" || req.N == "" || req.D == "" {
		writeError(w, http.StatusBadRequest, "нужны поля ciphertext, n, d")
		return
	}

	result, err := rsa.Decrypt(req.Ciphertext, req.N, req.D)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// RSAAttack — POST /api/rsa/attack (пока заглушка)
func RSAAttack(w http.ResponseWriter, r *http.Request) {
	var req attackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "неверный JSON")
		return
	}

	var (
		result *rsa.AttackResult
		err    error
	)

	switch req.Method {
	case "factor":
		result, err = rsa.AttackFactor(req.N, req.E, req.Ciphertext)
	case "fermat":
		result, err = rsa.AttackFermat(req.N, req.E, req.Ciphertext)
	case "small_e":
		result, err = rsa.AttackSmallE(req.N, req.E, req.Ciphertext)
	case "wiener":
		result, err = rsa.AttackWiener(req.N, req.E, req.Ciphertext)
	default:
		writeError(w, http.StatusBadRequest, "method: factor | fermat | small_e | wiener")
		return
	}

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
