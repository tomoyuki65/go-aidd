package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tomoyuki65/go-aidd/internal/config"
)

type Issue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Save downloaded image to "images" directory
func downloadImage(client *http.Client, url, token, filePath string) error {
	// Request configuration
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"failed to download file from %s: unexpected status code %d %s",
			req.URL.String(),
			resp.StatusCode,
			http.StatusText(resp.StatusCode),
		)
	}
	defer resp.Body.Close()

	// Create file
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("file create error: %v", err)
	}
	defer out.Close()

	// Write to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// Generate task.md from GitHub issues for the given repository and label
func generateTaskMdGitHub(client *http.Client, repository, label string) error {
	// Set up regex to extract image URLs
	re := regexp.MustCompile(`https://github\.com/[^\s\)\"]+/(?:user-attachments|assets|user-images)/[^\s\)\"]+`)

	// Retrieve authentication token from GitHub CLI
	cmdGhAuthToken := exec.Command("gh", "auth", "token")
	outputGhAuthToken, err := cmdGhAuthToken.Output()
	if err != nil {
		return fmt.Errorf("failed to fetch GitHub auth token: %w", err)
	}
	token := strings.TrimSpace(string(outputGhAuthToken))

	// Fetch the target issue in JSON format
	cmdGhIssueList := exec.Command("gh", "issue", "list",
		"-R", repository,
		"--label", label,
		"--json", "number,title,body",
	)
	outputGhIssueList, err := cmdGhIssueList.Output()
	if err != nil {
		return fmt.Errorf("failed to fetch task information: %w", err)
	}

	// Parse GitHub issue list JSON into Issue structs
	var issues []Issue
	if err := json.Unmarshal(outputGhIssueList, &issues); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Write to task.md
	file, err := os.Create("task.md")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintln(file, "| Number | Title | Body |")
	fmt.Fprintln(file, "| --- | --- | --- |")

	// Write tasks
	for _, issue := range issues {
		// Extract all URLs from the issue body
		urls := re.FindAllString(issue.Body, -1)

		// Replace newline characters (\n or \r\n) with <br>
		safeBody := strings.ReplaceAll(issue.Body, "\n", "<br>")
		safeBody = strings.ReplaceAll(safeBody, "\r", "")

		// Escape pipe characters (|)
		safeBody = strings.ReplaceAll(safeBody, "|", "\\|")
		safeTitle := strings.ReplaceAll(issue.Title, "|", "\\|")

		// Process each extracted URL
		for i, url := range urls {
			// Prepare local file path for saving issue images
			imgDir := fmt.Sprintf("./images/issue_%d", issue.Number)
			os.MkdirAll(imgDir, 0755)

			fileName := fmt.Sprintf("img_%d.png", i+1)
			filePath := filepath.Join(imgDir, fileName)

			// Download the image and save it to a local file
			if err := downloadImage(client, url, token, filePath); err == nil {
				// Replace the image URL in the issue body with the local image path
				relPath := fmt.Sprintf("./images/issue_%d/%s", issue.Number, fileName)
				safeBody = strings.ReplaceAll(safeBody, url, relPath)
			} else {
				return fmt.Errorf("failed to download image: %w", err)
			}
		}

		// Write task
		fmt.Fprintf(file, "| %d | %s | %s |\n", issue.Number, safeTitle, safeBody)
	}

	fmt.Println("Task information successfully written to task.md.")

	return nil
}

// Generate "task.md" from task information
func generateTaskMd(provider, repository, label string) error {
	// HTTP client configuration
	client := &http.Client{}

	// Switch processing by provider
	switch provider {
	case "GitHub":
		return generateTaskMdGitHub(client, repository, label)
	default:
		return errors.New("unsupported provider is set")
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Generate "src/task.md" from task info
	if err := generateTaskMd(cfg.Issue.Provider, cfg.GitHub.Repository, cfg.Issue.Label); err != nil {
		log.Fatalf("failed to generate file 'task.md': %v", err)
	}
}
