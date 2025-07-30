package tui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Pane struct {
	focused bool

	title         string
	titlePosition lipgloss.Position
	jumpNum       int

	width  int
	height int

	content string

	style lipgloss.Style
}

func NewPane(title, innerContent string) *Pane {
	return &Pane{
		focused: false,

		title: title,

		titlePosition: lipgloss.Left,
		jumpNum:       0,

		content: innerContent,

		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true),
	}
}

func (p *Pane) SetTitle(title string, position lipgloss.Position) {
	p.title = title
	p.titlePosition = position
}

func (p *Pane) SetJumpNum(num int) {
	p.jumpNum = num
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
		return styled.Render(contentStyle.Render(p.content))
	}

	renderedContent := contentStyle.Render(p.content)

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
