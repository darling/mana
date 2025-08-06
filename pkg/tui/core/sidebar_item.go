package core

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/darling/mana/pkg/tui/core/layout"
)

// SidebarPaneCmp represents one of the panes within the sidebar, like "Conversations".
type SidebarPaneCmp struct {
	focused bool
	title   string
	width   int
	height  int
	content string
}

func NewSidebarPaneCmp(title string) SidebarPaneCmp {
	return SidebarPaneCmp{
		title:   title,
		content: "...", // Placeholder content
	}
}

func (p SidebarPaneCmp) Init() tea.Cmd { return nil }

func (p SidebarPaneCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case layout.ComponentSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
	}
	return p, nil
}

// View renders the pane. The focused state determines the border color.
func (p SidebarPaneCmp) View() string {
	var boxStyle lipgloss.Style
	if p.focused {
		boxStyle = FocusedBox
	} else {
		boxStyle = BlurredBox
	}

	// Calculate size for internal content, accounting for border and padding.
	contentWidth := p.width - boxStyle.GetHorizontalPadding() - boxStyle.GetHorizontalFrameSize()
	contentHeight := p.height - boxStyle.GetVerticalPadding() - boxStyle.GetVerticalFrameSize()

	header := lipgloss.NewStyle().Bold(true).Render(p.title)
	body := lipgloss.NewStyle().Render(p.content)

	// Ensure content fits within the calculated dimensions.
	body = lipgloss.NewStyle().
		Width(contentWidth).
		Height(contentHeight - lipgloss.Height(header)).
		Render(body)

	view := lipgloss.JoinVertical(lipgloss.Top, header, body)

	return boxStyle.Width(p.width).Height(p.height).Render(view)
}

func (p SidebarPaneCmp) SetFocused(focused bool) (layout.Focusable, tea.Cmd) {
	p.focused = focused
	return p, nil
}

func (p SidebarPaneCmp) IsFocused() bool {
	return p.focused
}

func (p SidebarPaneCmp) Clone() layout.Focusable {
	return SidebarPaneCmp{
		focused: p.focused,
		title:   p.title,
		width:   p.width,
		height:  p.height,
		content: p.content,
	}
}
