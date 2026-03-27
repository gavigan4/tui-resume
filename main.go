package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%s", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%s", host, port)
	go func() {
		if err = s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && err != ssh.ErrServerClosed {
		log.Fatalln(err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, active := s.Pty()
	if !active {
		wish.Fatalln(s, "no active terminal, skipping")
		return nil, nil
	}
	m := model{
		term:   pty.Term,
		width:  pty.Window.Width,
		height: pty.Window.Height,
		view:   0,
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type model struct {
	term   string
	width  int
	height int
	view   int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "0", "esc", "b", "left":
			m.view = 0
		case "1":
			m.view = 1
		case "2":
			m.view = 2
		case "3":
			m.view = 3
		}
	}
	return m, nil
}

// --- ADVANCED LIPGLOSS STYLING ---
var (
	primary   = lipgloss.Color("#58a6ff") // GitHub Blue
	secondary = lipgloss.Color("#39d353") // Hacker Green
	accent    = lipgloss.Color("#f2cc60") // Yellow
	textDim   = lipgloss.Color("#8b949e")

	scarlet   = lipgloss.Color("#CE0F3D") // OSU Scarlet
	gray      = lipgloss.Color("#B6C1C6") // OSU Gray
	leafGreen = lipgloss.Color("#4A773C")

	appBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primary).
		Padding(1, 4)

	headerStyle = lipgloss.NewStyle().
		Foreground(secondary).
		Bold(true).
		Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
		Foreground(primary).
		Align(lipgloss.Center)

	menuItemStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c9d1d9")).
		PaddingLeft(2)

	hotkeyStyle = lipgloss.NewStyle().
		Foreground(accent).
		Bold(true)

	sectionTitleStyle = lipgloss.NewStyle().
		Foreground(primary).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(primary).
		MarginBottom(1)

	jobTitleStyle = lipgloss.NewStyle().Foreground(secondary).Bold(true)
	companyStyle  = lipgloss.NewStyle().Foreground(primary).Bold(true)
	dateStyle     = lipgloss.NewStyle().Foreground(textDim).Italic(true)
	bulletStyle   = lipgloss.NewStyle().Foreground(textDim).SetString("• ")

	// OSU Specific Styling
	osuBox = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(scarlet).
		Padding(2, 6).
		Align(lipgloss.Center)

	osuHeader = lipgloss.NewStyle().
		Background(scarlet).
		Foreground(lipgloss.Color("#FFFFFF")).
		Bold(true).
		Padding(0, 2).
		MarginBottom(1)
)

func (m model) View() string {
	// Global Header for Views 0, 1, 2
	logo := headerStyle.Render(
		"___  ___  ___  ___ __ __   ___  ___  _ _  _  ___  ___  _ _\n" +
			"| . \\| . || . >| . >\\ \\ /  /  _>| . || | || |/  _>| . || \\ |\n" +
			"|   /| | || . \\| . \\ \\ /   | <__|   || ' || || <__|   ||   |\n" +
			"|_\\_\\`___'|___/|___/ |_|   `___/|_|_||__/ |_|`___/|_|_||_\\_|\n")
	subtitle := subtitleStyle.Render("SOFTWARE ENGINEER | FULLSTACK & AI")
	topBanner := lipgloss.JoinVertical(lipgloss.Center, logo, subtitle, "")

	var content string

	switch m.view {
	case 0:
		// MENU
		menu := lipgloss.JoinVertical(lipgloss.Left,
			"Select a section to explore:\n",
			menuItemStyle.Render(hotkeyStyle.Render("[1]")+" Experience"),
			menuItemStyle.Render(hotkeyStyle.Render("[2]")+" Skills"),
			menuItemStyle.Render(hotkeyStyle.Render("[3]")+" Education"),
			"",
			menuItemStyle.Render(hotkeyStyle.Render("[q]")+" Quit connection"),
		)
		content = lipgloss.JoinVertical(lipgloss.Center, topBanner, menu)

	case 1:
		// EXPERIENCE
		exp := lipgloss.JoinVertical(lipgloss.Left,
			sectionTitleStyle.Render("EXPERIENCE"),
			lipgloss.JoinHorizontal(lipgloss.Left, companyStyle.Render("JPMorgan Chase"), " | ", jobTitleStyle.Render("Software Engineer"), " ", dateStyle.Render("(Jun 2025 - Present)")),
			bulletStyle.Render()+"Architect and develop fullstack applications and robust internal tooling.\n",
			lipgloss.JoinHorizontal(lipgloss.Left, companyStyle.Render("Falconbridge Corporation"), " | ", jobTitleStyle.Render("Software Engineer"), " ", dateStyle.Render("(Jun 2025 - Present)")),
			bulletStyle.Render()+"Engineered a Retrieval-Augmented Generation (RAG) chatbot system.\n",
			lipgloss.JoinHorizontal(lipgloss.Left, companyStyle.Render("Abbott Nutrition"), " | ", jobTitleStyle.Render("Data Science Intern"), " ", dateStyle.Render("(May 2024 - Aug 2024)")),
			bulletStyle.Render()+"Deployed analytical models and automated ETL pipelines.",
			bulletStyle.Render()+"Directed national media strategies via data-driven marketing insights.\n",
			lipgloss.JoinHorizontal(lipgloss.Left, companyStyle.Render("S2R Analytics & BrewDog"), " | ", jobTitleStyle.Render("Data Analytics"), " ", dateStyle.Render("(2023)")),
			bulletStyle.Render()+"Streamlined API data retrieval with Python dashboards.",
			bulletStyle.Render()+"Analyzed market data using R and Qualtrics for C-level executives.",
			"",
			lipgloss.NewStyle().Foreground(textDim).Render("Press [0] to go back"),
		)
		content = lipgloss.JoinVertical(lipgloss.Left, topBanner, exp)

	case 2:
		// SKILLS
		skills := lipgloss.JoinVertical(lipgloss.Left,
			sectionTitleStyle.Render("TECHNICAL SKILLS"),
			lipgloss.JoinHorizontal(lipgloss.Left, hotkeyStyle.Render("[SYS]   "), "Python, R, SQL, Java, C++, Node.js, .NET"),
			lipgloss.JoinHorizontal(lipgloss.Left, hotkeyStyle.Render("[TOOLS] "), "VS Code, PyCharm, SQLite, Nielsen DB, Power BI"),
			lipgloss.JoinHorizontal(lipgloss.Left, hotkeyStyle.Render("[MISC]  "), "OAuth 2.0, APIs, Agile, Jira, GitHub"),
			"",
			lipgloss.NewStyle().Foreground(textDim).Render("Press [0] to go back"),
		)
		content = lipgloss.JoinVertical(lipgloss.Left, topBanner, skills)

	case 3:
		// EDUCATION - OHIO STATE THEMED
		blockOStr := lipgloss.NewStyle().Foreground(scarlet).Bold(true).Render(
			"  ██████╗  \n" +
			" ██╔═══██╗ \n" +
			" ██║   ██║ \n" +
			" ██║   ██║ \n" +
			" ╚██████╔╝ \n" +
			"  ╚═════╝  ")

		leafStr := lipgloss.NewStyle().Foreground(leafGreen).Bold(true).Render(
			"      .    \n" +
			"     / \\   \n" +
			"    /   \\  \n" +
			"   |\\___/| \n" +
			"   \\  |  / \n" +
			"    \\ | /  \n" +
			"     \\|/   \n" +
			"      |    ")

		// Join the ASCII art side-by-side
		artBox := lipgloss.JoinHorizontal(lipgloss.Top, blockOStr, "      ", leafStr)

		osuBody := lipgloss.JoinVertical(lipgloss.Center,
			artBox,
			"",
			osuHeader.Render(" THE OHIO STATE UNIVERSITY "),
			lipgloss.NewStyle().Foreground(gray).Bold(true).Render("B.S. Business Administration & CS Minor"),
			lipgloss.NewStyle().Foreground(gray).Render("December 2024"),
			"",
			lipgloss.NewStyle().Foreground(gray).Render("GPA: 3.5 | Dean's List | Info Systems Specialization"),
			"",
			lipgloss.NewStyle().Foreground(textDim).Render("Press [0] to go back"),
		)

		// We wrap the OSU content directly in the OSU double-bordered box instead of the default box
		content = osuBox.Render(osuBody)
	}

	// Apply the main app box for everything except the Education view (which handles its own box)
	var ui string
	if m.view == 3 {
		ui = content
	} else {
		ui = appBox.Render(content)
	}

	// Use lipgloss Place to perfectly center the UI in the user's terminal
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, ui)
}
