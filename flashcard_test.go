// flashcard_test.go
package main

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestFlashcardDB(t *testing.T) {
	// Open a database connection
	db, err := sql.Open("sqlite3", "./flashcards.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Create a new flashcard database
	flashcardDB := NewFlashcardDB(db)

	// Clear the flashcards table before running the tests
	err = flashcardDB.ClearFlashcards()
	if err != nil {
		t.Fatal(err)
	}

	// Test adding a flashcard
	flashcard := Flashcard{
		DeckID:   "1",
		Question: "What is the capital of France?",
		Answer:   "Paris",
	}
	err = flashcardDB.AddFlashcard(flashcard)
	if err != nil {
		t.Fatal(err)
	}

	// Test getting all flashcards
	flashcards, err := flashcardDB.GetFlashcards()
	if err != nil {
		t.Fatal(err)
	}

	// Check the number of flashcards
	expectedCount := 1
	if len(flashcards) != expectedCount {
		t.Errorf("expected %d flashcards, got %d", expectedCount, len(flashcards))
	}

	// Check the flashcard data
	if flashcards[0].DeckID != flashcard.DeckID {
		t.Errorf("expected deck ID %s, got %s", flashcard.DeckID, flashcards[0].DeckID)
	}
	if flashcards[0].Question != flashcard.Question {
		t.Errorf("expected question %s, got %s", flashcard.Question, flashcards[0].Question)
	}
	if flashcards[0].Answer != flashcard.Answer {
		t.Errorf("expected answer %s, got %s", flashcard.Answer, flashcards[0].Answer)
	}
}

func TestMain(m *testing.M) {
	// Set up any test-specific setup here
	// ...

	// Run the tests
	fmt.Println("Running tests...")
	exitCode := m.Run()

	// Perform any cleanup or teardown here
	// ...

	// Exit with the appropriate exit code
	fmt.Println("Tests finished.")
	fmt.Println("Exit code:", exitCode)
}

func NewFlashcardDB(db *sql.DB) *FlashcardDB {
	return &FlashcardDB{db: db}
}

type FlashcardDB struct {
	db *sql.DB
}

func (f *FlashcardDB) AddFlashcard(card Flashcard) error {
	// Implement the code to add a flashcard to the database
	// ...

	return nil
}

func (f *FlashcardDB) GetFlashcards() ([]Flashcard, error) {
	// Implement the code to retrieve all flashcards from the database
	// ...

	return nil, nil
}

func (f *FlashcardDB) ClearFlashcards() error {
	// Implement the code to clear all flashcards from the database
	// ...

	return nil
}

// ...
