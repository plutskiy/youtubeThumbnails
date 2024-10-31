package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	pb "testTask/pkg/api"
)

type server struct {
	pb.UnimplementedThumbnailServiceServer
	cache     map[string][]byte
	db        *sql.DB
	cacheLock sync.Mutex
	dbLock    sync.Mutex
}

func NewThumbnailServer(db *sql.DB) *server {
	return &server{
		cache: make(map[string][]byte),
		db:    db,
	}
}

func (s *server) GetThumbnail(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error) {
	if (!strings.Contains(req.Url, "youtube.com") || !strings.Contains(req.Url, "youtu.be")) && !strings.Contains(req.Url, "v=") {
		return nil, errors.New("invalid url")
	}

	s.cacheLock.Lock()
	thumbnail, exist := s.cache[req.Url]
	s.cacheLock.Unlock()

	if exist {
		return &pb.ThumbnailResponse{Image: thumbnail}, nil
	}

	thumbnail, err := s.getThumbnailFromDB(req.Url)
	if err == nil {
		return &pb.ThumbnailResponse{Image: thumbnail}, nil
	}

	thumbnail, err = s.getThumbnailFromURL(req.Url)
	if err != nil {
		return nil, err
	}

	return &pb.ThumbnailResponse{Image: thumbnail}, nil
}

func (s *server) getThumbnailFromDB(ulr string) ([]byte, error) {
	var thumbnail []byte
	err := s.db.QueryRow("SELECT thumbnail FROM thumbnails WHERE url = $1", ulr).Scan(&thumbnail)
	if err != nil {
		return nil, err
	}

	s.cacheLock.Lock()
	s.cache[ulr] = thumbnail
	s.cacheLock.Unlock()

	return thumbnail, nil
}

func (s *server) saveThumbnailToDB(url string, thumbnail []byte) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()
	_, err := s.db.Exec("INSERT INTO thumbnails (url, thumbnail) VALUES (?, ?)", url, thumbnail)
	return err
}

func (s *server) getThumbnailFromURL(url string) ([]byte, error) {
	videoID, err := extractVideoID(url)
	if err != nil {
		return nil, err
	}
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID)

	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch thumbnail")
	}

	thumbnail, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	s.cacheLock.Lock()
	s.cache[url] = thumbnail
	s.cacheLock.Unlock()

	err = s.saveThumbnailToDB(url, thumbnail)

	return thumbnail, nil
}

func extractVideoID(url string) (string, error) {
	cleanedFromDomain := strings.Split(url, "v=")[1]
	if len(cleanedFromDomain) == 0 {
		return "", errors.New("video ID not found")
	}

	extractedVideoID := strings.Split(cleanedFromDomain, "&")[0]

	return extractedVideoID, nil
}
