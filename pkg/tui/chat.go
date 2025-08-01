package tui

import (
	"strings"

	"github.com/charmbracelet/glamour/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/darling/mana/pkg/llm"
)

// Conversation represents a chat session with message history
type Conversation struct {
	ID       string
	Messages []llm.Message // Full history, user + assistant
}

// renderHistory converts a slice of messages to a display string
func renderHistory(msgs []llm.Message) string {
	var b strings.Builder
	for _, m := range msgs {
		switch m.Role {
		case "user":
			b.WriteString(lipgloss.NewStyle().Bold(true).
				Render("you:\n") + m.Content + "\n\n")
		default:
			content, err := glamour.Render(m.Content, "dark")
			if err != nil {
				content = m.Content // Fallback to raw content if rendering fails
			}
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6cf")).
				Render("ai:\n") + content + "\n\n")
		}
	}
	return b.String()
}
