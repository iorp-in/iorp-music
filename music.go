package main // Declares that this file is part of the main package

import (
	"errors" // Package for handling errors
	"flag"   // Package for parsing command-line flags

	// Package for formatting strings
	"io"       // Package for handling I/O operations
	"net/http" // Package for HTTP client and server implementations
	"os/exec"  // Package for running external commands
	"regexp"   // Package for regular expressions
	"strconv"  // Package for converting strings to numbers
	"strings"  // Package for string manipulation

	"github.com/kkdai/youtube/v2" // External library for interacting with YouTube API
)

// Function to extract video ID from a YouTube search query
func getVideoIDFromQuery(query string) (string, error) {
	// Replace spaces with '+' to form a valid URL-encoded query
	encodedQuery := strings.ReplaceAll(query, " ", "+")

	// Construct the search URL
	searchURL := "https://www.youtube.com/results?q=" + encodedQuery

	// Fetch search results from YouTube
	response, err := http.Get(searchURL)
	if err != nil {
		return "", errors.New("error fetching search results: " + err.Error())
	}
	defer response.Body.Close() // Close the response body when function exits

	// Read the response body into a byte slice
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("error reading response body: " + err.Error())
	}

	// Define a regular expression to extract video ID from search results
	var videoIDRegex = regexp.MustCompile(`(?m){"webCommandMetadata":{"url":"\/watch\?v=(.*?)[\\|"]`)

	// Find sub matches in the response body using the regex
	matches := videoIDRegex.FindStringSubmatch(string(data))
	// Check if a video ID is found
	if len(matches) < 2 {
		return "", errors.New("no video found in search results")
	}

	// Return the extracted video ID
	return matches[1], nil
}

// HTTP handler for processing YouTube requests
func YoutubeHandler(w http.ResponseWriter, r *http.Request) {
	flag.Parse() // Parse command-line flags

	// Extract the query from the URL path
	query := strings.TrimPrefix(r.URL.Path, "/youtube/")

	// Get video ID corresponding to the query
	videoID, err := getVideoIDFromQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Initialize a YouTube client
	client := youtube.Client{}

	// Construct the full YouTube video URL
	videoURL := "https://www.youtube.com/watch?v=" + videoID

	// Get video metadata from YouTube
	video, err := client.GetVideo(videoURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get stream URL for the video
	streamURL, err := client.GetStreamURL(video, &video.Formats.AudioChannels(2).Quality("medium")[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Title", video.Title)
	w.Header().Set("X-Content-Duration", strconv.Itoa(int(video.Duration)))
	w.Header().Set("Connection", "keep-alive")

	// Find the path to the ffmpeg executable
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		http.Error(w, "ffmpeg not found!", http.StatusInternalServerError)
		return
	}

	// Prepare ffmpeg command
	ffmpegCmd := exec.Command(ffmpegPath, "-i", streamURL, "-acodec", "libmp3lame", "-f", "mp3", "-")
	ffmpegCmd.Stdout = w // Set the command's standard output to be the HTTP response writer

	// Execute ffmpeg command
	ffmpegCmd.Run()
}
