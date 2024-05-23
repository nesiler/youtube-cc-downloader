package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func loadChannels(filename string) []string {
	Info("Loading channels from file: ", filename)
	file, err := os.Open(filename)
	if err != nil {
		Err(os.Stderr, "Error opening channels file:\n", err) // Print error to stderr
		return nil
	}
	defer file.Close()

	var channels []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		channels = append(channels, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		Err(os.Stderr, "Error reading channels file:\n", err)
		return nil
	}

	Out("Loaded channel count: ", len(channels))
	return channels
}

func fetchVideoIDs(channelID, count string) []string {
	Info("Fetching video IDs for channel: ", channelID)
	cmd := exec.Command("youtube-dl", "--get-id", "--playlist-end", count, "https://www.youtube.com/channel/"+channelID)

	output, err := cmd.Output()
	if err != nil {
		Err(os.Stderr, "Error fetching video IDs:\n", err)
		return nil
	}

	videoIDs := strings.Split(strings.TrimSpace(string(output)), "\n")
	Out("Fetched video IDs:", videoIDs)
	return videoIDs
}

func downloadSubtitle(videoID string) (string, error) {
	Info("Downloading subtitle for video: ", videoID)
	outputDir := filepath.Join("tmp")
	os.MkdirAll(outputDir, 0755) // Ensure output directory exists

	cmd := exec.Command("youtube-dl", "--write-sub", "--write-auto-sub", "--sub-lang", "en", "--skip-download", "-o", filepath.Join(outputDir, "%(id)s.%(ext)s"), "https://www.youtube.com/watch?v="+videoID)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("\nerror downloading subtitle: \n", err)
	}

	// Find the subtitle file with consistent naming (based on youtube-dl's default)
	subtitleFilePattern := filepath.Join(outputDir, videoID+".en.*")
	subtitleFiles, err := filepath.Glob(subtitleFilePattern)
	if err != nil || len(subtitleFiles) == 0 {
		return "", fmt.Errorf("\nerror finding subtitle file: \n", err)
	}

	subtitleBytes, err := os.ReadFile(subtitleFiles[0])
	if err != nil {
		return "", fmt.Errorf("\nerror reading subtitle file: \n", err) // Wrap error with '%w'
	}

	// Convert subtitle content to plain text (assuming .vtt)
	subtitleText := string(subtitleBytes)
	subtitleText = regexp.MustCompile(`(WEBVTT|Kind:.+|Language:.+|\d\d:\d\d:\d\d\.\d\d\d --> \d\d:\d\d:\d\d\.\d\d\d.*\n)`).ReplaceAllString(subtitleText, "") // Remove VTT headers, timestamps, and blank lines after timestamps
	subtitleText = strings.TrimSpace(subtitleText)
	subtitleText = regexp.MustCompile("<.*?>").ReplaceAllString(subtitleText, "") // Remove tags

	// Clean up temporary files
	// os.RemoveAll(outputDir)

	return subtitleText, nil
}

func saveSubtitleToFile(subtitleContent, channelID, videoID string) {
	Info("Saving subtitle to file for video: ", videoID)
	channelDir := filepath.Join("subtitles", channelID)
	if err := os.MkdirAll(channelDir, 0755); err != nil {
		Err(os.Stderr, "Error creating channel directory: \n", err)
		return
	}

	filename := filepath.Join(channelDir, videoID+".txt") // Assuming plain text subtitles
	if err := os.WriteFile(filename, []byte(subtitleContent), 0644); err != nil {
		Err(os.Stderr, "Error saving subtitle to file: \n", err)
		return
	}
}
