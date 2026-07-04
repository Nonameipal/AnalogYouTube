package storage
import (
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"fmt"
	"io"

)
type VideoStorage struct {
	BasePath string
	BaseURL  string
}

func NewVideoStorage(basePath, baseURL string) *VideoStorage {
	return &VideoStorage{
		BasePath: basePath,
		BaseURL:  baseURL,
	}
}

func (s *VideoStorage) SaveFile(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	if err := os.MkdirAll(filepath.Join(s.BasePath, folder), os.ModePerm); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

		fullPath := filepath.Join(s.BasePath, folder, fileName)

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/" + filepath.ToSlash(filepath.Join(s.BaseURL, folder, fileName)), nil
}

func (s *VideoStorage) URLPath(fileURL string) string {
	cleanURL := strings.TrimPrefix(fileURL, "/")
	cleanURL = strings.TrimPrefix(cleanURL, s.BaseURL)

	return filepath.Join(s.BasePath, cleanURL)
}

func (s *VideoStorage) VideoQualitiesPaths(videoID int) (string, string) {
	folder := fmt.Sprintf("videos/%d", videoID)

	outputDir := filepath.Join(s.BasePath, folder)
	outputURLPrefix := "/" + filepath.ToSlash(filepath.Join(s.BaseURL, folder))

	return outputDir, outputURLPrefix
}