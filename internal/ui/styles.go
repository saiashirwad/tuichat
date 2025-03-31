package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Common styles
	appStyle = lipgloss.NewStyle().
		Margin(1, 2)

	// Message styles
	userMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")).
		Bold(true)

	botMsgStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("170"))

	// Input styles
	inputPromptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("63")).
		Bold(true)

	// Title styles
	titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		Bold(true)
) 