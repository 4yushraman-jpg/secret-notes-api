package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"secret-notes-app/database"
	"secret-notes-app/middleware"
	"secret-notes-app/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email or password can't be empty!", http.StatusBadRequest)
		return
	}

	var exists string
	err = database.DB.QueryRow("SELECT email FROM users WHERE email = ?", req.Email).Scan(&exists)
	if err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	if err != sql.ErrNoRows {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO users (email, password) VALUES (?, ?)`
	_, err = database.DB.Exec(query, req.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Signup successful",
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req models.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email or password can't be empty!", http.StatusBadRequest)
		return
	}

	var id int
	var storedPassword string
	query := `SELECT id, password FROM users WHERE email = ?`
	err = database.DB.QueryRow(query, req.Email).Scan(&id, &storedPassword)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &models.Claims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.Note
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content can't be empty!", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO notes (content, user_id) VALUES (?, ?)`
	_, err = database.DB.Exec(query, req.Content, userID)
	if err != nil {
		http.Error(w, "Failed to create note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Note created",
	})
}

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `SELECT id, content FROM notes WHERE user_id = ?`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notes := []models.GetNotesResponse{}
	for rows.Next() {
		var n models.GetNotesResponse
		err := rows.Scan(&n.ID, &n.Content)
		if err != nil {
			http.Error(w, "Error fetching rows", http.StatusInternalServerError)
			return
		}
		notes = append(notes, n)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "Database iteration error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}
