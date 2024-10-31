package client

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	pb "testTask/pkg/api"
)

type mockThumbnailServiceServer struct {
	pb.UnimplementedThumbnailServiceServer
	mockGetThumbnail func(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error)
}

func (m *mockThumbnailServiceServer) GetThumbnail(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error) {
	return m.mockGetThumbnail(ctx, req)
}

func TestDownloadAndSaveThumbnail(t *testing.T) {
	mockServer := &mockThumbnailServiceServer{
		mockGetThumbnail: func(ctx context.Context, req *pb.ThumbnailRequest) (*pb.ThumbnailResponse, error) {
			if req.Url == "https://valid.url" {
				return &pb.ThumbnailResponse{Image: []byte("thumbnail data")}, nil
			}
			return nil, errors.New("invalid URL")
		},
	}

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"Valid URL", "https://valid.url", false},
		{"Invalid URL", "https://invalid.url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirPath := "image"
			filePath := filepath.Join(dirPath, sanitizeFilename(tt.url)+".jpg")

			err := os.RemoveAll(dirPath)
			if err != nil {
				return
			}

			DownloadAndSaveThumbnail(mockServer, tt.url)

			_, err = os.Stat(filePath)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				err := os.RemoveAll(dirPath)
				if err != nil {
					return
				}
			}
		})
	}
}
