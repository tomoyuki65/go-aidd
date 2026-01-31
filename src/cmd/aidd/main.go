package main

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	"github.com/tomoyuki65/go-aidd/internal/config"
	mt "github.com/tomoyuki65/go-aidd/internal/module/task"
)

// Set up common components
var separator *tview.TextView = tview.NewTextView().
	SetDynamicColors(true).
	SetText("------------------------------------------------------------------------------------------")

// Display task details
func showTaskDetail(cfg *config.Config, app *tview.Application, pages *tview.Pages, task mt.Task) {
	description := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Would you like to run this task ?[-]")

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
			// Task running modal settings
			taskRunningModal := tview.NewModal().SetText("Task running......")
			pages.AddPage("task_running_modal", taskRunningModal, true, true)

			go func() {
				// タスク実行
				err := mt.RunTask(cfg, task)

				// Screen update settings
				app.QueueUpdateDraw(func() {
					pages.RemovePage("task_running_modal")

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
						SetText("Task completed successfully !!").
						AddButtons([]string{"Close"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							pages.RemovePage("success")
							pages.RemovePage("task_detail")
						})
					pages.AddPage("success", successModal, true, true)
				})
			}()
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
func renderTasks(cfg *config.Config, app *tview.Application, taskList *tview.List, pages *tview.Pages, tasks []mt.Task, currentPage, pageSize *int) {
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
			generatingModal := tview.NewModal().SetText("Generating...")
			pages.AddPage("generating_modal", generatingModal, true, true)

			go func() {
				// Generate "task.md"
				err := mt.GenerateTaskMd(cfg)

				// Screen update settings
				app.QueueUpdateDraw(func() {
					pages.RemovePage("generating_modal")

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
			tasks, err := mt.LoadTaskMd()
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
