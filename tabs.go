package main

import "github.com/charmbracelet/lipgloss"

func RenderTabs(m model) string {
	var renderedTabs []string

	// Calculate total width available for tabs
	totalWidth := m.windowSize[0] - (len(m.Tabs) * 2) // Subtract space for borders between tabs

	// Calculate the width for each tab
	tabWidth := totalWidth / len(m.Tabs)

	// Adjust in case of rounding issues
	extraSpace := totalWidth % len(m.Tabs)

	for i, t := range m.Tabs {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab

		// Give extra space to some tabs if there's leftover width
		widthAdjustment := tabWidth
		if extraSpace > 0 {
			widthAdjustment++
			extraSpace--
		}

		if isActive {
			style = ActiveTabStyle.Width(widthAdjustment)
		} else {
			style = InactiveTabStyle.Width(widthAdjustment)
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

		style = style.Border(border)

		tabName := t.name
		if t.notifications {
			tabName = "● " + tabName
		}

		renderedTabs = append(renderedTabs, style.Render(tabName))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}
