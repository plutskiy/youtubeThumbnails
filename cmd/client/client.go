package client

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	pb "testTask/pkg/api"
	"time"
)

func sanitizeFilename(url string) string {
	re := regexp.MustCompile(`[^\w]+`)
	sanitized := re.ReplaceAllString(url, "-")
	if len(sanitized) > 100 {
		hash := sha1.Sum([]byte(url))
		return hex.EncodeToString(hash[:])
	}
	return sanitized
}

func DownloadAndSaveThumbnail(server pb.ThumbnailServiceServer, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := pb.ThumbnailRequest{Url: url}
	resp, err := server.GetThumbnail(ctx, &req)
	if err != nil {
		log.Printf("Failed to get thumbnail for %s: %v", url, err)
		return
	}

	filename := fmt.Sprintf("%s.jpg", sanitizeFilename(url))
	dirPath := "image"
	filePath := filepath.Join(dirPath, filename)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			log.Printf("Failed to create directory %s: %v", dirPath, err)
			return
		}
	}

	err = os.WriteFile(filePath, resp.Image, 0644)
	if err != nil {
		log.Printf("Failed to save thumbnail for %s: %v", url, err)
		return
	}
	fmt.Printf("Thumbnail for %s saved as %s\n", url, filename)
}
