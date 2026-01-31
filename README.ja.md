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
  
### 1. コンフィグ設定用のファイルをリネーム
以下のコマンドを実行し、コンフィグ設定用のファイルをリネームします。  
```
cp src/config.example.yml src/config.yml
```
  
<br>
  
### 2.　「config.yml」の内容を修正
・GitHubから抽出する対象のIssueを指定する際や、タスク完了後のPR作成時に付与するラベル設定を変更したい場合は、issue.labelの値を修正して下さい。  
```
issue:
  label: "AI DD"
```  
> ※ ここで設定したラベルを対象リポジトリに作成して下さい。  
  
<br>
  
・GitHubからタスクを抽出する際や、タスク実行前に対象のリポジトリをクローンするため、`github.repository`の値を修正して下さい。  
```
github:
  repository: "オーナー名/リポジトリ名"
```  
  
<br>
  
・リポジトリをクローンする方法を変更したい場合はgithub.clone_typeの値を修正して下さい。  
```
github:
  clone_type: "SSH"
```
  
<br>
  
・リポジトリをクローンする際のブランチを変更したい場合はgithub.clone_branchの値を修正して下さい。  
```
github:
  clone_branch: "main"
```
  
<br>
  
・タスク完了後のPRをドラフトで作りたい場合は、github.pr_draftの値をtrueに変更して下さい。  
```
github:
  pr_draft: false
```  
  
<br>
  
・タスク実行時に利用するAIツールを変更したい場合は、ai.typeの値を修正して下さい。  
```
ai:
  type: "Gemini CLI"
```  
> ※ ただし、Gemini CLI以外は動作未検証です。  
  
<br>
  
### 3. makeコマンドでアプリ起動
ビルド済みのバイナリファイルを「/src/bin」に格納しているため、OSに合わせて以下のmakeコマンドを利用してアプリを起動して下さい。  
   
> ※ ただし、MacOS（ARM）以外は動作未検証です。  
  
<br>
  
・MacOS（ARM）の場合
```
make run-mac
```
  
<br>
  
・Linuxの場合
```
make run-linux
```
  
<br>
  
・Windowsの場合
```
make run-windows
```
  
<br>
  
### 4. メニューから処理を選択して実行
アプリを起動するとTUIでメニューが表示されます。   
以下のメニューがあり、それぞれ実行できます。  
  
#### 1. 「・Retrieve the issue information and create or update task.md.」
このメニューを選択するとコンフィグ設定の内容をもとにたタスク情報を取得し、`src/task.md`にタスク情報を集約します。  
  
> ※ タスク内容に画像が含まれていた場合、画像ファイルを「src/images」に.png形式でダウンロードします。そしてタスク内容に含まれている画像のURLをダウンロードした画像のファイルのパスに変換しています。  
  
> ※ タスク情報を集約するためのファイル「src/task.md」は手動で作っても大丈夫です。手動で作りたい場合はサンプルファイル「src/task.example.md」を格納しているため、ファイル名をリネーム後、中身のレイアウトを合わせてファイルを作成して下さい。  
  
<br>
  
#### 2. 「・Read task.md and display the list of tasks.」
このメニューを選択すると`src/task.md`からタスク情報を読み込んでタスク一覧を表示します。  
対象のタスクを選択するとタスクの詳細が表示され、TABキーでフォームを選択してタスクを実行できます。  
  
> ※ タスクを実行する際は、事前に対象のリポジトリをworkディレクトリ配下にクローンしてからタスクを実行するようにしています。  
  
<br>
  
#### 3. 「Quit」
このメニューを選択するとアプリを終了します。
  
<br>
  