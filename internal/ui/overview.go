package ui

import (
	"fmt"
	"github.com/thejawker/rennen/internal/types"
	"github.com/thejawker/rennen/internal/utils"
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
	processTable := renderProcessTable(m, maxLines-1, windowWidth)

	hint := hintStyle.Width(windowWidth).Render("←/→ tabs, (q)uit all")

	return fmt.Sprintf("%s\n%s\n%s", commandList, processTable, hint)
}

func renderCommandList(m types.ViewModelProvider, width int) interface{} {
	list := "Shortcuts:\n"
	for i, cmd := range m.GetViewModel().Commands {
		// use a space to separate the command from the description
		// to show selected you do > before the command
		prefix := " "
		if i == m.GetViewModel().SelectedCommand {
			prefix = "›"
		}

		status := "•"
		// we check whether the command is running or not and if it has been active within last 10 seconds
		if cmd.LastActivity.Add(3 * time.Second).After(time.Now()) {
			status = "✓"
		}

		list += fmt.Sprintf(" %s %s %s\n", prefix, status, cmd.Shortname)
	}

	return list
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

	return table.Render()
}
