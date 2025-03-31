package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/config"
)

// ChatView displays the conversation history
type ChatView struct {
	config        *config.Config
	messages      []string // This would be replaced with actual chat messages
	width, height int
	style         lipgloss.Style
}

// NewChatView creates a new chat view
func NewChatView(cfg *config.Config) *ChatView {
	return &ChatView{
		config: cfg,
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1),
		messages: []string{"Welcome to GoChat!", "Type your message below and press Enter to send."},
	}
}

// SetSize updates the size of the chat view
func (c *ChatView) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.style = c.style.Width(width - 2).Height(height - 2)
}

// Init initializes the chat view
func (c ChatView) Init() tea.Cmd {
	// No initial commands needed
	return nil
}

// Update handles events for the chat view
func (c *ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle messages and update state
	return c, nil
}

// View renders the chat view
func (c *ChatView) View() string {
	content := lipgloss.JoinVertical(lipgloss.Left, c.messages...)
	return c.style.Render(content)
}
