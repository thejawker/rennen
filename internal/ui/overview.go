package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/thejawker/rennen/internal/types"
	"github.com/thejawker/rennen/internal/utils"
	"strings"
	"time"
)

func renderOverview(m types.ViewModelProvider, maxLines int) string {
	//processes := "Processes: " + fmt.Sprintf("%d", len(m.GetViewModel().Processes))
	//runTime := "Started: " + m.GetRunTime()
	//
	//count := lipgloss.
	//	NewStyle().
	//	Foreground(lipgloss.AdaptiveColor{Light: "#555", Dark: "#555"}).
	//	Bold(true).
	//	Render(fmt.Sprintf("%s\n%s", processes, runTime))
	//
	//content := lipgloss.NewStyle().Height(maxLines - 1).Render(count)
	windowWidth := m.GetViewModel().WindowSize.Width - windowStyle.GetHorizontalFrameSize() - 2
	commandList := renderCommandList(m, windowWidth)
	commandLines := strings.Split(commandList, "\n")
	processTable := renderProcessTable(m, maxLines-len(commandLines)-2, windowWidth)

	hint := hintStyle.Width(windowWidth).Render("←/→ tabs, (q)uit all")

	return fmt.Sprintf("%s\n\n%s\n%s", commandList, processTable, hint)
}

func renderCommandList(m types.ViewModelProvider, width int) string {
	list := "shortcuts:\n"

	listStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderLeft(true).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false).
		PaddingLeft(1).
		BorderForeground(lipgloss.Color("253"))

	if len(m.GetViewModel().Commands) == 0 {
		list += lipgloss.NewStyle().
			Italic(true).
			Faint(true).
			PaddingTop(1).
			Render("consider adding a command so you can easily ren it\n")
		return listStyle.Render(list)
	}

	for i, cmd := range m.GetViewModel().Commands {
		style := lipgloss.NewStyle()

		prefix := " "
		if i == m.GetViewModel().SelectedCommand {
			style = style.Foreground(lipgloss.Color("#555"))
			prefix = "›"
		}

		// pink
		prefix = lipgloss.NewStyle().
			PaddingRight(1).
			PaddingRight(1).
			Foreground(lipgloss.Color("#ff00ff")).
			Render(prefix)

		status := "•"
		// we check whether the command is running or not and if it has been active within last 10 seconds
		if cmd.LastActivity.Add(3 * time.Second).After(time.Now()) {
			style = style.Foreground(lipgloss.Color("#15803d"))
			status = "✓"
		}

		list += prefix
		list += style.Render(fmt.Sprintf("%s %s", status, cmd.Shortname))
		list += "\n"
	}

	return listStyle.Render(list)
}

func renderProcessTable(m types.ViewModelProvider, height, width int) string {
	table := NewTable().
		SetColumns([]string{"command / process", "output", "status"}).
		SetColumnWidth("command / process", 20).
		SetColumnWidth("status", 10).
		SetTotalWidth(width - 2)

	for _, proc := range m.GetViewModel().Processes {
		status := ""

		if proc.StartedAt != nil {
			status = utils.RelativeTime(*proc.StartedAt)
		}

		if proc.IsStopped() {
			status = "Stopped"
		}

		table.AddRow([]string{proc.Shortname, proc.GetLastNonEmptyLine(), status})
	}

	return lipgloss.NewStyle().
		Height(height).
		Render(table.Render())
}
