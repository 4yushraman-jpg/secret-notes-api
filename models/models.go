package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Note struct {
	Content string `json:"content"`
}

type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

type GetNotesResponse struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}
