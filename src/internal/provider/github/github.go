package github

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	dl "github.com/tomoyuki65/go-aidd/internal/util/download"
)

type GitHubIssue struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Generate task.md from GitHub issues for the given repository and label
func GenerateTaskMd(repository, label string) error {
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
	var issues []GitHubIssue
	if err := json.Unmarshal(outputGhIssueList, &issues); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Create the src directory
	if err := os.MkdirAll("src", 0755); err != nil {
		return fmt.Errorf("failed to create src directory: %w", err)
	}

	// Write to task.md
	file, err := os.Create(filepath.Join("src", "task.md"))
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
			imgDir := filepath.Join("src", "images", fmt.Sprintf("issue_%d", issue.Number))
			os.MkdirAll(imgDir, 0755)

			fileName := fmt.Sprintf("img_%d.png", i+1)
			filePath := filepath.Join(imgDir, fileName)

			// Download the image and save it to a local file
			if err := dl.SaveImages("GitHub", url, token, filePath); err == nil {
				// Replace the image URL in the issue body with the local image path
				relPath := filepath.Join("images", fmt.Sprintf("issue_%d", issue.Number), fileName)
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
