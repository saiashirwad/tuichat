package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/chat"
	"github.com/saiashirwad/gochat/internal/config"
	"github.com/saiashirwad/gochat/internal/llm"
)

// ChatView displays the conversation history
type ChatView struct {
	config        *config.Config
	messages      []chat.Message
	llmClient     *llm.Client
	width, height int
	style         lipgloss.Style
}

// NewChatView creates a new chat view
func NewChatView(cfg *config.Config) *ChatView {
	return &ChatView{
		config:    cfg,
		llmClient: llm.NewClient(cfg),
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1),
		messages: []chat.Message{
			chat.NewMessage(chat.RoleAssistant, "Welcome to GoChat! Type your message below and press Enter to send."),
		},
	}
}

// SetSize updates the size of the chat view
func (c *ChatView) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.style = c.style.Width(width - 2).Height(height - 2)
}

// Init initializes the chat view
func (c *ChatView) Init() tea.Cmd {
	return nil
}

// sendMessageCmd creates a command to send a message to the LLM
func sendMessageCmd(client *llm.Client, messages []chat.Message) tea.Cmd {
	return func() tea.Msg {
		response, err := client.SendMessage(messages)
		if err != nil {
			return errMsg{err}
		}
		return newMessageMsg{
			message: chat.NewMessage(chat.RoleAssistant, response),
		}
	}
}

// Message types
type newMessageMsg struct {
	message chat.Message
}

type errMsg struct {
	err error
}

// Update handles events for the chat view
func (c *ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case newMessageMsg:
		c.messages = append(c.messages, msg.message)
		return c, nil
	case errMsg:
		c.messages = append(c.messages, chat.NewMessage(chat.RoleAssistant, fmt.Sprintf("Error: %v", msg.err)))
		return c, nil
	case userInputMsg:
		// Add user message to history
		userMessage := chat.NewMessage(chat.RoleUser, msg.input)
		c.messages = append(c.messages, userMessage)
		// Send to LLM
		return c, sendMessageCmd(c.llmClient, c.messages)
	}
	return c, nil
}

// formatMessage formats a single message for display
func (c *ChatView) formatMessage(msg chat.Message) string {
	var prefix string
	var msgStyle lipgloss.Style

	switch msg.Role {
	case chat.RoleUser:
		prefix = "You: "
		msgStyle = userMsgStyle
	case chat.RoleAssistant:
		prefix = "Assistant: "
		msgStyle = botMsgStyle
	default:
		prefix = "System: "
		msgStyle = lipgloss.NewStyle()
	}

	timestamp := ""
	if c.config.UI.ShowTimestamp {
		timestamp = msg.Timestamp.Format("15:04:05 ")
	}

	return msgStyle.Render(timestamp + prefix + msg.Content)
}

// View renders the chat view
func (c *ChatView) View() string {
	var formattedMessages []string
	for _, msg := range c.messages {
		formattedMessages = append(formattedMessages, c.formatMessage(msg))
	}

	content := strings.Join(formattedMessages, "\n\n")
	return c.style.Render(content)
}
