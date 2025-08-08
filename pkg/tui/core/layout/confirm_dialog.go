package layout

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/textinput"
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

// PromptDialog is a modal layer to capture a prompt from the user
type PromptDialog struct {
	focused bool
	width   int
	height  int
	input   textinput.Model
	keys    struct {
		Submit key.Binding
		Cancel key.Binding
	}
}

func NewPromptDialog(initial string) *PromptDialog {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.SetValue(initial)
	ti.Prompt = "> "
	ti.CursorEnd()
	ti.Focus()

	pd := &PromptDialog{input: ti}
	pd.keys.Submit = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "send"))
	pd.keys.Cancel = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "cancel"))
	return pd
}

func (p *PromptDialog) Init() tea.Cmd { return nil }

func (p *PromptDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(m, p.keys.Submit):
			text := p.input.Value()
			return p, func() tea.Msg { return PromptSubmittedMsg{Text: text} }
		case key.Matches(m, p.keys.Cancel):
			return p, func() tea.Msg { return CancelledMsg{} }
		}
	}
	var cmd tea.Cmd
	p.input, cmd = p.input.Update(msg)
	return p, cmd
}

func (p *PromptDialog) View() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(max(40, p.width/2)).
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("15"))

	controls := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Render("[Enter]") +
		" Send • " +
		lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("[Esc]") +
		" Cancel"

	return style.Render(p.input.View() + "\n\n" + controls)
}

func (p *PromptDialog) SetSize(width, height int) tea.Cmd {
	p.width, p.height = width, height
	boxWidth := max(40, width/2)
	inputWidth := boxWidth - 6 // account for border and padding
	if inputWidth < 10 {
		inputWidth = 10
	}
	p.input.SetWidth(inputWidth)
	return nil
}
func (p *PromptDialog) GetSize() (int, int) { return p.width, p.height }
func (p *PromptDialog) SetFocused(focused bool) (FocusScope, tea.Cmd) {
	p.focused = focused
	if focused {
		p.input.Focus()
	} else {
		p.input.Blur()
	}
	return p, nil
}
func (p *PromptDialog) IsFocused() bool         { return p.focused }
func (p *PromptDialog) Clone() FocusScope       { clone := *p; return &clone }
func (p *PromptDialog) Bindings() []key.Binding { return []key.Binding{p.keys.Submit, p.keys.Cancel} }
func (p *PromptDialog) LayerMeta() LayerMeta {
	return LayerMeta{
		ID:          "prompt",
		Z:           100,
		Modal:       true,
		CaptureKeys: true,
		DismissKeys: []string{"esc"},
		Scrim:       true,
		Pos:         Position{Anchor: Center},
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
