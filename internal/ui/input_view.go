package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/config"
)

var (
	// Style for the input box - no borders, minimal padding
	inputBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("233"))

	// Style for the input text
	inputTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))
)

// userInputMsg is sent when the user submits a message
type userInputMsg struct {
	input string
}

// InputView handles user input
type InputView struct {
	config    *config.Config
	input     string
	cursorPos int
	width     int
	focused   bool
}

// NewInputView creates a new input view
func NewInputView(cfg *config.Config) *InputView {
	return &InputView{
		config:  cfg,
		input:   "",
		focused: true,
	}
}

// SetWidth updates the width of the input view
func (i *InputView) SetWidth(width int) {
	i.width = width
	inputBoxStyle = inputBoxStyle.Width(width)
}

// Init initializes the input view
func (i *InputView) Init() tea.Cmd {
	return nil
}

// Update handles events for the input view
func (i *InputView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !i.focused {
		return i, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			i.focused = false
			return i, func() tea.Msg {
				return focusChatsMsg{}
			}
		case "enter":
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
		case "shift+enter":
			// Insert newline at cursor position
			i.input = i.input[:i.cursorPos] + "\n" + i.input[i.cursorPos:]
			i.cursorPos++
		case "backspace":
			if i.cursorPos > 0 {
				// Remove character before cursor
				i.input = i.input[:i.cursorPos-1] + i.input[i.cursorPos:]
				i.cursorPos--
			}
		case "left":
			if i.cursorPos > 0 {
				i.cursorPos--
			}
		case "right":
			if i.cursorPos < len(i.input) {
				i.cursorPos++
			}
		case " ":
			// Handle space explicitly
			i.input = i.input[:i.cursorPos] + " " + i.input[i.cursorPos:]
			i.cursorPos++
		default:
			if msg.Type == tea.KeyRunes {
				// Insert runes at cursor position
				i.input = i.input[:i.cursorPos] + string(msg.Runes) + i.input[i.cursorPos:]
				i.cursorPos += len(msg.Runes)
			}
		}
	}
	return i, nil
}

// View renders the input view
func (i *InputView) View() string {
	if !i.focused {
		return ""
	}

	// Add cursor to input
	displayText := i.input
	if i.cursorPos < len(displayText) {
		displayText = displayText[:i.cursorPos] + "│" + displayText[i.cursorPos:]
	} else {
		displayText = displayText + "│"
	}

	return inputBoxStyle.Render(
		inputTextStyle.Render(displayText),
	)
}

// Focus sets the input view as focused
func (i *InputView) Focus() {
	i.focused = true
}

// Blur removes focus from the input view
func (i *InputView) Blur() {
	i.focused = false
}
