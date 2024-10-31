package server

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testTask/pkg/api"
	"testing"
)

func TestGetThumbnail_InvalidURL(t *testing.T) {
	db := &sql.DB{}
	s := NewThumbnailServer(db)

	req := &api.ThumbnailRequest{Url: "https://example.com/video"}
	_, err := s.GetThumbnail(context.Background(), req)
	if err == nil || err.Error() != "invalid url" {
		t.Errorf("expected error 'invalid url', got %v", err)
	}
}

func TestGetThumbnail_CacheHit(t *testing.T) {
	db := &sql.DB{}
	s := NewThumbnailServer(db)
	testURL := "https://youtube.com/watch?v=test123"
	testImage := []byte("test image")

	s.cache[testURL] = testImage

	req := &api.ThumbnailRequest{Url: testURL}
	resp, err := s.GetThumbnail(context.Background(), req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if string(resp.Image) != string(testImage) {
		t.Errorf("expected image %v, got %v", testImage, resp.Image)
	}
}

func TestGetThumbnailFromDB(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("failed to close database: %v", err)
		}
	}(db)

	_, err = db.Exec("CREATE TABLE thumbnails (url TEXT PRIMARY KEY, thumbnail BLOB)")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	testURL := "https://youtube.com/watch?v=test123"
	testThumbnail := []byte("test thumbnail")
	_, err = db.Exec("INSERT INTO thumbnails (url, thumbnail) VALUES (?, ?)", testURL, testThumbnail)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	s := NewThumbnailServer(db)

	thumbnail, err := s.getThumbnailFromDB(testURL)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if string(thumbnail) != string(testThumbnail) {
		t.Errorf("expected thumbnail %v, got %v", testThumbnail, thumbnail)
	}

	cachedThumbnail, exists := s.cache[testURL]
	if !exists || string(cachedThumbnail) != string(testThumbnail) {
		t.Errorf("expected cached thumbnail %v, got %v", testThumbnail, cachedThumbnail)
	}
}

func TestGetThumbnailFromURL_Success(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	_, err = db.Exec("CREATE TABLE thumbnails (url TEXT PRIMARY KEY, thumbnail BLOB)")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}

	s := NewThumbnailServer(db)

	url := "https://www.youtube.com/watch?v=u4PMwpMGTdY"
	thumbnail, err := s.getThumbnailFromURL(url)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	cachedThumbnail, exists := s.cache[url]
	if !exists || string(cachedThumbnail) != string(thumbnail) {
		t.Errorf("expected euqal cached thumbnail and thumbnail from method', got %v", cachedThumbnail)
	}

	var thumbnailFromDB []byte
	err = db.QueryRow("SELECT thumbnail FROM thumbnails WHERE url = ?", url).Scan(&thumbnailFromDB)
	if err != nil {
		t.Errorf("expected thumbnail to be saved in database, but got error: %v", err)
	}
	if string(thumbnailFromDB) != string(thumbnail) {
		t.Errorf("expected thumbnail 'u4PMwpMGTdY' in database, got %v", thumbnailFromDB)
	}
}

func TestExtractVideoID_Valid(t *testing.T) {
	url := "https://youtube.com/watch?v=test123&other_param=value"
	videoID, err := extractVideoID(url)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if videoID != "test123" {
		t.Errorf("expected video ID 'test123', got %v", videoID)
	}
}
