package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testMultipartFile struct {
	*strings.Reader
}

func (f testMultipartFile) Close() error {
	return nil
}

func TestVideoStorageSaveFileAndPaths(t *testing.T) {
	basePath := t.TempDir()
	storage := NewVideoStorage(basePath, "uploads")

	file := testMultipartFile{Reader: strings.NewReader("video-data")}
	header := &multipart.FileHeader{Filename: "clip.MP4"}

	url, err := storage.SaveFile(file, header, "videos/originals")
	if err != nil {
		t.Fatalf("SaveFile returned error: %v", err)
	}
	if !strings.HasPrefix(url, "/uploads/videos/originals/") || !strings.HasSuffix(url, ".mp4") {
		t.Fatalf("unexpected url: %q", url)
	}

	savedPath := storage.URLPath(url)
	data, err := os.ReadFile(savedPath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	if string(data) != "video-data" {
		t.Fatalf("unexpected saved data: %q", string(data))
	}

	outputDir, outputURLPrefix := storage.VideoQualitiesPaths(10)
	if outputDir != filepath.Join(basePath, "videos", "10") {
		t.Fatalf("unexpected qualities dir: %q", outputDir)
	}
	if outputURLPrefix != "/uploads/videos/10" {
		t.Fatalf("unexpected qualities url prefix: %q", outputURLPrefix)
	}

	archivePath, archiveURL := storage.VideoArchivePath(10)
	if archivePath != filepath.Join(basePath, "videos", "10", "archive", "144p.mp4") {
		t.Fatalf("unexpected archive path: %q", archivePath)
	}
	if archiveURL != "/uploads/videos/10/archive/144p.mp4" {
		t.Fatalf("unexpected archive url: %q", archiveURL)
	}
}

func TestVideoStorageRemoveFileByURL(t *testing.T) {
	basePath := t.TempDir()
	storage := NewVideoStorage(basePath, "uploads")
	file := testMultipartFile{Reader: strings.NewReader("data")}
	url, err := storage.SaveFile(file, &multipart.FileHeader{Filename: "file.txt"}, "files")
	if err != nil {
		t.Fatalf("SaveFile returned error: %v", err)
	}

	if err := storage.RemoveFileByURL(url); err != nil {
		t.Fatalf("RemoveFileByURL returned error: %v", err)
	}
	if _, err := os.Stat(storage.URLPath(url)); !os.IsNotExist(err) {
		t.Fatalf("expected file to be removed, got err=%v", err)
	}
	if err := storage.RemoveFileByURL(url); err != nil {
		t.Fatalf("removing missing file should be safe: %v", err)
	}
	if err := storage.RemoveFileByURL(""); err != nil {
		t.Fatalf("empty url should be ignored: %v", err)
	}
}

func TestVideoStorageRemoveDir(t *testing.T) {
	basePath := t.TempDir()
	storage := NewVideoStorage(basePath, "uploads")
	dir := filepath.Join(basePath, "videos", "1")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		t.Fatalf("failed to prepare dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "360p.mp4"), []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to prepare file: %v", err)
	}

	if err := storage.RemoveDir(dir); err != nil {
		t.Fatalf("RemoveDir returned error: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("expected dir to be removed, got err=%v", err)
	}
	if err := storage.RemoveDir(""); err != nil {
		t.Fatalf("empty dir should be ignored: %v", err)
	}
}

func TestVideoStorageRejectsRemovingOutsideBasePath(t *testing.T) {
	storage := NewVideoStorage(t.TempDir(), "uploads")
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "file.txt")
	if err := os.WriteFile(outsideFile, []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to prepare outside file: %v", err)
	}

	if err := storage.RemoveDir(outsideDir); err == nil {
		t.Fatal("expected error when removing outside base path")
	}
	if err := storage.removeInsideBasePath(outsideFile, false); err == nil {
		t.Fatal("expected error when removing outside file")
	}
}

var _ multipart.File = testMultipartFile{}
var _ io.Closer = testMultipartFile{}
