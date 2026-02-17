package main

import (
	"log"
	"net/http"
	"secret-notes-app/database"
	"secret-notes-app/handlers"
	"secret-notes-app/middleware"
)

func main() {
	database.InitDB()

	http.HandleFunc("/signup", handlers.SignupHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	http.HandleFunc("POST /notes", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.CreateNoteHandler)))
	http.HandleFunc("GET /notes", middleware.RateLimitMiddleware(middleware.AuthMiddleware(handlers.GetNotesHandler)))

	log.Println("Server is running on port 8080..")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
