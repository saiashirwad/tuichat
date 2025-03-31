package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/config"
)

// userInputMsg is sent when the user submits a message
type userInputMsg struct {
	input string
}

// InputView handles user input
type InputView struct {
	config     *config.Config
	input      string
	cursorPos  int
	width      int
	style      lipgloss.Style
	inputStyle lipgloss.Style
}

// NewInputView creates a new input view
func NewInputView(cfg *config.Config) *InputView {
	inputStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	return &InputView{
		config: cfg,
		input:  "",
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1),
		inputStyle: inputStyle,
	}
}

// SetWidth updates the width of the input view
func (i *InputView) SetWidth(width int) {
	i.width = width
	i.style = i.style.Width(width - 2)
}

// Init initializes the input view
func (i *InputView) Init() tea.Cmd {
	return nil
}

// Update handles events for the input view
func (i *InputView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if i.input != "" {
				input := strings.TrimSpace(i.input)
				if input != "" {
					// Clear the input
					oldInput := i.input
					i.input = ""
					i.cursorPos = 0
					// Return the message as a command
					return i, func() tea.Msg {
						return userInputMsg{input: oldInput}
					}
				}
			}
			return i, nil
		case tea.KeyBackspace:
			if i.cursorPos > 0 {
				// Remove character before cursor
				i.input = i.input[:i.cursorPos-1] + i.input[i.cursorPos:]
				i.cursorPos--
			}
		case tea.KeyLeft:
			if i.cursorPos > 0 {
				i.cursorPos--
			}
		case tea.KeyRight:
			if i.cursorPos < len(i.input) {
				i.cursorPos++
			}
		case tea.KeySpace:
			// Insert space at cursor position
			i.input = i.input[:i.cursorPos] + " " + i.input[i.cursorPos:]
			i.cursorPos++
		case tea.KeyRunes:
			// Insert runes at cursor position
			i.input = i.input[:i.cursorPos] + string(msg.Runes) + i.input[i.cursorPos:]
			i.cursorPos += len(msg.Runes)
		}
	}
	return i, nil
}

// View renders the input view
func (i *InputView) View() string {
	// Add cursor to input
	displayText := i.input
	if i.cursorPos < len(displayText) {
		displayText = displayText[:i.cursorPos] + "│" + displayText[i.cursorPos:]
	} else {
		displayText = displayText + "│"
	}

	return i.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			"Type your message (Enter to send):",
			i.inputStyle.Render(displayText),
		),
	)
}
