# Expense Tracker CLI

An interactive, terminal-based expense tracker built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss). Quickly view, add, and persist your daily expenses—all from a polished, responsive TUI!

---

## 🖥️ Table of Contents

- [Features](#-features)  
- [Prerequisites](#-prerequisites)  
- [Installation](#-installation)  
- [Usage](#-usage)  
- [Configuration](#-configuration)  
- [Development](#-development)  
- [Contributing](#-contributing)  
- [License](#-license)  

---

## 🚀 Features

- **Rich, Responsive List View**  
  Navigate with ↑/↓, filter by typing, paginate, and toggle help overlays out of the box.
- **Quick‑Add Expense Overlay**  
  Press <kbd>Ctrl+N</kbd> to open a neat text‑entry prompt: `Name Amount Date` (e.g. `coffee 2.50 12-04-2023`).
- **Input Validation & Safe Cancel**  
  Rejects malformed entries (wrong fields, non‑numeric amount, invalid date) and restores the main list without error.
- **Progress Bar Animation**  
  After adding an expense, watch a smooth progress bar fill from 0→100% to confirm persistence.
- **Persistent Storage**  
  All entries are appended to `expenses.txt` (created automatically), one record per line, ideal for backups or CSV import.
- **Auto‑Resizing Layout**  
  List always fills your terminal; windows resizes are handled automatically.
- **Cross‑Platform**  
  Pure Go binary—works on Windows, macOS, and Linux.

---

## 🛠️ Prerequisites

- [Go 1.18+](https://golang.org/dl/) installed and in your `PATH`
- A Unix-like terminal or Windows PowerShell / Command Prompt

---

## 📥 Installation

#dowmload precompiled binary :)
```bash
# Clone this repository
git clone https://github.com/your-username/expense-tracker-cli.git
cd expense-tracker-cli

# Build for your current OS
go build -o expenses.exe main.go

# (Optional) Cross‑compile for Windows on macOS/Linux:
GOOS=windows GOARCH=amd64 go build -o expenses.exe main.go
