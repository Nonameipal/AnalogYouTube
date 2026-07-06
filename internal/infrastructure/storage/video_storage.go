package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	cleanURL = strings.TrimPrefix(cleanURL, s.BaseURL+"/")

	return filepath.Join(s.BasePath, cleanURL)
}
func (s *VideoStorage) VideoQualitiesPaths(videoID int) (string, string) {
	folder := fmt.Sprintf("videos/%d", videoID)

	outputDir := filepath.Join(s.BasePath, folder)
	outputURLPrefix := "/" + filepath.ToSlash(filepath.Join(s.BaseURL, folder))

	return outputDir, outputURLPrefix
}
func (s *VideoStorage) VideoArchivePath(videoID int) (string, string) {
	folder := fmt.Sprintf("videos/%d/archive", videoID)
	fileName := "144p.mp4"

	outputPath := filepath.Join(s.BasePath, folder, fileName)
	archiveURL := "/" + filepath.ToSlash(filepath.Join(s.BaseURL, folder, fileName))

	return outputPath, archiveURL
}

func (s *VideoStorage) RemoveFileByURL(fileURL string) error {
	if fileURL == "" {
		return nil
	}

	return s.removeInsideBasePath(s.URLPath(fileURL), false)
}

func (s *VideoStorage) RemoveDir(path string) error {
	if path == "" {
		return nil
	}

	return s.removeInsideBasePath(path, true)
}

func (s *VideoStorage) removeInsideBasePath(path string, recursive bool) error {
	basePath, err := filepath.Abs(s.BasePath)
	if err != nil {
		return err
	}

	targetPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	relativePath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return err
	}

	if relativePath == "." || relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("storage path is outside base path")
	}

	if recursive {
		return os.RemoveAll(targetPath)
	}

	err = os.Remove(targetPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
