package tui

import (
	"github.com/charmbracelet/lipgloss"
)

type Pane struct {
	focused    bool
	width      int
	height     int
	content    string
	style      lipgloss.Style
	focusCol   string
	unfocusCol string
}

func NewPane(innerContent string) *Pane {
	return &Pane{
		focused:    false,
		content:    innerContent,
		focusCol:   "39",
		unfocusCol: "240",
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true),
	}
}

func (p *Pane) SetFocus(focused bool) {
	p.focused = focused
}

func (p *Pane) SetSize(width, height int) {
	p.width = width
	p.height = height
}

func (p *Pane) SetContent(innerContent string) {
	p.content = innerContent
}

func (p *Pane) Render() string {
	borderColor := p.unfocusCol
	if p.focused {
		borderColor = p.focusCol
	}
	styled := p.style.BorderForeground(lipgloss.Color(borderColor))

	h, v := styled.GetFrameSize()
	innerWidth := max(0, p.width-h)
	innerHeight := max(0, p.height-v)

	contentStyle := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "15"})

	if !p.focused {
		contentStyle = contentStyle.Foreground(lipgloss.AdaptiveColor{Light: "245", Dark: "245"})
	}

	return styled.Render(contentStyle.Render(p.content))
}
