package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/thejawker/rennen/internal/types"
)

func renderOverview(m types.ViewModelProvider, maxLines int) string {
	processes := "Processes: " + fmt.Sprintf("%d", len(m.GetViewModel().Processes))
	runTime := "Started: " + m.GetRunTime()

	count := lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#555", Dark: "#555"}).
		Bold(true).
		Render(fmt.Sprintf("%s\n%s", processes, runTime))

	content := lipgloss.NewStyle().Height(maxLines - 1).Render(count)
	windowWidth := m.GetViewModel().WindowSize.Width - windowStyle.GetHorizontalFrameSize() - 2
	hint := hintStyle.Width(windowWidth).Render("←/→ tabs, (q)uit all")

	return fmt.Sprintf("%s\n%s", content, hint)
}
