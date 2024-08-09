package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type process struct {
	Shortname    string `json:"shortname"`
	Command      string `json:"command"`
	Description  string `json:"description"`
	output       string
	cmd          *exec.Cmd
	lastActivity time.Time
}

type tab struct {
	name          string
	notifications bool
}

type model struct {
	Tabs       []tab
	activeTab  int
	windowSize []int
	processes  []*process
	mutex      sync.Mutex
	program    *tea.Program
}

type config struct {
	Processes []process `json:"processes"`
}

func (m model) Init() tea.Cmd {
	log.Println("Initializing model")
	return tea.Batch(m.startAllProcesses()...)
}

type outputMsg struct {
	process *process
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			m.clearNotification(m.activeTab)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			m.clearNotification(m.activeTab)
			return m, nil
		}

	case processStartedMsg:
		log.Printf("Process started: %s\n", msg.process.Shortname)
		return m, m.updateNotifications()

	case processErrorMsg:
		log.Printf("Error starting process %s: %v\n", msg.process.Shortname, msg.err)
		return m, nil

	case tea.WindowSizeMsg:
		m.windowSize = []int{msg.Width, msg.Height}
		return m, nil

	case outputMsg:
		return m, m.updateNotifications()
	}

	return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

//goland:noinspection ALL
func (m *model) View() string {
	doc := strings.Builder{}

	tabs := RenderTabs(*m)
	doc.WriteString(tabs)
	doc.WriteString("\n")

	windowStyle := WindowStyle.Width((m.windowSize[0] - WindowStyle.GetHorizontalFrameSize())).Height(m.windowSize[1] - ActiveTabStyle.GetVerticalFrameSize() - 2)

	process := m.getActiveProcess()
	if process != nil {
		// Define adaptive styles for header and output
		commandStyle := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#a7c7e7", Dark: "#8394a7"}).
			Bold(true)
		descriptionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#a7e7a7", Dark: "#8aa78a"})
		dividerStyle := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#d3d3d3", Dark: "#5c5c5c"}).
			PaddingTop(0).
			PaddingBottom(1)
		outputStyle := lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#606060", Dark: "#e0e0e0"})

		// Construct the window content with command, description, and output
		header := commandStyle.Render(fmt.Sprintf("$ %s", process.Command)) + "\n"
		header += descriptionStyle.Render(fmt.Sprintf("%s", process.Description)) + "\n"
		divider := dividerStyle.Render(strings.Repeat("─", m.windowSize[0]-WindowStyle.GetHorizontalFrameSize()-2))
		windowContent := header + divider + "\n" + outputStyle.Render(process.output)

		// Render the command, description, and output with padding
		doc.WriteString(windowStyle.Padding(0, 1, 1, 1).Align(lipgloss.Top, lipgloss.Left).Render(windowContent))

	} else {
		doc.WriteString(windowStyle.Render("No process selected"))
	}

	return DocStyle.Render(doc.String())
}

func (m *model) startProcess(p *process) tea.Cmd {
	return func() tea.Msg {
		log.Printf("Starting process: %s\n", p.Shortname)
		args := strings.Fields(p.Command)
		cmd := exec.Command(args[0], args[1:]...)
		p.cmd = cmd

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return processErrorMsg{process: p, err: err}
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return processErrorMsg{process: p, err: err}
		}

		if err := cmd.Start(); err != nil {
			return processErrorMsg{process: p, err: err}
		}

		go m.handleOutput(p, io.MultiReader(stdout, stderr))

		return processStartedMsg{process: p}
	}
}

func (m *model) startAllProcesses() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.processes))
	for i, p := range m.processes {
		cmds[i] = m.startProcess(p)
	}
	return cmds
}

func (m *model) handleOutput(p *process, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		m.mutex.Lock()
		p.output += line + "\n"
		p.lastActivity = time.Now()
		log.Printf("Process %s output: %s\n", p.Shortname, line)
		m.mutex.Unlock()

		// Send a message to trigger UI update
		m.program.Send(outputMsg{process: p})
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error reading output from process %s: %v\n", p.Shortname, err)
	}
}

func (m *model) updateNotifications() tea.Cmd {
	return func() tea.Msg {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		for i, t := range m.Tabs {
			if i == 0 {
				continue // Skip the "all" tab
			}
			for _, p := range m.processes {
				if p.Shortname == t.name {
					m.Tabs[i].notifications = time.Since(p.lastActivity) < time.Minute
					break
				}
			}
		}
		return nil
	}
}

func (m *model) clearNotification(tabIndex int) {
	if tabIndex > 0 && tabIndex < len(m.Tabs) {
		m.Tabs[tabIndex].notifications = false
	}
}

type processStartedMsg struct {
	process *process
}

type processErrorMsg struct {
	process *process
	err     error
}

func main() {
	// Set up logging to a file
	logFile, err := os.Create("debug.log")
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("Starting application")

	cfg, err := loadConfig("ren.json")
	if err != nil {
		log.Printf("Error loading config: %v\n", err)
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	processes := make([]*process, len(cfg.Processes))
	for i := range cfg.Processes {
		processes[i] = &cfg.Processes[i]
	}

	tabs := generateTabs(processes)

	m := &model{
		Tabs:       tabs,
		processes:  processes,
		activeTab:  0,
		windowSize: []int{0, 0},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	m.program = p // Store program reference in model

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			// duration .2s
			duration := 200 * time.Millisecond

			p.Send(tea.Tick(duration, func(t time.Time) tea.Msg {
				return outputMsg{}
			}))
		}
	}()

	if _, err := p.Run(); err != nil {
		log.Printf("Error running program: %v", err)
		fmt.Printf("Error running program: %v", err)
		return
	}
}

func (m model) getActiveTab() tab {
	return m.Tabs[m.activeTab]
}

func (m model) getActiveProcess() *process {
	// find by name
	for _, p := range m.processes {
		if p.Shortname == m.getActiveTab().name {
			return p
		}
	}

	return nil
}

func loadConfig(filename string) (config, error) {
	var cfg config
	file, err := os.Open(filename)
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	return cfg, err
}

func generateTabs(processes []*process) []tab {
	var tabs []tab

	tabs = append(tabs, tab{
		name:          "all",
		notifications: false,
	})

	for _, p := range processes {
		tabs = append(tabs, tab{
			name:          p.Shortname,
			notifications: false,
		})
	}
	return tabs
}
