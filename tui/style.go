package tui

import (
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var (
	width, height, err = term.GetSize(0)
	viewportWidth      = termWidth - 5
	viewportHeight     = termHeight - 8
	termWidth          = int(float32(width) / 2)
	termHeight         = int(float32(height) / 1.5)

	yankedStyle = lipgloss.NewStyle().Background(lipgloss.Color("205"))
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#E95678"))
	appStyle    = lipgloss.NewStyle().
			Width(termWidth).
			Height(termHeight).
			Margin((height-termHeight)/2-2, 0, 0, (width-termWidth)/2).
			Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#E7E7E7"))
	viewportFocusedStyle = lipgloss.NewStyle().Background(lipgloss.Color("#E5E5F5"))
	notionRoleStyle      = lipgloss.NewStyle().Background(lipgloss.Color("#D5D5F5"))
	errorStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#EE3000"))
	helpStyle            = lipgloss.NewStyle().Margin(0, 0, 0, (width-termWidth)/2).Foreground(lipgloss.Color("#323253"))
)
