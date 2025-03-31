package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
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
	config    *config.Config
	textInput textinput.Model
	width     int
}

// NewInputView creates a new input view
func NewInputView(cfg *config.Config) *InputView {
	ti := textinput.New()
	ti.Placeholder = "Type your message and press Enter..."
	ti.CharLimit = 4096 // Reasonable limit for LLM context
	ti.Width = 40       // Will be adjusted by SetWidth
	ti.Focus()

	// Style the text input
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.PromptStyle = lipgloss.NewStyle().Background(lipgloss.Color("233"))
	ti.PlaceholderStyle = ti.TextStyle.Copy().Foreground(lipgloss.Color("240"))

	return &InputView{
		config:    cfg,
		textInput: ti,
	}
}

// SetWidth updates the width of the input view
func (i *InputView) SetWidth(width int) {
	i.width = width
	i.textInput.Width = width
}

// Init initializes the input view
func (i *InputView) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles events for the input view
func (i *InputView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return i, func() tea.Msg {
				return focusChatsMsg{}
			}
		case tea.KeyEnter:
			if input := strings.TrimSpace(i.textInput.Value()); input != "" {
				oldInput := input
				i.textInput.Reset()
				return i, func() tea.Msg {
					return userInputMsg{input: oldInput}
				}
			}
		}
	}

	i.textInput, cmd = i.textInput.Update(msg)
	return i, cmd
}

// View renders the input view
func (i *InputView) View() string {
	return i.textInput.View()
}

// Focus sets the input view as focused
func (i *InputView) Focus() {
	i.textInput.Focus()
}

// Blur removes focus from the input view
func (i *InputView) Blur() {
	i.textInput.Blur()
}
