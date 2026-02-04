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

type CompletedTask struct {
	BranchName string
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

func getFilePath(fileName string) (string, error) {
	searchPaths := []string{
		fileName,
		filepath.Join("src", fileName),
	}

	var taskMdPath string
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			taskMdPath = path
			break
		}
	}

	if taskMdPath == "" {
		return "", fmt.Errorf("no %s found in search paths", fileName)
	}

	return taskMdPath, nil
}

// Create commands for Git Clone
func createCmdForGitClone(cfg *config.Config, branchName string) (*exec.Cmd, error) {
	switch cfg.GitHub.CloneType {
	case "SSH":
		repositoryURL := fmt.Sprintf("git@github.com:%s.git", cfg.GitHub.Repository)
		cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repositoryURL)
		return cmd, nil
	case "HTTPS":
		repositoryURL := fmt.Sprintf("https://github.com/%s.git", cfg.GitHub.Repository)
		cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", repositoryURL)
		return cmd, nil
	case "GitHub CLI":
		cmd := exec.Command("gh", "repo", "clone", cfg.GitHub.Repository, "--branch", branchName, "--single-branch")
		return cmd, nil
	default:
		return nil, errors.New("unsupported clone type is set")
	}
}

// Create commands for AI processing
func createCmdForAiProcessing(cfg *config.Config, prompt string) (*exec.Cmd, error) {
	switch cfg.AI.Type {
	case "Gemini CLI":
		cmd := exec.Command("gemini", "-p", prompt, "-y")
		if len(cfg.AI.Model) > 0 {
			cmd.Args = append(cmd.Args, "-m", cfg.AI.Model)
		}
		return cmd, nil
	case "Claude Code":
		cmd := exec.Command("claude", "-p", prompt, "-y")
		if len(cfg.AI.Model) > 0 {
			cmd.Args = append(cmd.Args, "-m", cfg.AI.Model)
		}
		return cmd, nil
	case "Codex":
		cmd := exec.Command("codex", "-y", prompt)
		if len(cfg.AI.Model) > 0 {
			cmd = exec.Command("codex", "-y", "--model", cfg.AI.Model, prompt)
		}
		return cmd, nil
	case "GitHub Copilot CLI":
		cmd := exec.Command("copilot", "-p", prompt, "--allow-all-tools")
		if len(cfg.AI.Model) > 0 {
			cmd.Args = append(cmd.Args, "--model", cfg.AI.Model)
		}
		return cmd, nil
	default:
		return nil, errors.New("unsupported AI type is set")
	}
}

// Add the processed branch name to completed_tasks.txt
func addCompletedTaskToTxt(currentDir, branchName string) error {
	path := filepath.Join(currentDir, "src", "completed_tasks.txt")

	// Open the file in append mode (create it if it doesn't exist)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.New("failed to open completed_tasks.txt")
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("%s\n", branchName)); err != nil {
		return errors.New("failed to write to completed_tasks.txt")
	}

	return nil
}

// Load task information from task.md
func LoadTaskMd() ([]Task, error) {
	// Open task.md
	taskMdPath, err := getFilePath("task.md")
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

// Load the branch names of completed tasks from completed_tasks.txt
func LoadCompletedTasks() ([]CompletedTask, error) {
	// Open completed_tasks.txt
	completedTasksTxtPath, err := getFilePath("completed_tasks.txt")
	if err != nil {
		return nil, err
	}

	file, err := os.Open(completedTasksTxtPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// /Get the branch name of a completed task
	var completedTasks []CompletedTask
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNum++

		// Skip empty lines
		if line == "" {
			continue
		}

		// Remove the newline character
		branchName := strings.ReplaceAll(line, "\r\n", "")

		completedTasks = append(completedTasks, CompletedTask{
			BranchName: branchName,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return completedTasks, nil
}

// Task execution process
func RunTask(cfg *config.Config, task Task) error {
	// Skip if the task’s skip_run_task in the config is true
	if cfg.Task.SkipRunTask {
		return nil
	}

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
	cmdGitClone, err := createCmdForGitClone(cfg, cfg.GitHub.CloneBranch)
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to create cmdGitClone: %w", err)
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
	cmdRunTask, err := createCmdForAiProcessing(cfg, task.Body)
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to create cmdRunTask: %w", err)
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

		// Append the pushed branch name to completed_tasks.txt
		if err := addCompletedTaskToTxt(currentDir, branchName); err != nil {
			os.Chdir(currentDir)
			return fmt.Errorf("failed to addCompletedTaskToTxt: %w", err)
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

// Execute additional revision process
func ExecuteAdditionalRevision(cfg *config.Config, branchName, revisionDetails string) error {
	// Skip if the task’s skip_exec_revision in the config is true
	if cfg.Task.SkipExecRevision {
		return nil
	}

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create and move to the work directory
	timestamp := time.Now().Format("20060102_150405")
	taskName := strings.Split(branchName, "/")[1]
	workDir := filepath.Join(".", "work", fmt.Sprintf("revision_%s_%s", taskName, timestamp))
	// workDir := filepath.Join(".", "work", fmt.Sprintf("%s_%s", branchName, timestamp))
	os.MkdirAll(workDir, 0755)
	os.Chdir(workDir)

	// Clone the target repository and move into its directory
	cmdGitClone, err := createCmdForGitClone(cfg, branchName)
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to create cmdGitClone: %w", err)
	}

	_, err = cmdGitClone.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	repoName := strings.Split(cfg.GitHub.Repository, "/")[1]
	os.Chdir(repoName)

	// Execute re revise
	cmdReRevise, err := createCmdForAiProcessing(cfg, revisionDetails)
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to create cmdReRevise: %w", err)
	}

	_, err = cmdReRevise.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to run re revise process: %w", err)
	}

	// Commit process
	cmdGitAdd := exec.Command("git", "add", "-A")
	_, err = cmdGitAdd.Output()
	if err != nil {
		os.Chdir(currentDir)
		return fmt.Errorf("failed to git add files: %w", err)
	}

	commitMsg := fmt.Sprintf("aidd: [%s_%s] Revision", branchName, timestamp)
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

		// Add a comment with the correction details to the PR (assuming the PR has already been created)
		if cfg.GitHub.CreatePrOnComplete {
			bodyText := fmt.Sprintf("【Revision details】\n%s", revisionDetails)
			cmdAddCommentToPR := exec.Command("gh", "pr", "comment", "--body", bodyText)
			_, err = cmdAddCommentToPR.Output()
			if err != nil {
				os.Chdir(currentDir)
				return fmt.Errorf("failed to add comment to PR: %w", err)
			}
		}
	}

	// Return to the current directory
	os.Chdir(currentDir)

	return nil
}
