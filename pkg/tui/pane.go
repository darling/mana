package tui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PaneModel struct {
	id      string // <-- Added: Unique ID for the pane
	focused bool

	title         string
	titlePosition lipgloss.Position
	jumpNum       int

	width  int
	height int

	content string

	style lipgloss.Style

	viewport viewport.Model
}

// Pane messages
type (
	PaneContentMsg string
	// Renamed from PaneFocusMsg to be more specific and avoid conflicts.
	// This message carries the ID of the pane to be focused.
	FocusPaneMsg string
	PaneSizeMsg  struct{ Width, Height int }
)

// NewPane now accepts an ID.
func NewPane(id, title, innerContent string) *PaneModel {
	pane := &PaneModel{
		id:            id, // <-- Added: Store the ID
		focused:       false,
		title:         title,
		titlePosition: lipgloss.Left,
		jumpNum:       0,
		content:       innerContent,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true),
		viewport: viewport.New(0, 0),
	}

	pane.viewport.SetContent(pane.content)

	return pane
}

func (p *PaneModel) SetTitle(title string, position lipgloss.Position) {
	p.title = title
	p.titlePosition = position
}

func (p *PaneModel) SetJumpNum(num int) {
	p.jumpNum = num
}

func (p *PaneModel) Init() tea.Cmd { return nil }

func (p *PaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	// THIS CASE IS REMOVED. RootModel will send a specific PaneSizeMsg instead.
	// case tea.WindowSizeMsg:
	//     p.handleResize(msg.Width, msg.Height)

	case PaneContentMsg:
		p.content = string(msg)
		p.viewport.SetContent(p.content)
		p.viewport.GotoBottom() // Go to bottom on new content

	// This pane will set its focus state based on whether its ID matches the message.
	case FocusPaneMsg:
		p.focused = (p.id == string(msg))

	case PaneSizeMsg:
		p.handleResize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if p.focused {
			cmd = p.handleKeys(msg)
		}

	case tea.MouseMsg:
		// Pass to viewport first for its internal handling (like scrolling).
		p.viewport, cmd = p.viewport.Update(msg)
		cmds = append(cmds, cmd)

		// Then, handle our custom mouse logic (like focus on click).
		cmd = p.handleMouse(msg)
		cmds = append(cmds, cmd)
	}

	return p, tea.Batch(cmds...)
}

func (p *PaneModel) View() string {
	borderColor := BorderNormal()
	if p.focused {
		borderColor = BorderFocused()
	}
	styled := p.style.BorderForeground(borderColor)

	titleStr := p.title
	if p.jumpNum > 0 {
		titleStr = fmt.Sprintf("[%d] %s", p.jumpNum, titleStr)
	}

	h, v := styled.GetFrameSize()
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

	if titleStr == "" {
		return styled.Width(p.width).Height(p.height).Render(contentStyle.Render(p.viewport.View()))
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

	startInEffective := int(math.Round(float64(effectiveFill) * float64(p.titlePosition)))
	leftNum := minSide + startInEffective
	rightNum := minSide + (effectiveFill - startInEffective)

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

	return fullView
}

// DELETE these deprecated methods. They encourage an anti-pattern.
/*
func (p *PaneModel) SetContent(innerContent string) {
	p.content = innerContent
	p.viewport.SetContent(p.content)
}

func (p *PaneModel) SetFocus(focused bool) {
	p.focused = focused
}

func (p *PaneModel) SetSize(width, height int) {
	p.width = width
	p.height = height
}
*/

func (p *PaneModel) handleResize(width, height int) {
	p.width = width
	p.height = height

	h, v := p.style.GetFrameSize()
	innerWidth := max(0, width-h)
	innerHeight := max(0, height-v)

	p.viewport.Width = innerWidth
	p.viewport.Height = innerHeight
}

func (p *PaneModel) handleScrolling(direction int) {
	if direction > 0 {
		p.viewport.LineUp(1)
	} else {
		p.viewport.LineDown(1)
	}
}

func (p *PaneModel) handleKeys(msg tea.KeyMsg) tea.Cmd {
	if msg.Type == tea.KeyUp || msg.String() == "k" {
		p.handleScrolling(1)
	}
	if msg.Type == tea.KeyDown || msg.String() == "j" {
		p.handleScrolling(-1)
	}
	return nil
}

func (p *PaneModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if msg.Type == tea.MouseLeft {
		return func() tea.Msg { return FocusPaneMsg(p.id) }
	}
	if msg.Type == tea.MouseWheelUp {
		p.handleScrolling(1)
	}
	if msg.Type == tea.MouseWheelDown {
		p.handleScrolling(-1)
	}
	return nil
}
