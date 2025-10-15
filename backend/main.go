package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"notes-app/db"
	"notes-app/handlers"
)

func main() {
	// Database connection and migration would typically go here

	// Подключаемся к базе и создаём таблицу
	database, err := db.Connect("./notes.db")
	if err != nil {
		log.Fatal("db connect:", err)
	}
	if err := db.Migrate(database); err != nil {
		log.Fatal("db migrate:", err)
	}

	r := chi.NewRouter()

	// Set up CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for simplicity
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Define a simple health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/notes", handlers.GetNotes(database))
	r.Post("/notes", handlers.CreateNote(database))
	r.Put("/notes/{id}", handlers.UpdateNote(database))
	r.Delete("/notes/{id}", handlers.DeleteNote(database))

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
