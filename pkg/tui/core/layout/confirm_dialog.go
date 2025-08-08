package layout

import (
	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// ConfirmedMsg is sent when the user confirms the dialog
type ConfirmedMsg struct{}

// CancelledMsg is sent when the user cancels the dialog
type CancelledMsg struct{}

// ConfirmDialog is an example layer implementation
type ConfirmDialog struct {
	focused bool
	width   int
	height  int
	text    string
	keys    struct {
		Confirm key.Binding
		Cancel  key.Binding
	}
}

// NewConfirmDialog creates a new confirmation dialog
func NewConfirmDialog(text string) *ConfirmDialog {
	cd := &ConfirmDialog{
		text: text,
	}

	cd.keys.Confirm = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	)
	cd.keys.Cancel = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	)

	return cd
}

// Init implements tea.Model
func (c *ConfirmDialog) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (c *ConfirmDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, c.keys.Confirm):
			return c, func() tea.Msg { return ConfirmedMsg{} }
		case key.Matches(msg, c.keys.Cancel):
			return c, func() tea.Msg { return CancelledMsg{} }
		}
	}
	return c, nil
}

// View implements tea.Model
func (c *ConfirmDialog) View() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(40).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("15"))

	content := c.text + "\n\n" +
		lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("[Enter]") +
		" Confirm • " +
		lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("[Esc]") +
		" Cancel"

	return style.Render(content)
}

// SetSize implements Sizeable
func (c *ConfirmDialog) SetSize(width, height int) tea.Cmd {
	c.width = width
	c.height = height
	return nil
}

// GetSize implements Sizeable
func (c *ConfirmDialog) GetSize() (int, int) {
	return c.width, c.height
}

// SetFocused implements FocusScope
func (c *ConfirmDialog) SetFocused(focused bool) (FocusScope, tea.Cmd) {
	newDialog := *c
	newDialog.focused = focused
	return &newDialog, nil
}

// IsFocused implements FocusScope
func (c *ConfirmDialog) IsFocused() bool {
	return c.focused
}

// Clone implements FocusScope
func (c *ConfirmDialog) Clone() FocusScope {
	clone := *c
	return &clone
}

// Bindings implements Help
func (c *ConfirmDialog) Bindings() []key.Binding {
	return []key.Binding{c.keys.Confirm, c.keys.Cancel}
}

// LayerMeta implements Layer
func (c *ConfirmDialog) LayerMeta() LayerMeta {
	return LayerMeta{
		ID:          "confirm",
		Z:           100,
		Modal:       true,
		CaptureKeys: true,
		DismissKeys: []string{"esc"},
		Scrim:       true,
		Pos: Position{
			Anchor: Center,
		},
	}
}
