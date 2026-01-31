# Go-AIDD

![Go](https://img.shields.io/badge/Go-blue?logo=go&logoColor=white)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
[![X (formerly Twitter) URL](https://img.shields.io/twitter/url?url=https%3A%2F%2Fx.com%2Ftomoyuki65&label=%40tomoyuki65)](https://x.com/tomoyuki65)

[English](README.md) | 日本語  

Go-AIDD (AI-Driven Development) は、AI駆動開発用のタスクを管理・実行するためのツールです。  
タスク情報を`task.md`に表形式で集約し、TUI上で一覧表示・選択実行できます。  
実行したタスクはローカルのAIツールを利用して処理できます。  
  
<br>
  
## 主な機能

* **タスク情報集約**: GitHub Issue等から特定のラベルが付いたタスクを抽出し、`task.md` へ表形式で集約。  
* **TUI表示**: ターミナルUIでタスク一覧を表示。  
* **タスク実行**: ローカルで使用中のAIツール（Gemini CLIなど）を利用してタスク実行。  
* **PR自動作成**: タスク完了後、プルリクエストを自動作成。  
  
<br>
  
## 開発の背景
今後のソフトウェア開発では、**「AI駆動開発」** が中心的な手法となることが予想されます。  
効率的にAI駆動開発を進めるには、設計段階でタスクをAIによる自動処理がしやすい形に分割し、それらを適切にAIツールで処理することが重要です。  
  
本ツールでは、分割した複数のAIタスクを`task.md`に集約し、対象のタスクを選択するだけでAIに処理させることができます。さらに、タスクを並列実行することで、開発スピードを大幅に向上させることが可能です。  
  
また、ローカル環境で利用中のAIツールを活用するため、余分なネットワークコストがかからず、セキュリティ面でも安心してタスクを実行できます。  
  
<br>
  
## 前提条件

### 開発環境
* **Go**: 1.25.6+
* **Docker / Docker Compose**: Required
  
### 実行用ツール
* **Git**: Required
* **GitHub CLI**: Required
* **ローカルAIツール**: Required
  * ローカルで動作する Gemini CLI 等が必要です。
  
  <br>
  
## 開発用ツール

* **staticcheck**: [dominikh/go-tools](https://github.com/dominikh/go-tools)  
  
<br>
  
## 使用ライブラリ
  
* **TUI**: [rivo/tview](https://github.com/rivo/tview)  
* **Configuration:** [knadh/koanf](https://github.com/knadh/koanf)  
  
<br>
  
## 使い方
  
<br>
  