// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

// Flashcard represents a flashcard with a question and an answer.
type Flashcard struct {
	ID       int    `json:"id"`
	DeckID   string `json:"deckId"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

var db *sql.DB

func main() {
	// Open a database connection
	var err error
	db, err = sql.Open("sqlite3", "./flashcards.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the flashcards table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS flashcards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		deckId TEXT NOT NULL,
		question TEXT NOT NULL,
		answer TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Set up the HTTP server
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/api/flashcards", getFlashcardsHandler).Methods("GET")
	router.HandleFunc("/api/flashcards", addFlashcardHandler).Methods("POST")
	router.HandleFunc("/api/flashcards/{id}", updateFlashcardHandler).Methods("PUT")
	router.HandleFunc("/api/flashcards/{id}", deleteFlashcardHandler).Methods("DELETE")

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handler for the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// Handler for getting all flashcards
func getFlashcardsHandler(w http.ResponseWriter, r *http.Request) {
	flashcards, err := getFlashcards()
	if err != nil {
		http.Error(w, "Failed to get flashcards", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flashcards)
}

// Handler for adding a flashcard
func addFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	var flashcard Flashcard
	err := json.NewDecoder(r.Body).Decode(&flashcard)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = addFlashcard(flashcard)
	if err != nil {
		http.Error(w, "Failed to add flashcard", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Handler for updating a flashcard
func updateFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	flashcardID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid flashcard ID", http.StatusBadRequest)
		return
	}

	// Retrieve the existing flashcard from the database
	existingFlashcard, err := getFlashcardByID(flashcardID)
	if err != nil {
		http.Error(w, "Failed to retrieve flashcard", http.StatusInternalServerError)
		return
	}

	// Decode the request body into a temporary flashcard struct
	var tempFlashcard struct {
		DeckID   string `json:"deckId"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
	err = json.NewDecoder(r.Body).Decode(&tempFlashcard)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the flashcard with the new values while preserving the original deckId
	flashcard := Flashcard{
		ID:       flashcardID,
		DeckID:   existingFlashcard.DeckID,
		Question: tempFlashcard.Question,
		Answer:   tempFlashcard.Answer,
	}

	err = updateFlashcard(flashcardID, flashcard)
	if err != nil {
		http.Error(w, "Failed to update flashcard", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Function to retrieve a flashcard by ID from the database
func getFlashcardByID(id int) (*Flashcard, error) {
	row := db.QueryRow("SELECT id, deckId, question, answer FROM flashcards WHERE id = ?", id)

	var flashcard Flashcard
	err := row.Scan(&flashcard.ID, &flashcard.DeckID, &flashcard.Question, &flashcard.Answer)
	if err != nil {
		return nil, err
	}

	return &flashcard, nil
}

// Handler for deleting a flashcard
func deleteFlashcardHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	flashcardID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid flashcard ID", http.StatusBadRequest)
		return
	}

	err = deleteFlashcard(flashcardID)
	if err != nil {
		http.Error(w, "Failed to delete flashcard", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Function to retrieve all flashcards from the database
func getFlashcards() ([]Flashcard, error) {
	rows, err := db.Query("SELECT id, deckId, question, answer FROM flashcards")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flashcards []Flashcard

	for rows.Next() {
		var flashcard Flashcard
		err := rows.Scan(&flashcard.ID, &flashcard.DeckID, &flashcard.Question, &flashcard.Answer)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}

	return flashcards, nil
}

// Function to add a flashcard to the database
func addFlashcard(flashcard Flashcard) error {
	_, err := db.Exec("INSERT INTO flashcards (deckId, question, answer) VALUES (?, ?, ?)",
		flashcard.DeckID, flashcard.Question, flashcard.Answer)
	return err
}

// Function to update a flashcard in the database
func updateFlashcard(flashcardID int, flashcard Flashcard) error {
	_, err := db.Exec("UPDATE flashcards SET deckId = ?, question = ?, answer = ? WHERE id = ?",
		flashcard.DeckID, flashcard.Question, flashcard.Answer, flashcardID)
	return err
}

// Function to delete a flashcard from the database
func deleteFlashcard(flashcardID int) error {
	_, err := db.Exec("DELETE FROM flashcards WHERE id = ?", flashcardID)
	return err
}
