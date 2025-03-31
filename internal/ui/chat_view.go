package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/chat"
	"github.com/saiashirwad/gochat/internal/config"
	"github.com/saiashirwad/gochat/internal/llm"
)

var (
	// Style for the entire chat area - no borders
	chatStyle = lipgloss.NewStyle()

	// Style for message boxes - solid background
	messageBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Background(lipgloss.Color("233")).
			Margin(0, 0, 0, 0). // No margin
			Padding(0, 1)

	// Style for focused message box - solid background
	focusedMessageBoxStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("205")).
				Background(lipgloss.Color("234")).
				Margin(0, 0, 0, 0). // No margin
				Padding(0, 1)

	// Style for user messages
	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	// Style for assistant messages
	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))
)

// ChatView displays the conversation history
type ChatView struct {
	config      *config.Config
	messages    []chat.Message
	llmClient   *llm.Client
	viewport    viewport.Model
	width       int
	height      int
	focusIndex  int  // Index of currently focused message
	focusActive bool // Whether message focus is active
}

// NewChatView creates a new chat view
func NewChatView(cfg *config.Config) *ChatView {
	c := &ChatView{
		config:    cfg,
		llmClient: llm.NewClient(cfg),
		messages: []chat.Message{
			chat.NewMessage(chat.RoleAssistant, "Welcome to GoChat! Type your message below and press Enter to send."),
		},
	}

	// Initialize viewport with minimum size
	c.viewport = viewport.New(10, 10)
	c.viewport.Style = lipgloss.NewStyle()

	return c
}

// SetSize updates the size of the chat view
func (c *ChatView) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.viewport.Width = width   // No need to subtract borders
	c.viewport.Height = height // No need to subtract borders
	chatStyle = chatStyle.Width(width)
	messageBoxStyle = messageBoxStyle.Width(width - 4) // Account for message box margins
	focusedMessageBoxStyle = focusedMessageBoxStyle.Width(width - 4)

	// Update content after resize
	c.updateContent()
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if !c.focusActive {
				c.focusActive = true
				c.focusIndex = len(c.messages) - 1
			} else {
				c.focusIndex++
				if c.focusIndex >= len(c.messages) {
					c.focusIndex = 0
				}
			}
			c.updateContent()
		case "shift+tab":
			if !c.focusActive {
				c.focusActive = true
				c.focusIndex = len(c.messages) - 1
			} else {
				c.focusIndex--
				if c.focusIndex < 0 {
					c.focusIndex = len(c.messages) - 1
				}
			}
			c.updateContent()
		case "esc":
			if c.focusActive {
				c.focusActive = false
				c.updateContent()
			}
		}

	case newMessageMsg:
		c.messages = append(c.messages, msg.message)
		if c.focusActive {
			c.focusIndex = len(c.messages) - 1
		}
		c.updateContent()
		return c, nil
	case errMsg:
		c.messages = append(c.messages, chat.NewMessage(chat.RoleAssistant, fmt.Sprintf("Error: %v", msg.err)))
		c.updateContent()
		return c, nil
	case userInputMsg:
		// Add user message to history
		userMessage := chat.NewMessage(chat.RoleUser, msg.input)
		c.messages = append(c.messages, userMessage)
		if c.focusActive {
			c.focusIndex = len(c.messages) - 1
		}
		c.updateContent()
		// Send to LLM
		return c, sendMessageCmd(c.llmClient, c.messages)
	}

	// Handle viewport messages
	c.viewport, cmd = c.viewport.Update(msg)
	return c, cmd
}

// updateContent updates the viewport content with formatted messages
func (c *ChatView) updateContent() {
	var formattedMessages []string

	for i, msg := range c.messages {
		var content string
		var style lipgloss.Style

		// Choose message style based on role
		switch msg.Role {
		case chat.RoleUser:
			content = userStyle.Render(msg.Content)
		case chat.RoleAssistant:
			content = assistantStyle.Render(msg.Content)
		}

		// Add header based on role
		header := "LLM Message"
		if msg.Role == chat.RoleUser {
			header = "My message"
		}
		content = lipgloss.JoinVertical(lipgloss.Left, header, content)

		// Apply message box style based on focus
		if c.focusActive && i == c.focusIndex {
			style = focusedMessageBoxStyle
		} else {
			style = messageBoxStyle
		}

		formattedMessages = append(formattedMessages, style.Render(content))
	}

	content := strings.Join(formattedMessages, "\n")
	c.viewport.SetContent(content)

	// Scroll to focused message if focus is active
	if c.focusActive {
		// Calculate approximate position of focused message
		var pos int
		for i := 0; i < c.focusIndex; i++ {
			pos += 5 // Approximate height of each message box
		}
		c.viewport.YOffset = pos
	} else {
		// Otherwise scroll to bottom for new messages
		c.viewport.GotoBottom()
	}
}

// View renders the chat view
func (c *ChatView) View() string {
	return chatStyle.Render(c.viewport.View())
}
