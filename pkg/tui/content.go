package tui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/darling/mana/pkg/llm"
	"github.com/google/uuid"
)

// LLM message types
type (
	LLMResponseMsg struct{ Response llm.Message }
	LLMErrorMsg    struct{ Err error }
)

// Content interface for the content pane
type Content interface {
	Init() tea.Cmd
	Update(tea.Msg) (Content, tea.Cmd)
	View() string
	SetSize(w, h int)
	SetFocused(bool)
}

// generateCmd creates a command to generate LLM response asynchronously
func generateCmd(mgr *llm.Manager, hist []llm.Message) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		res, err := mgr.Generate(ctx, hist)
		if err != nil {
			return LLMErrorMsg{err}
		}
		return LLMResponseMsg{res}
	}
}

type ContentModel struct {
	Pane
	llmManager *llm.Manager
	convo      *Conversation
	loading    bool
	content    string
}

func NewContentModel(llmManager *llm.Manager) Content {
	initialContent := "Welcome to Mana!"

	if llmManager == nil {
		initialContent += "\n\nNo LLM manager found. Please configure an LLM provider to start a conversation."
	} else {
		initialContent += "\n\nPress 'c' to start a conversation with the AI."
	}

	// TODO: Custom components for empty state for rich, fun, interactive experience

	return &ContentModel{
		Pane:       NewPane("Chat", initialContent),
		llmManager: llmManager,
		convo:      &Conversation{ID: "default", Messages: []llm.Message{}},
		loading:    false,
	}
}

// Init implements the Content interface
func (m ContentModel) Init() tea.Cmd {
	return nil
}

// SetSize implements the Content interface
func (m *ContentModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.handleResize(w, h)
}

// SetFocused sets the focus state of the content model
func (m *ContentModel) SetFocused(focused bool) {
	m.focused = focused
}

// SetLoading sets the loading state of the content model
func (m *ContentModel) SetLoading(loading bool) {
	m.loading = loading
}

func (m *ContentModel) Update(msg tea.Msg) (Content, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case PromptSubmittedMsg:
		if m.llmManager != nil {
			user := llm.Message{
				ID:       uuid.NewString(),
				Provider: "user",
				Role:     "user",
				Content:  msg.Prompt,
			}
			m.convo.Messages = append(m.convo.Messages, user)
			m.SetLoading(true)
			m.syncViewportContent()
			cmds = append(cmds, generateCmd(m.llmManager, m.convo.Messages))
		}

	case LLMResponseMsg:
		m.convo.Messages = append(m.convo.Messages, msg.Response)
		m.SetLoading(false)
		m.syncViewportContent()

	case LLMErrorMsg:
		m.SetLoading(false)
		m.syncViewportContent()
		// cmds = append(cmds, func() tea.Msg {
		// 	return StatusErrorMsg(msg.Err)
		// })

	case tea.KeyMsg:
		if m.focused {
			switch {
			case key.Matches(msg, keys.Prompt):
				cmd := func() tea.Msg {
					return OpenDialogMsg{Model: NewPromptDialog()}
				}
				cmds = append(cmds, cmd)
			}
		}
	}

	var paneCmd tea.Cmd
	m.Pane, paneCmd = (&m.Pane).Update(msg)
	cmds = append(cmds, paneCmd)

	return m, tea.Batch(cmds...)
}

func (m ContentModel) View() string {
	return m.Pane.View()
}

func (m *ContentModel) syncViewportContent() {
	content := ""
	if m.convo != nil && len(m.convo.Messages) > 0 {
		content = renderHistory(m.convo.Messages)
	}
	if m.loading {
		content += "\n• AI is thinking..."
	}
	if content != m.content {
		m.content = content
		m.viewport.SetContent(content)
		m.viewport.GotoBottom()
	}
}
