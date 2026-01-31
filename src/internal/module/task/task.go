package task

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tomoyuki65/go-aidd/internal/config"
	"github.com/tomoyuki65/go-aidd/internal/provider/container"
	"github.com/tomoyuki65/go-aidd/internal/provider/github"
)

type Task struct {
	Number int
	Title  string
	Body   string
}

// Generate "task.md" from task information
func GenerateTaskMd(cfg *config.Config) error {
	// Switch processing by provider
	switch cfg.Issue.Provider {
	case "GitHub":
		return github.GenerateTaskMd(cfg.GitHub.Repository, cfg.Issue.Label)
	case "container":
		return container.ExecContainer()
	default:
		return errors.New("unsupported provider is set")
	}
}

func getTaskMdPath() (string, error) {
	searchPaths := []string{
		"task.md",
		filepath.Join("src", "task.md"),
	}

	var taskMdPath string
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			taskMdPath = path
			break
		}
	}

	if taskMdPath == "" {
		return "", errors.New("no task.md found in search paths")
	}

	return taskMdPath, nil
}

// Load task information from task.md
func LoadTaskMd() ([]Task, error) {
	// Open task.md
	taskMdPath, err := getTaskMdPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(taskMdPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Retrieve task information
	var tasks []Task
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		// Skip header and separator lines
		if lineNum <= 2 {
			continue
		}

		// Skip empty lines
		if line == "" {
			continue
		}

		// Split by |
		cols := strings.Split(line, "|")
		// In Markdown tables, cells at the start and end are blank, so fewer than 4 cells is an error
		if len(cols) < 4 {
			return nil, errors.New("invalid table row at line " + strconv.Itoa(lineNum))
		}

		// Trim each item
		numberStr := strings.TrimSpace(cols[1])
		title := strings.TrimSpace(cols[2])
		body := strings.TrimSpace(cols[3])

		// Convert to number
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return nil, fmt.Errorf("invalid number at line %d: %v", lineNum, err)
		}

		// Convert <br> in body to newline ("\n")
		body = strings.ReplaceAll(body, "<br>", "\n")

		tasks = append(tasks, Task{
			Number: number,
			Title:  title,
			Body:   body,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Task execution process
func RunTask(cfg *config.Config, task Task) error {
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create and move to the work directory
	timestamp := time.Now().Format("20060102_150405")
	workDir := filepath.Join(".", "work", fmt.Sprintf("task_%d_%s", task.Number, timestamp))
	os.MkdirAll(workDir, 0755)
	os.Chdir(workDir)

	// Clone the target repository and move into its directory
	var cmdGitClone *exec.Cmd
	switch cfg.GitHub.CloneType {
	case "SSH":
		repositoryURL := fmt.Sprintf("git@github.com:%s.git", cfg.GitHub.Repository)
		cmdGitClone = exec.Command("git", "clone", "-b", cfg.GitHub.CloneBranch, "--single-branch", repositoryURL)
	case "HTTPS":
		repositoryURL := fmt.Sprintf("https://github.com/%s.git", cfg.GitHub.Repository)
		cmdGitClone = exec.Command("git", "clone", "-b", cfg.GitHub.CloneBranch, "--single-branch", repositoryURL)
	case "GitHub CLI":
		cmdGitClone = exec.Command("gh", "repo", "clone", cfg.GitHub.Repository, "--branch", cfg.GitHub.CloneBranch, "--single-branch")
	default:
		// Return to the current directory
		os.Chdir(currentDir)
		return errors.New("unsupported clone type is set")
	}

	_, err = cmdGitClone.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	repoName := strings.Split(cfg.GitHub.Repository, "/")[1]
	os.Chdir(repoName)

	// Check if the branch exists
	branchName := fmt.Sprintf("aidd/task_%d", task.Number)
	cmdCheckBranch := exec.Command("git", "ls-remote", "--heads", "origin", branchName)
	out, err := cmdCheckBranch.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to check branch: %w", err)
	} else if len(out) > 0 {
		os.Chdir(currentDir)
		return fmt.Errorf("branch '%s' already exists", branchName)
	}

	// Create the branch
	cmdGitCheckout := exec.Command("git", "checkout", "-b", branchName)
	_, err = cmdGitCheckout.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to create branch: %w", err)
	}

	// Execute the task
	var cmdRunTask *exec.Cmd
	switch cfg.AI.Type {
	case "Gemini CLI":
		cmdRunTask = exec.Command("gemini", "-p", task.Body, "-y")
		if len(cfg.AI.Model) > 0 {
			cmdRunTask.Args = append(cmdRunTask.Args, "-m", cfg.AI.Model)
		}
	default:
		os.Chdir(currentDir)
		return errors.New("unsupported AI type is set")
	}

	_, err = cmdRunTask.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to run task: %w", err)
	}

	// Commit process
	cmdGitAdd := exec.Command("git", "add", "-A")
	_, err = cmdGitAdd.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to git add files: %w", err)
	}

	commitMsg := fmt.Sprintf("aidd: [task_%d] %s", task.Number, task.Title)
	cmdGitCommit := exec.Command("git", "commit", "-m", commitMsg)
	_, err = cmdGitCommit.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to git commit: %w", err)
	}

	// Push to GitHub
	if cfg.GitHub.PushBranchOnComplete {
		cmdGitPush := exec.Command("git", "push", "-u", "origin", branchName)
		_, err = cmdGitPush.Output()
		if err != nil {
			os.Chdir(currentDir)
			return fmt.Errorf("failed to git push: %w", err)
		}

		// Create a pull request
		if cfg.GitHub.CreatePrOnComplete {
			bodyText := fmt.Sprintf("【Task Detail】\n%s", task.Body)

			cmdCreatePullRequest := exec.Command("gh", "pr", "create",
				"--base", cfg.GitHub.CloneBranch,
				"--head", branchName,
				"--title", commitMsg,
				"--body", bodyText,
				"--label", cfg.Issue.Label,
			)

			if cfg.GitHub.PrDraft {
				cmdCreatePullRequest.Args = append(cmdCreatePullRequest.Args, "--draft")
			}

			_, err = cmdCreatePullRequest.Output()
			if err != nil {
				os.Chdir(currentDir)
				return fmt.Errorf("failed to create pull request: %w", err)
			}
		}
	}

	// Return to the current directory
	os.Chdir(currentDir)

	return nil
}
