# Robby Gavigan - Terminal Resume

Welcome to my interactive terminal resume! This project is a fully functional Terminal User Interface (TUI) built in Go using the [Charm.sh](https://charm.sh/) ecosystem (`bubbletea`, `lipgloss`, and `wish`).

It allows anyone to SSH directly into my resume without installing any software.

## 🚀 How to View (Live)

You can interact with my resume right now directly from your terminal. No installation required!

Open your Mac or Linux terminal (or PowerShell on Windows) and run:

```bash
ssh -p 36741 guest@194.213.18.204
```
*(Note: If the terminal warns you about host authenticity, just type `yes`)*

### Navigation
Once connected, use your keyboard to navigate the TUI:
* `1` - View Experience
* `2` - View Technical Skills
* `3` - View Education (Themed for The Ohio State University! 🌰)
* `0` - Return to the main menu
* `q` or `Ctrl+C` - Disconnect

---

## 🛠️ How to Run Locally (For Developers)

If you want to run the code yourself or use this as a template for your own resume, you can easily build it using Go.

### Prerequisites
* [Go](https://golang.org/doc/install) (1.20 or later)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/gavigan4/tui-resume.git
   cd tui-resume
   ```
2. Download dependencies:
   ```bash
   go mod tidy
   ```
3. Start the SSH server:
   ```bash
   go run main.go
   ```

The server will start on port `23234`. You can test it locally by opening a second terminal window and running:
```bash
ssh -p 23234 localhost
```

## 🏗️ Technologies Used
* **[Go](https://go.dev/)**: Core application logic.
* **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: The Elm-architecture TUI framework.
* **[Lipgloss](https://github.com/charmbracelet/lipgloss)**: For the custom styling, colors, and layout (including the Ohio State theme!).
* **[Wish](https://github.com/charmbracelet/wish)**: The SSH server middleware that serves the Bubble Tea app over the network.
