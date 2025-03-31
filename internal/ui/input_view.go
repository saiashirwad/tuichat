package ui

import (
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
	config *config.Config
	input  string
	cursor int
	width  int
	style  lipgloss.Style
}

// NewInputView creates a new input view
func NewInputView(cfg *config.Config) *InputView {
	return &InputView{
		config: cfg,
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")),
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
				// Create a message from the current input
				input := i.input
				// Clear the input
				i.input = ""
				i.cursor = 0
				// Return the message as a command
				return i, func() tea.Msg {
					return userInputMsg{input: input}
				}
			}
		case tea.KeyBackspace:
			if i.cursor > 0 {
				i.input = i.input[:i.cursor-1] + i.input[i.cursor:]
				i.cursor--
			}
		case tea.KeyRunes:
			i.input = i.input[:i.cursor] + string(msg.Runes) + i.input[i.cursor:]
			i.cursor += len(msg.Runes)
		case tea.KeyLeft:
			if i.cursor > 0 {
				i.cursor--
			}
		case tea.KeyRight:
			if i.cursor < len(i.input) {
				i.cursor++
			}
		}
	}
	return i, nil
}

// View renders the input view
func (i *InputView) View() string {
	prompt := "> "
	cursor := i.input
	if i.cursor < len(i.input) {
		cursor = i.input[:i.cursor] + "|" + i.input[i.cursor:]
	} else {
		cursor = i.input + "|"
	}
	return i.style.Render(prompt + cursor)
}
