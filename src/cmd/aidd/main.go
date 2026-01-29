package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

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
func generateTaskMd(cfg *config.Config) error {
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

// Load task information from task.md
func loadTaskMd() ([]Task, error) {
	// Set path for task.md
	exePath, _ := os.Executable()
	binDir := filepath.Dir(exePath)
	path := filepath.Join(binDir, "..", "task.md")

	if os.Getenv("ENV") == "local" {
		path = "task.md"
	}

	// Open task.md
	file, err := os.Open(path)
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

// func runTask(cfg *config.Config, task Task) {
// 	// 実行する処理を記述
// 	//
// }

// Display task details
func showTaskDetail(cfg *config.Config, app *tview.Application, pages *tview.Pages, task Task) {
	description := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Would you like to run this task ?[-]")

	separator := tview.NewTextView().
		SetDynamicColors(true).
		SetText("------------------------------------------------------------------------------------------")

	taskInfoText := fmt.Sprintf("Number: %d\nTitle: %s\n\nBody:\n-----\n%s", task.Number, task.Title, task.Body)

	taskInfo := tview.NewTextView().
		SetDynamicColors(true).
		SetWordWrap(true).
		SetText(taskInfoText)

	// Confirmation form settings
	confirmForm := tview.NewForm().
		AddButton("Back", func() {
			// Remove task details from the page settings and return
			pages.RemovePage("task_detail")
		}).
		AddButton("Run", func() {
			// 実行する処理を記述
			//
			//
		})

	// Allow arrow key navigation between buttons
	confirmForm.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown, tcell.KeyRight:
			return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
		case tcell.KeyUp, tcell.KeyLeft:
			return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
		}
		return event
	})

	// Set task information height
	taskInfoHeight := strings.Count(task.Body, "\n") + 6

	taskDetailView := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(taskInfo, taskInfoHeight, 1, false).
		AddItem(confirmForm, 0, 1, true)

	//  Set task details height to 18 or higher
	taskDetailHeight := 18 + taskInfoHeight

	taskDetail := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(description, 2, 1, false).
		AddItem(separator, 1, 1, false).
		AddItem(taskDetailView, taskDetailHeight, 1, true).
		AddItem(separator, 1, 1, false).
		AddItem(nil, 0, 1, false)
	taskDetail.SetBorder(true).SetTitle(" Task detail ")

	// Add task details to the page settings and display them (focus on the confirmation form)
	pages.AddPage("task_detail", taskDetail, true, true)
	app.SetFocus(confirmForm)
}

// Task list display process
func renderTasks(cfg *config.Config, app *tview.Application, taskList *tview.List, pages *tview.Pages, tasks []Task, currentPage, pageSize *int) {
	taskList.Clear()

	// Calculate the page range
	start := *currentPage * *pageSize
	end := start + *pageSize
	if end > len(tasks) {
		end = len(tasks)
	}

	// Display the task list for the current page
	for _, t := range tasks[start:end] {
		task := t
		taskList.AddItem(fmt.Sprintf("%d. %s", task.Number, task.Title), "", 0, func() {
			// Display task details
			showTaskDetail(cfg, app, pages, task)
		})
	}

	// Set up pagination
	pagination := fmt.Sprintf("[green]-- page: %d / %d --[-]", *currentPage+1, (len(tasks)-1) / *pageSize + 1)
	taskList.AddItem(pagination, "", 0, nil)

	if end < len(tasks) {
		taskList.AddItem("▶ Next page", "", 'n', func() {
			*currentPage++
			renderTasks(cfg, app, taskList, pages, tasks, currentPage, pageSize)
		})
	}
	if *currentPage > 0 {
		taskList.AddItem("◀ Back page", "", 'b', func() {
			*currentPage--
			renderTasks(cfg, app, taskList, pages, tasks, currentPage, pageSize)
		})
	}

	// Handle returning to the main menu
	taskList.AddItem("Return to the main menu", "", 'r', func() {
		pages.SwitchToPage("main_menu")
	})
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Define the app using tview (mouse enabled)
	app := tview.NewApplication().EnableMouse(true)

	// Define pages
	pages := tview.NewPages()

	// Configure task list page management
	taskCurrentPage := 0
	taskPageSize := cfg.Task.ListPageSize

	// Set up common components
	separator := tview.NewTextView().
		SetDynamicColors(true).
		SetText("------------------------------------------------------------------------------------------")

	// -- Task List Settings --
	taskDescription := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Please select the task you want to run.[-]")

	taskSelectList := tview.NewList()

	taskMenu := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(taskDescription, 2, 1, false).
		AddItem(separator, 1, 1, false).
		AddItem(taskSelectList, 18, 1, true).
		AddItem(separator, 1, 1, false).
		AddItem(nil, 0, 1, false)
	taskMenu.SetBorder(true).SetTitle(" Task list menu ")

	// -- Main Menu Settings --
	mainDescription := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Please select from the following menu.[-]")

	mainSelectList := tview.NewList().
		AddItem("[::b]・Retrieve the issue information and create or update task.md.[::-]", "", '1', func() {
			// Generating modal settings
			generating := tview.NewModal().SetText("Generating...")
			pages.AddPage("generating", generating, true, true)

			go func() {
				// Generate "task.md"
				err := generateTaskMd(cfg)

				// Screen update settings
				app.QueueUpdateDraw(func() {
					pages.RemovePage("generating")

					// Force redraw to fix UI corruption
					app.Sync()

					// In case of an error
					if err != nil {
						errorModal := tview.NewModal().
							SetText(fmt.Sprintf("[yellow][::b]An error occurred !![::-]\n\n%v", err)).
							AddButtons([]string{"OK"}).
							SetDoneFunc(func(buttonIndex int, buttonLabel string) {
								pages.RemovePage("error")
							})
						pages.AddPage("error", errorModal, true, true)
						return
					}

					// Success message
					successModal := tview.NewModal().
						SetText("task.md has been generated successfully !!").
						AddButtons([]string{"Close"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							pages.RemovePage("success")
						})
					pages.AddPage("success", successModal, true, true)
				})
			}()
		}).
		AddItem("[::b]・Read task.md and display the list of tasks.[::-]", "", '2', func() {
			// Initialize the current page of the task list to 0
			taskCurrentPage = 0

			// Load task information from task.md
			tasks, err := loadTaskMd()
			if err != nil {
				errorModal := tview.NewModal().
					SetText(fmt.Sprintf("[yellow][::b]An error occurred !![::-]\n\n%v", err)).
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(buttonIndex int, buttonLabel string) {
						pages.RemovePage("error")
					})
				pages.AddPage("error", errorModal, true, true)
				return
			}

			// Display the task list
			renderTasks(cfg, app, taskSelectList, pages, tasks, &taskCurrentPage, &taskPageSize)
			pages.SwitchToPage("task_menu")
		}).
		AddItem("Quit", "", 'q', func() {
			app.Stop()
		})

	mainMenu := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainDescription, 2, 1, false).
		AddItem(separator, 1, 1, false).
		AddItem(mainSelectList, 6, 1, true).
		AddItem(separator, 1, 1, false).
		AddItem(nil, 0, 1, false)
	mainMenu.SetBorder(true).SetTitle(" Main menu ")

	// -- Screen setup --
	pages.AddPage("main_menu", mainMenu, true, true)
	pages.AddPage("task_menu", taskMenu, true, false)

	// App startup process
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
