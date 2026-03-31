package main

import (
"context"
"fmt"
"log"
"os"
"os/signal"
"syscall"
"time"

"github.com/charmbracelet/bubbles/help"
"github.com/charmbracelet/bubbles/key"
"github.com/charmbracelet/bubbles/spinner"
"github.com/charmbracelet/bubbles/viewport"
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

type tickMsg time.Time

func main() {
srv, err := wish.NewServer(
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
if err = srv.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
log.Fatalln(err)
}
}()

sig := <-done
log.Printf("Received signal: %s. Stopping server...", sig)
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil && err != ssh.ErrServerClosed {
log.Fatalln(err)
}
}

func teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
pty, _, active := sess.Pty()
if !active {
wish.Fatalln(sess, "no active terminal, skipping")
return nil, nil
}

sp := spinner.New()
sp.Spinner = spinner.Dot
sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

vp := viewport.New(pty.Window.Width-10, 10)
vp.SetContent(fmt.Sprintf("%s\n%s\n%s", 
"• JPMorgan Chase | Fullstack Engineer (New York, NY)",
"• Falconbridge | RAG Chatbot Developer (Columbus, OH)",
"• Abbott Nutrition | Data Science Intern (Columbus, OH)"))

m := model{
width:    pty.Window.Width,
height:   pty.Window.Height,
spinner:  sp,
viewport: vp,
help:     help.New(),
keys:     keys,
}
return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type keyMap struct {
Up    key.Binding
Down  key.Binding
Back  key.Binding
Quit  key.Binding
One   key.Binding
Two   key.Binding
Three key.Binding
}

var keys = keyMap{
Up:    key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "scroll up")),
Down:  key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "scroll down")),
Back:  key.NewBinding(key.WithKeys("esc", "0", "left"), key.WithHelp("esc/0", "back")),
Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
One:   key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "exp")),
Two:   key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "skills")),
Three: key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "edu")),
}

func (k keyMap) ShortHelp() []key.Binding {
return []key.Binding{k.One, k.Two, k.Three, k.Back, k.Quit}
}
func (k keyMap) FullHelp() [][]key.Binding { return [][]key.Binding{{k.One, k.Two, k.Three, k.Up, k.Down}, {k.Back, k.Quit}} }

type model struct {
width    int
height   int
view     int
frame    int
spinner  spinner.Model
viewport viewport.Model
help     help.Model
keys     keyMap
}

func (m model) Init() tea.Cmd { return tea.Batch(m.spinner.Tick, m.tick()) }
func (m model) tick() tea.Cmd { return tea.Tick(time.Second/10, func(t time.Time) tea.Msg { return tickMsg(t) }) }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
switch msg := msg.(type) {
case tea.WindowSizeMsg:
m.width, m.height = msg.Width, msg.Height
m.viewport.Width, m.viewport.Height = msg.Width-10, 10
case tickMsg:
m.frame++
return m, m.tick()
case spinner.TickMsg:
var cmd tea.Cmd
m.spinner, cmd = m.spinner.Update(msg)
return m, cmd
case tea.KeyMsg:
switch {
case key.Matches(msg, m.keys.Quit): return m, tea.Quit
case key.Matches(msg, m.keys.Back): m.view = 0
case key.Matches(msg, m.keys.One): m.view = 1
case key.Matches(msg, m.keys.Two): m.view = 2
case key.Matches(msg, m.keys.Three): m.view = 3
}
}
var cmd tea.Cmd
m.viewport, cmd = m.viewport.Update(msg)
return m, cmd
}

var (
primary   = lipgloss.Color("#58a6ff")
scarlet   = lipgloss.Color("#CE0F3D")
osuBox    = lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).BorderForeground(scarlet).Padding(1, 4)
headerStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
)

func (m model) View() string {
shimmerColors := []string{"#39d353", "#26a641", "#006d21", "#26a641"}
curColor := shimmerColors[m.frame%len(shimmerColors)]

logo := headerStyle.Copy().Foreground(lipgloss.Color(curColor)).Render(
"___  ___  ___  ___ __ __   ___  ___  _ _  _  ___  ___  _ _\n" +
"| . \\| . || . >| . >\\ \\ /  /  _>| . || | || |/  _>| . || \\ |\n" +
"|   /| | || . \\| . \\ \\ /   | <__|   || ' || || <__|   ||   |\n" +
"|_\\_\\`___'|___/|___/ |_|   `___/|_|_||__/ |_|`___/|_|_||_\\_|\n")

subtitle := lipgloss.NewStyle().Foreground(primary).Render("SOFTWARE ENGINEER | FULLSTACK & AI " + m.spinner.View())
header := lipgloss.JoinVertical(lipgloss.Center, logo, subtitle)

var body string
switch m.view {
case 0:
body = "\nWelcome! Select a section below to explore."
case 1:
body = lipgloss.JoinVertical(lipgloss.Left, "[ EXPERIENCE ] (Scrollable)", m.viewport.View())
case 2:
body = "[ SKILLS ]\n\nGo, Python, React, Docker, AI/ML, Prompt Engineering"
case 3:
leafCols := []string{"#4A773C", "#5DA44B", "#CE0F3D", "#5DA44B"}
leafCol := leafCols[(m.frame/2)%len(leafCols)]
blockO := lipgloss.NewStyle().Foreground(scarlet).Bold(true).Render("  #######  \n ######### \n###     ###\n###     ###\n ######### \n  #######  ")
leaf := lipgloss.NewStyle().Foreground(lipgloss.Color(leafCol)).Render("    \\|/    \n  -- * --  \n    /|\\    ")
body = osuBox.Render(lipgloss.JoinVertical(lipgloss.Center, "THE OHIO STATE UNIVERSITY", "", lipgloss.JoinHorizontal(lipgloss.Center, blockO, "  ", leaf)))
}

footer := "\n" + m.help.View(m.keys)
return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, header, body, footer))
}
