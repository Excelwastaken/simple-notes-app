package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"notes-app/models"
)

func GetNotes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, title, content, created_at, updated_at FROM notes")
		if err != nil {
			http.Error(w, "Failed to fetch notes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var notes []models.Note
		for rows.Next() {
			var note models.Note
			if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt); err != nil {
				http.Error(w, "Failed to scan note", http.StatusInternalServerError)
				return
			}
			notes = append(notes, note)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(notes)
	}
}

func CreateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type payload struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		var p payload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if p.Title == "" || p.Content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", p.Title, p.Content)
		if err != nil {
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()

		var note models.Note
		err = db.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?", id).
			Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to fetch created note", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

func UpdateNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/notes/"):]
		type payload struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		var p payload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if p.Title == "" || p.Content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		result, err := db.Exec("UPDATE notes SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", p.Title, p.Content, idStr)
		if err != nil {
			http.Error(w, "Failed to update note", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		var note models.Note
		err = db.QueryRow("SELECT id, title, content, created_at, updated_at FROM notes WHERE id = ?", idStr).
			Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to fetch updated note", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(note)

	}
}

func DeleteNote(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Path[len("/notes/"):]

		result, err := db.Exec("DELETE FROM notes WHERE id = ?", idStr)
		if err != nil {
			http.Error(w, "Failed to delete note", http.StatusInternalServerError)
			return
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			http.Error(w, "Note not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
