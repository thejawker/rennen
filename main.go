package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"strings"
)

type process struct {
	shortname   string
	command     string
	description string
	output      string
}

type tab struct {
	name          string
	notifications bool
}

type model struct {
	Tabs       []tab
	activeTab  int
	windowSize []int
	processes  []process
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.windowSize = []int{msg.Width, msg.Height}
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
func (m model) View() string {
	doc := strings.Builder{}

	tabs := RenderTabs(m)
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
		header := commandStyle.Render(fmt.Sprintf("$ %s", process.command)) + "\n"
		header += descriptionStyle.Render(fmt.Sprintf("%s", process.description)) + "\n"
		divider := dividerStyle.Render(strings.Repeat("─", m.windowSize[0]-WindowStyle.GetHorizontalFrameSize()-2))
		windowContent := header + divider + "\n" + outputStyle.Render(process.output)

		// Render the command, description, and output with padding
		doc.WriteString(windowStyle.Padding(0, 1, 1, 1).Align(lipgloss.Top, lipgloss.Left).Render(windowContent))

	} else {
		doc.WriteString(windowStyle.Render("Not implemented yet"))
	}

	return DocStyle.Render(doc.String())
}

func main() {
	//tabs := []string{"Lip Gloss", "Blush", "Eye Shadow", "Mascara", "Foundation"}
	//tabContent := []string{"Lip Gloss Tab", "Blush Tab", "Eye Shadow Tab", "Mascara Tab", "Foundation Tab"}
	processes := []process{
		{
			shortname:   "fronend",
			command:     "yarn start",
			description: "Start the frontend server",
			output:      "bunch of frontend output here",
		},
		{
			shortname:   "mail",
			command:     "open http://localhost:8025 && mailhog",
			description: "Start the mailhog server and open the web interface",
			output:      "bunch of mailhog output here",
		},
		{
			shortname:   "stripe",
			command:     "yarn stripe:listen",
			description: "Start the stripe webhook listener",
			output:      "bunch of stripe output here",
		},
		{
			shortname:   "queue",
			command:     "php artisan queue:work",
			description: "Start the queue worker",
			output:      "bunch of queue output here",
		},
		{
			shortname:   "schedule",
			command:     "php artisan schedule:work",
			description: "Start the scheduler",
			output:      "bunch of schedule output here",
		},
	}

	tabs := generateTabs(processes)

	m := model{
		Tabs:       tabs,
		processes:  processes,
		activeTab:  0,
		windowSize: []int{0, 0},
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) getActiveTab() tab {
	return m.Tabs[m.activeTab]
}

func (m model) getActiveProcess() *process {
	// find by name
	for _, p := range m.processes {
		if p.shortname == m.getActiveTab().name {
			return &p
		}
	}

	return nil
}

func generateTabs(processes []process) []tab {
	var tabs []tab

	tabs = append(tabs, tab{
		name:          "all",
		notifications: false,
	})

	for idx, p := range processes {
		tabs = append(tabs, tab{
			name:          p.shortname,
			notifications: idx == 3,
		})
	}
	return tabs
}
