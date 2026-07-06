package ffmpeg

import (
	"fmt"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"os"
	"os/exec"
	"path/filepath"
)

type FFmpegSettings struct {
	BinaryPath string
}

func NewFFmpegSettings(binaryPath string) *FFmpegSettings {
	return &FFmpegSettings{
		BinaryPath: binaryPath,
	}
}

func (s *FFmpegSettings) GenerateVideoQualities(inputPath string, outputDir string, outputURLPrefix string) ([]domain.VideoQuality, error) {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}
	qualities := []struct {
		Name   string
		Height int
	}{
		{Name: "1080p", Height: 1080},
		{Name: "720p", Height: 720},
		{Name: "480p", Height: 480},
		{Name: "360p", Height: 360},
	}
	result := make([]domain.VideoQuality, 0, len(qualities))

	for _, quality := range qualities {
		fileName := quality.Name + ".mp4"
		outputPath := filepath.Join(outputDir, fileName)
		outputURL := outputURLPrefix + "/" + fileName

		args := []string{
			"-y",
			"-i", inputPath,
			"-vf", fmt.Sprintf("scale=-2:%d", quality.Height),
			"-c:v", "libx264",
			"-preset", "veryfast",
			"-crf", "23",
			"-c:a", "aac",
			"-b:a", "128k",
			outputPath,
		}

		cmd := exec.Command(s.BinaryPath, args...)
		if err := cmd.Run(); err != nil {
			return nil, err
		}

		result = append(result, domain.VideoQuality{
			Quality:  quality.Name,
			VideoURL: outputURL,
		})
	}

	return result, nil
}

func (s *FFmpegSettings) GenerateArchive144p(inputPath string, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}

	args := []string{
		"-y",
		"-i", inputPath,
		"-vf", "scale=-2:144",
		"-c:v", "libx264",
		"-preset", "veryfast",
		"-crf", "30",
		"-c:a", "aac",
		"-b:a", "64k",
		outputPath,
	}

	cmd := exec.Command(s.BinaryPath, args...)
	return cmd.Run()
}
