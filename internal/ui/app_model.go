package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/config"
)

// AppModel is the main application model
type AppModel struct {
	config        *config.Config
	chatView      *ChatView
	inputView     *InputView
	finderActive  bool
	finderView    *FinderView
	width, height int
}

// NewAppModel creates a new instance of the application model
func NewAppModel(cfg *config.Config) *AppModel {
	return &AppModel{
		config:       cfg,
		chatView:     NewChatView(cfg),
		inputView:    NewInputView(cfg),
		finderActive: false,
		finderView:   NewFinderView(cfg),
	}
}

// Init initializes the model
func (m *AppModel) Init() tea.Cmd {
	return nil
}

// Update handles events and updates the model
func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+f":
			// Toggle finder
			m.finderActive = !m.finderActive
			if m.finderActive {
				// Initialize finder search
				return m, m.finderView.Init()
			}
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		// Calculate heights
		inputHeight := 1                     // Input box height (just content, no borders)
		chatHeight := m.height - inputHeight // No extra space needed
		if chatHeight < 5 {
			chatHeight = 5 // Minimum chat height
		}

		// Update sub-component sizes
		m.chatView.SetSize(msg.Width, chatHeight)
		m.inputView.SetWidth(msg.Width)
		m.finderView.SetSize(msg.Width, msg.Height)
	}

	// Handle updates for sub-components
	if m.finderActive {
		// Update finder
		newFinderModel, cmd := m.finderView.Update(msg)
		if newModel, ok := newFinderModel.(*FinderView); ok {
			m.finderView = newModel
		}
		cmds = append(cmds, cmd)
	} else {
		// Update chat view
		newChatModel, cmd := m.chatView.Update(msg)
		if newModel, ok := newChatModel.(*ChatView); ok {
			m.chatView = newModel
		}
		cmds = append(cmds, cmd)

		// Update input view
		newInputModel, cmd := m.inputView.Update(msg)
		if newModel, ok := newInputModel.(*InputView); ok {
			m.inputView = newModel
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m *AppModel) View() string {
	if m.finderActive {
		return m.finderView.View()
	}

	// Join views without extra spacing
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.chatView.View(),
		m.inputView.View(),
	)
}
