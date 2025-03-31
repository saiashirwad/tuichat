package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/saiashirwad/gochat/internal/config"
)

// FinderView provides fuzzy search for chat history
type FinderView struct {
	config        *config.Config
	query         string
	results       []string // This would be replaced with actual chat results
	cursor        int
	width, height int
	style         lipgloss.Style
}

// NewFinderView creates a new finder view
func NewFinderView(cfg *config.Config) *FinderView {
	return &FinderView{
		config: cfg,
		style: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("170")).
			Padding(1),
		results: []string{"Chat 1", "Chat 2", "Chat 3"}, // Placeholder results
	}
}

// SetSize updates the size of the finder view
func (f *FinderView) SetSize(width, height int) {
	f.width = width
	f.height = height
	f.style = f.style.Width(width - 2).Height(height - 2)
}

// Init initializes the finder view
func (f *FinderView) Init() tea.Cmd {
	// In a real implementation, this would load chat history
	return nil
}

// Update handles events for the finder view
func (f *FinderView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if f.cursor > 0 {
				f.cursor--
			}
		case "down":
			if f.cursor < len(f.results)-1 {
				f.cursor++
			}
		case "enter":
			// Select the current chat
			// Would dispatch a command to load the selected chat
		}
	}
	return f, nil
}

// View renders the finder view
func (f *FinderView) View() string {
	var content string
	content += "Search: " + f.query + "\n\n"

	for i, result := range f.results {
		if i == f.cursor {
			content += "> " + result + "\n"
		} else {
			content += "  " + result + "\n"
		}
	}

	return f.style.Render(content)
}
