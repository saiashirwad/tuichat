package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/chat"
	"github.com/saiashirwad/gochat/internal/config"
	"github.com/saiashirwad/gochat/internal/llm"
)

var (
	// Style for the entire chat area - no borders
	chatStyle = lipgloss.NewStyle()

	// Base message style - minimal with just a separator line
	baseMessageStyle = lipgloss.NewStyle().
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder())

	// Style for user messages - pink separator
	userMessageStyle = baseMessageStyle.Copy().
				BorderForeground(lipgloss.Color("205"))

	// Style for LLM messages - blue separator
	llmMessageStyle = baseMessageStyle.Copy().
			BorderForeground(lipgloss.Color("39"))

	// Style for focused message - highlighted separator
	focusedMessageStyle = baseMessageStyle.Copy().
				BorderForeground(lipgloss.Color("99"))

	// Header styles - minimal
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	// Markdown renderer
	markdownRenderer *glamour.TermRenderer
)

func init() {
	// Initialize markdown renderer with dark theme
	var err error
	markdownRenderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0), // Will be set dynamically based on width
	)
	if err != nil {
		panic(err)
	}
}

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
	c.viewport.Width = width
	c.viewport.Height = height
	chatStyle = chatStyle.Width(width)

	// Calculate message width to fill the viewport
	messageWidth := width
	baseMessageStyle = baseMessageStyle.Width(messageWidth)
	userMessageStyle = userMessageStyle.Width(messageWidth)
	llmMessageStyle = llmMessageStyle.Width(messageWidth)
	focusedMessageStyle = focusedMessageStyle.Width(messageWidth)
	headerStyle = headerStyle.Width(messageWidth)

	// Update markdown renderer with new width
	markdownRenderer, _ = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(messageWidth-2), // Account for minimal padding
	)

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

// Message type for focusing chats
type focusChatsMsg struct{}

// Update handles events for the chat view
func (c *ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.focusActive {
			switch msg.String() {
			case "j", "tab":
				c.focusIndex++
				if c.focusIndex >= len(c.messages) {
					c.focusIndex = 0
				}
				c.updateContent()
			case "k", "shift+tab":
				c.focusIndex--
				if c.focusIndex < 0 {
					c.focusIndex = len(c.messages) - 1
				}
				c.updateContent()
			case "esc":
				c.focusActive = false
				c.updateContent()
			}
		}

	case newMessageMsg:
		c.messages = append(c.messages, msg.message)
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
		c.updateContent()
		// Send to LLM
		return c, sendMessageCmd(c.llmClient, c.messages)
	case focusChatsMsg:
		c.focusActive = true
		c.focusIndex = len(c.messages) - 1
		c.updateContent()
	}

	// Handle viewport messages
	c.viewport, cmd = c.viewport.Update(msg)
	return c, cmd
}

// updateContent updates the viewport content with formatted messages
func (c *ChatView) updateContent() {
	var formattedMessages []string
	var totalHeight int

	// First pass: format messages and calculate heights
	messageHeights := make([]int, len(c.messages))
	for i, msg := range c.messages {
		var content string
		var style lipgloss.Style

		// Render content as markdown
		rendered, err := markdownRenderer.Render(msg.Content)
		if err != nil {
			rendered = msg.Content // Fallback to plain text if markdown rendering fails
		}
		rendered = strings.TrimSpace(rendered) // Remove extra newlines from glamour

		// Add header based on role
		header := "LLM Message"
		if msg.Role == chat.RoleUser {
			header = "My message"
		}
		header = headerStyle.Render(header)

		// Join header and content without padding
		content = lipgloss.JoinVertical(lipgloss.Left, header, rendered)

		// Apply appropriate style based on role and focus
		if c.focusActive && i == c.focusIndex {
			style = focusedMessageStyle
		} else if msg.Role == chat.RoleUser {
			style = userMessageStyle
		} else {
			style = llmMessageStyle
		}

		// Apply the style and add to messages
		formattedMsg := style.Render(content)
		formattedMessages = append(formattedMessages, formattedMsg)

		// Calculate height of this message (count newlines + 1)
		height := strings.Count(formattedMsg, "\n") + 1
		messageHeights[i] = height
		totalHeight += height
	}

	// Join messages and set content
	content := strings.Join(formattedMessages, "\n")
	c.viewport.SetContent(content)

	// Adjust scrolling only when necessary
	if c.focusActive {
		// Calculate the position of the focused message
		var focusedMsgTop int
		for i := 0; i < c.focusIndex; i++ {
			focusedMsgTop += messageHeights[i]
		}
		focusedMsgBottom := focusedMsgTop + messageHeights[c.focusIndex]

		// Only scroll if the focused message is not fully visible
		if focusedMsgBottom > c.viewport.YOffset+c.viewport.Height {
			// Message is below viewport - scroll down
			c.viewport.YOffset = focusedMsgBottom - c.viewport.Height
		} else if focusedMsgTop < c.viewport.YOffset {
			// Message is above viewport - scroll up
			c.viewport.YOffset = focusedMsgTop
		}
	} else {
		// When not focused, stay at bottom for new messages
		if totalHeight > c.viewport.Height {
			c.viewport.YOffset = totalHeight - c.viewport.Height
		}
	}
}

// View renders the chat view
func (c *ChatView) View() string {
	return chatStyle.Render(c.viewport.View())
}
