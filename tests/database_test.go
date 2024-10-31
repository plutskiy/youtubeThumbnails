package tests

import (
	"testTask/cmd/database"
	"testing"
)

func TestInitDB(t *testing.T) {
	db := database.InitDB()
	defer database.CloseDB()

	if db == nil {
		t.Fatal("expected non-nil database connection")
	}

	var name string
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='thumbnails';").Scan(&name)
	if err != nil {
		t.Fatalf("expected table thumbnails to exist, got error: %v", err)
	}
	if name != "thumbnails" {
		t.Fatalf("expected table name to be thumbnails, got: %s", name)
	}
}

func TestCloseDB(t *testing.T) {
	db := database.InitDB()
	database.CloseDB()

	err := db.Ping()
	if err == nil {
		t.Fatal("expected error when pinging closed database, got nil")
	}
}
