# Go-AIDD

![Go](https://img.shields.io/badge/Go-blue?logo=go&logoColor=white)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
[![X (formerly Twitter) URL](https://img.shields.io/twitter/url?url=https%3A%2F%2Fx.com%2Ftomoyuki65&label=%40tomoyuki65)](https://x.com/tomoyuki65)

English | [日本語](README.ja.md)  

Go-AIDD (AI-Driven Development) is a tool for managing and executing tasks for AI-driven development.  
It consolidates task information in a tabular format in `task.md` and allows you to view and select tasks in a TUI.  
Executed tasks can be processed using local AI tools.  
  
<br>
  
## Key Features

* **Task Aggregation**: Extract tasks with specific labels from sources like GitHub Issues and consolidate them in a tabular format in `task.md`.  
* **TUI Display**: Show a list of tasks in a terminal-based UI.  
* **Task Execution**: Execute tasks using local AI tools (e.g., Gemini CLI).  
* **Automatic PR Creation**: Automatically create a pull request after task completion.  
  
<br>
  
## Motivation
AI-driven development is expected to become a central approach in future software development.  
To work efficiently in AI-driven development, it is important to break down tasks during the design phase into a format that AI can easily process, and handle them appropriately using AI tools.  
  
This tool aggregates multiple divided AI tasks into `task.md`, allowing you to select a task and have it processed by AI immediately. By executing tasks in parallel, you can significantly accelerate your development speed.  
  
Additionally, by leveraging AI tools running locally, you can avoid unnecessary network costs and execute tasks securely.  
  
<br>
  
## Prerequisites
  
### Environment
* **Go**: 1.25.6+
* **Docker / Docker Compose**: Required
  
### Execution tools
* **Git**: Required
* **GitHub CLI**: Required
* **Local AI Tools**: Required
  * A Gemini CLI or similar tool that runs locally is required.
  
<br>
  
## Tools
  
* **staticcheck**: [dominikh/go-tools](https://github.com/dominikh/go-tools)  
  
<br>
  
## Libraries
  
* **TUI**: [rivo/tview](https://github.com/rivo/tview)  
* **Configuration:** [knadh/koanf](https://github.com/knadh/koanf)  
  
<br>
  
## Usage
  
### 1. Rename the configuration file
Run the following command to rename the configuration file.  
```
cp src/config.example.yml src/config.yml
```
  
<br>
  
### 2.　Edit the contents of 「config.yml」
・To specify which GitHub issues to extract, or to change the labels applied when creating a PR after completing a task, modify the issue.label value.  
```
issue:
  label: "AI DD"
```  
> ※ Make sure the label you set here exists in the target repository.  
  
<br>
  
・To extract tasks from GitHub or clone the target repository before executing tasks, modify the github.repository value.  
```
github:
  repository: "owner/repository-name"
```  
  
<br>
  
・To change the method used to clone the repository, modify the github.clone_type value.  
```
github:
  clone_type: "SSH"
```
  
<br>
  
・To change the branch to clone, modify the github.clone_branch value.  
```
github:
  clone_branch: "main"
```
  
<br>
  
・To create PRs as drafts after completing a task, set github.pr_draft to true.  
```
github:
  pr_draft: false
```  
  
<br>
  
・To change the AI tool used for task execution, modify the ai.type value.  
```
ai:
  type: "Gemini CLI"
```  
> ※ Only "Gemini CLI" has been tested. Other tools may not work properly.  
  
<br>
  
### 3. Start the app using make
The pre-built binary files are stored in `/src/bin`. Use the following make commands according to your OS to start the application.  
   
> ※ Only MacOS (ARM) has been tested.  
  
<br>
  
・For MacOS (ARM)  
```
make run-mac
```
  
<br>
  
・For Linux  
```
make run-linux
```
  
<br>
  
・For Windows  
```
make run-windows
```
  
<br>
  
### 4. Select and execute actions from the menu
After launching the app, a TUI menu will be displayed.  
The following options are available.  
  
#### 1. 「・Retrieve the issue information and create or update task.md.」
Selecting this option retrieves task information based on your configuration and consolidates it into `src/task.md`.  
  
> ※ You can also manually create the file. A sample file 「src/task.example.md」 is provided; rename it and adjust the layout to create your own file if needed.  
  
<br>
  
#### 2. 「・Read task.md and display the list of tasks.」
This option reads task information from `src/task.md` and displays a task list.  
Selecting a task shows its details, and you can execute it by selecting forms using the TAB key.  
  
> ※ Before executing tasks, make sure to clone the target repository under the work directory.  
  
<br>
  
#### 3. 「Quit」
This option exits the application.  
  
<br>
  