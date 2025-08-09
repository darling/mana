package core

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/glamour/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/darling/mana/pkg/llm"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type MainCmp struct {
	focused    bool
	width      int
	height     int
	vp         viewport.Model
	messages   []llm.Message
	llmManager *llm.Manager
	keys       mainKeyMap
	renderer   *glamour.TermRenderer
}

// ChatResponseMsg is delivered when the LLM returns a response
type ChatResponseMsg struct {
	Message llm.Message
	Err     error
}

func NewMainCmp(manager *llm.Manager) MainCmp {
	// Initialize with a sane default renderer; will be resized on first ComponentSizeMsg
	var r *glamour.TermRenderer
	if tmp, err := glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(80),
	); err == nil {
		r = tmp
	}
	return MainCmp{
		keys:       DefaultMainKeyMap,
		llmManager: manager,
		renderer:   r,
	}
}

func (m MainCmp) Init() tea.Cmd { return nil }

func (m MainCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newM := m // Copy

	switch msg := msg.(type) {
	case layout.ComponentSizeMsg:
		sidebarWidth := msg.Width / 4
		newM.width = msg.Width - sidebarWidth
		newM.height = msg.Height
		innerW, innerH := newM.innerDimensions()
		newM.vp = viewport.New(
			viewport.WithWidth(innerW),
			viewport.WithHeight(innerH),
		)
		// (Re)create markdown renderer to match inner width
		if r, err := glamour.NewTermRenderer(
			glamour.WithEnvironmentConfig(),
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(innerW),
		); err == nil {
			newM.renderer = r
		} else {
			newM.renderer = nil
		}
		newM.vp.SetContent(newM.renderMessages(innerW))
	case layout.ConfirmedMsg:
		// no-op in chat view
	case layout.CancelledMsg:
		// no-op in chat view
	case layout.PromptSubmittedMsg:
		text := strings.TrimSpace(msg.Text)
		if text == "" {
			return newM, nil
		}
		// Append user message
		newM.messages = append(newM.messages, llm.Message{Role: "user", Content: text})
		innerW, _ := newM.innerDimensions()
		newM.vp.SetContent(newM.renderMessages(innerW))
		newM.vp.GotoBottom()

		// If we have an LLM, fire off generation
		if newM.llmManager != nil {
			history := append([]llm.Message(nil), newM.messages...)
			cmd := func() tea.Msg {
				resp, err := newM.llmManager.Generate(context.Background(), history)
				return ChatResponseMsg{Message: resp, Err: err}
			}
			return newM, cmd
		}
	case ChatResponseMsg:
		if msg.Err == nil && msg.Message.Content != "" {
			newM.messages = append(newM.messages, llm.Message{Role: "assistant", Content: msg.Message.Content, Provider: msg.Message.Provider, ID: msg.Message.ID})
			innerW, _ := newM.innerDimensions()
			newM.vp.SetContent(newM.renderMessages(innerW))
			newM.vp.GotoBottom()
		}
	case tea.KeyPressMsg:
		if !m.focused {
			return newM, nil
		}

		switch {
		case key.Matches(msg, m.keys.Redraw):
			// force refresh
			innerW, _ := newM.innerDimensions()
			newM.vp.SetContent(newM.renderMessages(innerW))
		case key.Matches(msg, m.keys.Create):
			return newM, func() tea.Msg { return layout.ShowPromptDialogMsg{} }
		case key.Matches(msg, m.keys.ShowDialog):
			return newM, func() tea.Msg { return layout.ShowPromptDialogMsg{} }
		}

		// Pass other keypresses to viewport for scrolling
		var vpCmd tea.Cmd
		newM.vp, vpCmd = newM.vp.Update(msg)
		return newM, vpCmd

	case tea.MouseMsg:
		// Forward mouse events (e.g., wheel) to viewport for scrolling
		var vpCmd tea.Cmd
		newM.vp, vpCmd = newM.vp.Update(msg)
		return newM, vpCmd
	}

	return newM, nil
}

func (m MainCmp) View() string {
	content := m.vp.View()

	var boxStyle lipgloss.Style
	if m.focused {
		boxStyle = FocusedBox
	} else {
		boxStyle = BlurredBox
	}

	// Render within a fixed-size box, clip and nowrap to avoid layout push
	innerW, innerH := m.innerDimensions()
	clipped := lipgloss.NewStyle().
		Width(innerW).Height(innerH).
		MaxWidth(innerW).MaxHeight(innerH).
		// ensure no extra newlines/padding sneak in
		Align(lipgloss.Left).
		Render(content)
	return boxStyle.Width(m.width).Height(m.height).Render(clipped)
}

func (m MainCmp) SetFocused(focused bool) (layout.Focusable, tea.Cmd) {
	newM := m
	newM.focused = focused
	return newM, nil
}

func (m MainCmp) IsFocused() bool { return m.focused }

func (m MainCmp) Clone() layout.Focusable {
	return MainCmp{
		focused:    m.focused,
		width:      m.width,
		height:     m.height,
		vp:         m.vp,
		messages:   append([]llm.Message(nil), m.messages...),
		llmManager: m.llmManager,
		keys:       m.keys,
		renderer:   m.renderer,
	}
}

func (m MainCmp) Bindings() []key.Binding {
	return []key.Binding{m.keys.Redraw, m.keys.Create, m.keys.ShowDialog}
}

func (m MainCmp) renderMessages(innerWidth int) string {
	if len(m.messages) == 0 {
		return ""
	}
	var b strings.Builder
	for i, msg := range m.messages {
		if i > 0 {
			b.WriteString("\n\n")
		}
		// role header
		role := msg.Role
		if role == "" {
			role = "assistant"
		}
		b.WriteString(fmt.Sprintf("%s:\n", role))
		if m.renderer != nil {
			if out, err := m.renderer.Render(msg.Content); err == nil {
				b.WriteString(out)
			} else {
				b.WriteString(hardWrap(msg.Content, innerWidth))
			}
		} else {
			b.WriteString(hardWrap(msg.Content, innerWidth))
		}
	}
	return b.String()
}

func (m MainCmp) innerDimensions() (int, int) {
	// Compute inner dimensions based on the outer box style chrome.
	// Focused and blurred styles currently share the same padding/frame sizes.
	s := FocusedBox
	innerW := m.width - s.GetHorizontalPadding() - s.GetHorizontalFrameSize()
	innerH := m.height - s.GetVerticalPadding() - s.GetVerticalFrameSize()
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}
	return innerW, innerH
}

// hardWrap wraps a string to the given width by rune, avoiding control runes.
func hardWrap(s string, width int) string {
	if width <= 0 || s == "" {
		return s
	}
	var lines []string
	var line []rune
	col := 0
	for _, r := range []rune(s) {
		if r == '\n' {
			lines = append(lines, string(line))
			line = line[:0]
			col = 0
			continue
		}
		// skip control characters except tab (treated as single space)
		if unicode.IsControl(r) && r != '\t' {
			continue
		}
		line = append(line, r)
		col++
		if col >= width {
			lines = append(lines, string(line))
			line = line[:0]
			col = 0
		}
	}
	if len(line) > 0 {
		lines = append(lines, string(line))
	}
	return strings.Join(lines, "\n")
}
