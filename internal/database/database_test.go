package database

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	path := "./database.test.json"
	db, err := NewDB(path)
	defer os.Remove(path)
	if err != nil {
		t.Fatalf("NewDB(%s) resulted in an error %v", path, err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("NewDB(%s) did not create a file", path)
	}

	if db.path != path {
		t.Errorf("DB path does not match, got: %s, want: %s", db.path, path)
	}

	if db.mux == nil {
		t.Errorf("DB mutex is nil")
	}
}
func TestLoadDB(t *testing.T) {
	path := "./database.test.json"
	db, _ := NewDB(path)
	defer os.Remove(path)

	data := `{"chirps": {"1": {"id": 1, "body": "test chirp"}}}`
	os.WriteFile(path, []byte(data), 0644)

	chirps, err := db.loadDB()
	if err != nil {
		t.Fatalf("loadDB resulted in an error: %v", err)
	}

	if len(chirps.Chirps) != 1 {
		t.Errorf("Expected 1 chirp, got %d", len(chirps.Chirps))
	}

	chirp, ok := chirps.Chirps[1]
	if !ok {
		t.Errorf("Chirp with ID 1 not found")
	}

	if chirp.Body != "test chirp" {
		t.Errorf("Expected chirp body to be 'test chirp', got '%s'", chirp.Body)
	}
}

func TestWriteDB(t *testing.T) {
	path := "./database.test.json"
	db, _ := NewDB(path)
	defer os.Remove(path)

	dbStructure := DBStructure{
		Chirps: map[int]Chirp{1: {Id: 1, Body: "test chirp"}},
	}

	err := db.writeDB(dbStructure)

	if err != nil {
		t.Fatalf("writeDB resulted in an error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("writeDB(%s) did not create a file", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Could not read file: %v", err)
	}

	expectedData := `{"chirps":{"1":{"id":1,"body":"test chirp"}}}`
	if string(data) != expectedData {
		t.Errorf("Expected data to be '%s', got '%s'", expectedData, string(data))
	}
}

func TestCreateChirp(t *testing.T) {
	path := "./database.test.json"
	db, _ := NewDB(path)
	defer os.Remove(path)

	// Create a chirp
	body := "test chirp"
	chirp, err := db.CreateChirp(body, 1)
	if err != nil {
		t.Fatalf("CreateChirp resulted in an error: %v", err)
	}

	if chirp.Body != body {
		t.Errorf("Expected chirp body to be '%s', got '%s'", body, chirp.Body)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Could not read file: %v", err)
	}

	expectedData := `{"chirps":{"1":{"id":1,"body":"test chirp"}}}`
	if string(data) != expectedData {
		t.Errorf("Expected data to be '%s', got '%s'", expectedData, string(data))
	}
}

func TestGetChirps(t *testing.T) {
	path := "./database.test.json"
	db, _ := NewDB(path)
	defer os.Remove(path)

	// Create a chirp
	body := "test chirp"
	_, err := db.CreateChirp(body, 1)
	if err != nil {
		t.Fatalf("CreateChirp resulted in an error: %v", err)
	}

	// Get the chirps
	chirps, err := db.GetChirps()
	if err != nil {
		t.Fatalf("GetChirps resulted in an error: %v", err)
	}

	// Check if the number of chirps is correct
	if len(chirps) != 1 {
		t.Errorf("Expected 1 chirp, got %d", len(chirps))
	}

	// Check if the chirp is correct
	if chirps[0].Body != body {
		t.Errorf("Expected chirp body to be '%s', got '%s'", body, chirps[0].Body)
	}
}
