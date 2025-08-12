package core

import (
	"context"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/viewport"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/glamour/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/darling/mana/pkg/chat"
	"github.com/darling/mana/pkg/tui/core/layout"
)

type MainCmp struct {
	focused     bool
	width       int
	height      int
	vp          viewport.Model
	messages    []chat.Message
	chatService chat.Service
	convID      chat.ConversationID
	streamCh    <-chan chat.StreamEvent
	cancel      context.CancelFunc
	keys        mainKeyMap
	renderer    *glamour.TermRenderer
}

// ChatStreamEventMsg carries streaming events from the chat service
type ChatStreamEventMsg struct{ Ev chat.StreamEvent }
type ChatStreamClosedMsg struct{}

func NewMainCmp(service chat.Service) MainCmp {
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
		keys:        DefaultMainKeyMap,
		chatService: service,
		renderer:    r,
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
		if newM.chatService != nil {
			ctx := context.Background()
			if newM.convID == "" {
				if id, err := newM.chatService.NewConversation(ctx, nil); err == nil {
					newM.convID = id
				} else {
					return newM, nil
				}
			}
			// Update local view first
			newM.messages = append(newM.messages, chat.Message{Role: chat.RoleUser, Content: text})
			// Append an assistant placeholder locally so we can stream into it
			newM.messages = append(newM.messages, chat.Message{Role: chat.RoleAssistant, Content: ""})
			innerW, _ := newM.innerDimensions()
			newM.vp.SetContent(newM.renderMessages(innerW))
			newM.vp.GotoBottom()

			// Persist user message and start streaming in the background
			_, _ = newM.chatService.AddUserMessage(ctx, newM.convID, text, nil)
			events, cancel, err := newM.chatService.StartAssistantStream(ctx, newM.convID, chat.GenerateOptions{})
			if err == nil {
				newM.streamCh = events
				newM.cancel = cancel
				return newM, newM.readNextStreamCmd()
			}
		}
	case ChatStreamEventMsg:
		// Append delta to the last assistant message locally
		if len(newM.messages) > 0 {
			last := len(newM.messages) - 1
			if newM.messages[last].Role == chat.RoleAssistant && msg.Ev.Delta != "" {
				newM.messages[last].Content += msg.Ev.Delta
				innerW, _ := newM.innerDimensions()
				newM.vp.SetContent(newM.renderMessages(innerW))
				newM.vp.GotoBottom()
			}
		}
		if msg.Ev.Err != nil || msg.Ev.Done {
			if newM.cancel != nil {
				newM.cancel()
			}
			newM.streamCh = nil
			newM.cancel = nil
			return newM, nil
		}
		return newM, newM.readNextStreamCmd()
	case ChatStreamClosedMsg:
		if newM.cancel != nil {
			newM.cancel()
		}
		newM.streamCh = nil
		newM.cancel = nil
		return newM, nil
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
		focused:     m.focused,
		width:       m.width,
		height:      m.height,
		vp:          m.vp,
		messages:    append([]chat.Message(nil), m.messages...),
		chatService: m.chatService,
		convID:      m.convID,
		streamCh:    m.streamCh,
		cancel:      m.cancel,
		keys:        m.keys,
		renderer:    m.renderer,
	}
}

func (m MainCmp) Bindings() []key.Binding {
	return []key.Binding{m.keys.Redraw, m.keys.Create, m.keys.ShowDialog}
}

// readNextStreamCmd reads a single event from the current stream channel.
// It returns ChatStreamClosedMsg when the stream channel is nil or closed.
func (m MainCmp) readNextStreamCmd() tea.Cmd {
	ch := m.streamCh
	return func() tea.Msg {
		if ch == nil {
			return ChatStreamClosedMsg{}
		}
		ev, ok := <-ch
		if !ok {
			return ChatStreamClosedMsg{}
		}
		return ChatStreamEventMsg{Ev: ev}
	}
}

func (m MainCmp) renderMessages(innerWidth int) string {
	if len(m.messages) == 0 {
		return ""
	}
	var b strings.Builder
	for _, msg := range m.messages {
		role := string(msg.Role)
		if role == "" {
			role = "assistant"
		}

		// Style the role label
		var label string
		switch role {
		case string(chat.RoleUser):
			label = MessageLabelUser.Render("User:")
		case string(chat.RoleSystem):
			label = MessageLabelSystem.Render("System:")
		default:
			label = MessageLabelAssistant.Render("Assistant:")
		}

		var body string
		if m.renderer != nil {
			if out, err := m.renderer.Render(msg.Content); err == nil {
				body = out
			} else {
				body = hardWrap(msg.Content, innerWidth)
			}
		} else {
			body = hardWrap(msg.Content, innerWidth)
		}

		// Compose a message block with label and content
		block := lipgloss.JoinVertical(
			lipgloss.Left,
			label,
			body,
		)

		b.WriteString("\n" + MessageBlock.Render(block))
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
