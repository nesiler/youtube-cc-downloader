package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getVideoMetadata(videoID string) (string, string, error) {
	cmd := exec.Command("youtube-dl", "--get-title", "--get-description", "--get-filename", "https://www.youtube.com/watch?v="+videoID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("error getting video metadata: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) < 3 {
		return "", "", fmt.Errorf("incomplete metadata received")
	}

	title := lines[1]
	description := lines[2]

	return title, description, nil
}

func generateCSV(filename string) {
	Info("Generating CSV file: ", filename)
	file, err := os.Create(filename)
	if err != nil {
		Err(os.Stderr, "Error creating CSV file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"video_id", "title", "channel_id", "description", "subtitle"}
	if err := writer.Write(header); err != nil {
		Err(os.Stderr, "Error writing CSV header: %v\n", err)
		return
	}

	err = filepath.Walk("subtitles", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			videoID := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			channelTitle := filepath.Base(filepath.Dir(path))

			subtitleBytes, err := os.ReadFile(path)
			if err != nil {
				Err(os.Stderr, "Error reading subtitle file %s: %v\n", path, err)
				return nil
			}
			subtitleContent := string(subtitleBytes)
			// Add this line to remove newlines from subtitleContent
			subtitleContent = strings.ReplaceAll(subtitleContent, "\n", " ")

			title, description, err := getVideoMetadata(videoID)
			if err != nil {
				Warn("Error getting metadata for video %s: %v\n Metadata will be empty.", videoID, err)
				// Leave the fields empty if metadata retrieval fails
				title, description = "", ""
			}

			record := []string{videoID, title, channelTitle, description, subtitleContent}
			if err := writer.Write(record); err != nil {
				Err(os.Stderr, "Error writing CSV row: %v\n", err)
			}
		}
		return nil
	})

	if err != nil {
		Err(os.Stderr, "Error walking subtitles directory: %v\n", err)
		return
	}

	Head("=== CSV file generated successfully ===")
}
