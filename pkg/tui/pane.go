package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Pane struct {
	focused  bool
	title    string
	content  string
	width    int
	height   int
	viewport viewport.Model
}

// Messages for pane communication
type (
	SetContentMsg struct{ Content string }
	FocusMsg      struct{ Focused bool }
	SizeMsg       struct{ Width, Height int }
)

func NewPane(title, content string) Pane {
	pane := Pane{
		focused:  false,
		title:    title,
		content:  content,
		viewport: viewport.New(0, 0),
	}

	pane.viewport.SetContent(pane.content)
	return pane
}

func (p *Pane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case SetContentMsg:
		p.content = msg.Content
		p.viewport.SetContent(p.content)

	case FocusMsg:
		p.focused = msg.Focused

	case SizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.handleResize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if p.focused {
			cmd = p.handleKeys(msg)
		}

	case tea.MouseMsg:
		p.viewport, cmd = p.viewport.Update(msg)
	}
	return *p, cmd
}

func (p *Pane) View() string {
	borderColor := BorderNormal()
	if p.focused {
		borderColor = BorderFocused()
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(borderColor)

	h, v := style.GetFrameSize()
	borderDef := lipgloss.RoundedBorder()
	innerWidth := max(0, p.width-h)
	innerHeight := max(0, p.height-v)

	contentStyle := lipgloss.NewStyle().
		Width(innerWidth).
		Height(innerHeight).
		Foreground(ContentFgActive())

	if !p.focused {
		contentStyle = contentStyle.Foreground(ContentFgInactive())
	}

	titleStr := p.title
	if titleStr == "" {
		return style.Width(p.width).Height(p.height).Render(contentStyle.Render(p.viewport.View()))
	}

	renderedContent := contentStyle.Render(p.viewport.View())

	fillStyle := lipgloss.NewStyle().Foreground(borderColor)
	titleStyle := fillStyle.Bold(p.focused).Italic(p.focused).Padding(0, 1)

	titleStyled := titleStyle.Render(titleStr)
	tlen := lipgloss.Width(titleStyled)

	minSide := 1
	maxTitleWidth := innerWidth - 2*minSide
	if tlen > maxTitleWidth {
		truncLen := max(maxTitleWidth-3, 0)
		runes := []rune(titleStr)
		truncPlain := ""
		if len(runes) > truncLen {
			truncPlain = string(runes[:truncLen]) + "..."
		} else {
			truncPlain = titleStr
		}
		titleStyled = titleStyle.Render(truncPlain)
		tlen = lipgloss.Width(titleStyled)
	}

	fill := innerWidth - tlen
	effectiveFill := max(fill-2*minSide, 0)

	leftNum := minSide + effectiveFill/2
	rightNum := minSide + (effectiveFill - effectiveFill/2)

	leftFill := fillStyle.Render(strings.Repeat(borderDef.Top, leftNum))
	rightFill := fillStyle.Render(strings.Repeat(borderDef.Top, rightNum))

	innerTop := leftFill + titleStyled + rightFill

	topLeft := fillStyle.Render(borderDef.TopLeft)
	topRight := fillStyle.Render(borderDef.TopRight)
	topLine := topLeft + innerTop + topRight

	sideStyle := lipgloss.NewStyle().
		Border(borderDef, false, true, true, true).
		BorderForeground(borderColor)

	sidesContent := sideStyle.Render(renderedContent)

	fullView := lipgloss.JoinVertical(lipgloss.Left, topLine, sidesContent)

	// Ensure the returned string *exactly* occupies the negotiated rectangle.
	return lipgloss.NewStyle().
		Width(p.width).
		Height(p.height).
		Render(fullView)
}

func (p *Pane) handleResize(width, height int) {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder(), true)
	h, v := style.GetFrameSize()
	innerWidth := max(0, width-h)
	innerHeight := max(0, height-v)

	p.viewport.Width = innerWidth
	p.viewport.Height = innerHeight
}

func (p *Pane) handleScrolling(direction int) {
	if direction > 0 {
		p.viewport.ScrollUp(1)
	} else {
		p.viewport.ScrollDown(1)
	}
}

func (p *Pane) handleKeys(msg tea.KeyMsg) tea.Cmd {
	if msg.Type == tea.KeyUp || msg.String() == "k" {
		p.handleScrolling(1)
	}
	if msg.Type == tea.KeyDown || msg.String() == "j" {
		p.handleScrolling(-1)
	}
	return nil
}
