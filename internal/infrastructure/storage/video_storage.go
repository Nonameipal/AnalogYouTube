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
