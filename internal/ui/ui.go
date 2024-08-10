package ui

import (
	"fmt"
	"github.com/thejawker/rennen/internal/utils"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/thejawker/rennen/internal/types"
)

var (
	inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(0, 0, 0, 0)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#3f3f46", Dark: "#475569"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(0, 1, 0, 1).Align(lipgloss.Top, lipgloss.Left).Border(lipgloss.RoundedBorder()).UnsetBorderTop()
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func RenderView(m types.ViewModelProvider) string {
	doc := strings.Builder{}

	// Render tabs
	doc.WriteString(renderTabs(m.GetViewModel()))
	doc.WriteString("\n")

	// window style
	windowWidth := m.GetViewModel().WindowSize.Width - windowStyle.GetHorizontalFrameSize() + 2
	windowHeight := m.GetViewModel().WindowSize.Height - activeTabStyle.GetVerticalFrameSize() - 2

	// Render content
	content, shouldCenter := renderContent(m, windowHeight)

	ws := windowStyle.Width(windowWidth).Height(windowHeight)

	if shouldCenter {
		ws = ws.Align(lipgloss.Center, lipgloss.Center)
	}

	doc.WriteString(ws.Render(content))

	return doc.String()
}

func renderTabs(vm types.Model) string {
	var renderedTabs []string

	// Calculate total width available for tabs
	totalWidth := vm.WindowSize.Width - (len(vm.Tabs) * 2) // Subtract space for borders between tabs

	// Calculate the width for each tab
	tabWidth := totalWidth / len(vm.Tabs)

	// Adjust in case of rounding issues
	extraSpace := totalWidth % len(vm.Tabs)

	for i, t := range vm.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(vm.Tabs)-1, i == vm.ActiveTab

		// Give extra space to some tabs if there's leftover width
		widthAdjustment := tabWidth
		if extraSpace > 0 {
			widthAdjustment++
			extraSpace--
		}

		if isActive {
			style = activeTabStyle.Width(widthAdjustment)
		} else {
			style = inactiveTabStyle.Width(widthAdjustment)
		}

		border, _, _, _, _ := style.GetBorder()

		if isFirst {
			if isActive {
				border.BottomLeft = "│"
			} else {
				border.BottomLeft = "├"
			}
		}

		if isLast {
			if isActive {
				border.BottomRight = "│"
			} else {
				border.BottomRight = "┤"
			}
		}

		style = style.
			Border(border)

		tabName := t.Name
		if t.Notification {
			tabName = "● " + tabName
		}

		truncatedName := utils.SmartTruncate(tabName, widthAdjustment-2, "")

		renderedTabs = append(renderedTabs, style.Render(truncatedName))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func renderContent(m types.ViewModelProvider, maxLines int) (string, bool) {
	process := m.GetActiveProcess()
	if process == nil {
		return fmt.Sprintf("Viewing tab: %s", m.GetActiveTabName()), true
	}

	// Define adaptive styles for header and output
	commandStyle := lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#a7c7e7", Dark: "#8394a7"}).
		Bold(true)
	descriptionStyle := lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#a7e7a7", Dark: "#8aa78a"})
	dividerStyle := lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#d3d3d3", Dark: "#5c5c5c"}).
		PaddingTop(0).
		PaddingBottom(1)
	outputStyle := lipgloss.
		NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#606060", Dark: "#e0e0e0"})

	// Construct the window content with command, description, and output
	header := commandStyle.Render(fmt.Sprintf("$ %s", process.Command)) + "\n"
	header += descriptionStyle.Render(fmt.Sprintf("%s", process.Description)) + "\n"

	windowWidth := m.GetViewModel().WindowSize.Width - windowStyle.GetHorizontalFrameSize() - 2
	divider := dividerStyle.Render(strings.Repeat("─", windowWidth))
	output := process.GetOutput()

	if output == "" {
		output = "No output yet..."
	}

	// Calculate available height for viewport
	headerHeight := lipgloss.Height(header)
	dividerHeight := lipgloss.Height(divider)
	viewportHeight := maxLines - headerHeight - dividerHeight

	// Create a viewport for scrollable content
	vp := viewport.New(windowWidth, viewportHeight)
	vp.SetContent(outputStyle.Render(output))
	vp.GotoBottom()

	// Combine all elements
	content := fmt.Sprintf("%s\n%s\n%s", header, divider, vp.View())

	return content, false
}
