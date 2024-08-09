package main

import "github.com/charmbracelet/lipgloss"

var (
	InactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = tabBorderWithBottom("┘", " ", "└")
	DocStyle          = lipgloss.NewStyle().Padding(0, 0, 0, 0)
	HighlightColor    = lipgloss.AdaptiveColor{Light: "#3f3f46", Dark: "#475569"}
	InactiveTabStyle  = lipgloss.NewStyle().Border(InactiveTabBorder, true).BorderForeground(HighlightColor).Padding(0, 1)
	ActiveTabStyle    = InactiveTabStyle.Border(ActiveTabBorder, true)
	WindowStyle       = lipgloss.NewStyle().BorderForeground(HighlightColor).Padding(2, 0).Align(lipgloss.Center, lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()
)
