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
  
<br>
  