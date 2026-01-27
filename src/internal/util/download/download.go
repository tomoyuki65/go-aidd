package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// Save downloaded image to "images" directory
func SaveImages(provider, url, token, filePath string) error {
	// Request configuration
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if provider == "GitHub" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	// Execute request
	client := &http.Client{}
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
