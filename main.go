package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
)

var (
	Head              = color.New(color.FgHiMagenta).Add(color.Bold).Add(color.Underline).Add(color.BgHiWhite).PrintlnFunc()
	Out               = color.New(color.FgHiWhite).PrintlnFunc()
	Info              = color.New(color.FgHiCyan).PrintlnFunc()
	Warn              = color.New(color.FgHiYellow).Add(color.Bold).PrintlnFunc()
	Err               = color.New(color.FgHiRed).Add(color.Bold).FprintfFunc()
	Ok                = color.New(color.FgHiGreen).PrintlnFunc()
	downloadErrorsLog *os.File
)

func init() {
	// Open or create the log file
	var err error
	downloadErrorsLog, err = os.OpenFile("download_errors.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		Err(os.Stderr, "Error opening download_errors.log: %v\n", err)
		os.Exit(1) // Exit if the log file cannot be opened
	}
}

func logFailedDownload(channelID, videoID, errorMessage string) {
	logEntry := fmt.Sprintf("Channel ID: %s, Video ID: %s, Error: %s\n", channelID, videoID, errorMessage)
	if _, err := downloadErrorsLog.WriteString(logEntry); err != nil {
		Err(os.Stderr, "Error writing to download_errors.log: %v\n", err)
	}
}
func clearFiles() {
	Info("Clearing all downloaded files and directories...")
	if err := os.RemoveAll("tmp"); err != nil {
		Err(os.Stderr, "Error clearing temporary files: %v\n", err) // Print error to stderr
	}
	if err := os.RemoveAll("subtitles"); err != nil {
		Err(os.Stderr, "Error clearing subtitles: %v\n", err) // Print error to stderr
	}
}

func skipDownloaded(videoIDs []string, channelID string) []string {
	var filteredVideoIDs []string
	for _, videoID := range videoIDs {
		// Construct the expected file path for this video
		channelDir := filepath.Join("subtitles", channelID)
		filePath := filepath.Join(channelDir, videoID+".txt")

		// Check if the file exists
		if _, err := os.Stat(filePath); err == nil {
			Info("Skipping already downloaded video:", videoID)
			continue // Skip this video
		}
		filteredVideoIDs = append(filteredVideoIDs, videoID) // Add it if not found
	}
	return filteredVideoIDs
}

func downloadSubtitles(fileName, count string) {
	channelIDs := loadChannels(fileName)
	if len(channelIDs) == 0 {
		Warn("No channels loaded. Exiting.")
		return
	}

	if count == "" {
		count = "10"
	}

	var wg sync.WaitGroup // WaitGroup to synchronize goroutines
	for _, channelID := range channelIDs {
		wg.Add(1) // Increment WaitGroup counter for each channel

		go func(channelID string) { // Start a goroutine for each channel
			defer wg.Done() // Decrement counter when the goroutine finishes

			videoIDs := fetchVideoIDs(channelID, count)
			if len(videoIDs) == 0 {
				Warn("No videos found for channel, skipping: ", channelID)
				return
			}

			videoIDs = skipDownloaded(videoIDs, channelID)

			for _, videoID := range videoIDs {
				subtitleContent, err := downloadSubtitle(videoID)
				if err != nil {
					Warn("Error downloading subtitle for video: ", videoID, err)
					logFailedDownload(channelID, videoID, err.Error()) // Log the error
					continue
				}
				saveSubtitleToFile(subtitleContent, channelID, videoID)
				Ok("Subtitle for video: " + videoID + " saved successfully.")
			}
		}(channelID) // Pass channelID to the goroutine
	}

	wg.Wait()

	// check if there are any subtitles downloaded
	if _, err := os.Stat("subtitles"); os.IsNotExist(err) {
		Warn("No subtitles downloaded. Exiting.")
		return
	}

	// Close the log file at the end of the program
	if err := downloadErrorsLog.Close(); err != nil {
		Err(os.Stderr, " Error closing download_errors.log: \n", err)
	}

	Head("=== Subtitles downloaded successfully ===")
}

func main() {
	args := os.Args[1:] // Exclude the program name from arguments

	if len(args) > 0 {
		if args[0] == "clear" {
			Head("=== Cleaner Mode ===")
			clearFiles()
			Ok("=== All files and directories cleared successfully ===")
			return
		}

		if args[0] == "csv" {
			Head("=== CSV Mode ===")
			generateCSV("output.csv")
			return
		}

		if args[0] == "download" {
			Head("=== Downloader Mode ===")
			downloadSubtitles("channels.txt", "15")
			return
		}
	}

	Head("=== Default Mode ===")
	downloadSubtitles("channels.txt", "15")
	generateCSV("output.csv")
}
