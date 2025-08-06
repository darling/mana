package layout

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/darling/mana/pkg/tui/core/components"
)

// FocusChangedMsg is a message used to signal that the focused component has changed.
// It is used to trigger updates in other components, such as the help view in the status bar.
type FocusChangedMsg struct{}

type Focusable interface {
	components.Component

	SetFocused(bool) (Focusable, tea.Cmd)
	IsFocused() bool
	Clone() Focusable
}

type FocusManager struct {
	components   []Focusable
	focusedIndex int
	wrap         bool
}

func NewFocusManager(components []Focusable, wrap bool) FocusManager {
	return FocusManager{
		components:   cloneComponents(components),
		focusedIndex: -1,
		wrap:         wrap,
	}
}

func cloneComponents(cs []Focusable) []Focusable {
	cloned := make([]Focusable, len(cs))
	for i, c := range cs {
		cloned[i] = c.Clone()
	}
	return cloned
}

func (fm FocusManager) blurCurrent() (FocusManager, tea.Cmd) {
	if fm.focusedIndex < 0 || fm.focusedIndex >= len(fm.components) {
		return fm, nil
	}

	newFM := fm
	var cmd tea.Cmd
	newFM.components[newFM.focusedIndex], cmd = newFM.components[newFM.focusedIndex].SetFocused(false)
	return newFM, cmd
}

// focus is an internal helper that handles the logic of changing focus. It blurs
// the previously focused component, focuses the new one, and dispatches a
// FocusChangedMsg to signal the change.
func (fm FocusManager) focus(index int) (FocusManager, tea.Cmd) {
	if fm.focusedIndex == index && fm.focusedIndex != -1 {
		return fm, nil // No change, no command.
	}

	// Get a blurred manager and the blur command.
	blurredFM, blurCmd := fm.blurCurrent()

	// Use this new manager state to set focus.
	focusedFM := blurredFM
	focusedFM.focusedIndex = index
	var focusCmd tea.Cmd
	focusedFM.components[focusedFM.focusedIndex], focusCmd = focusedFM.components[focusedFM.focusedIndex].SetFocused(true)

	// Command to signal that focus has changed.
	focusChangedCmd := func() tea.Msg { return FocusChangedMsg{} }

	return focusedFM, tea.Batch(blurCmd, focusCmd, focusChangedCmd)
}

func (fm FocusManager) FocusNext() (FocusManager, tea.Cmd) {
	if len(fm.components) == 0 {
		return fm, nil
	}

	if fm.focusedIndex == -1 {
		return fm.focus(0)
	}

	if fm.focusedIndex == len(fm.components)-1 && !fm.wrap {
		return fm, nil
	}

	nextIndex := (fm.focusedIndex + 1) % len(fm.components)
	return fm.focus(nextIndex)
}

func (fm FocusManager) FocusPrev() (FocusManager, tea.Cmd) {
	if len(fm.components) == 0 {
		return fm, nil
	}

	if fm.focusedIndex == -1 {
		return fm.focus(len(fm.components) - 1)
	}

	if fm.focusedIndex == 0 && !fm.wrap {
		return fm, nil
	}

	prevIndex := fm.focusedIndex - 1
	if prevIndex < 0 {
		prevIndex = len(fm.components) - 1
	}
	return fm.focus(prevIndex)
}

func (fm FocusManager) Focus(index int) (FocusManager, tea.Cmd, error) {
	if index < 0 || index >= len(fm.components) {
		return fm, nil, errors.New("index out of bounds")
	}

	newFM, cmd := fm.focus(index)
	return newFM, cmd, nil
}

func (fm FocusManager) UpdateFocused(msg tea.Msg) (FocusManager, tea.Cmd) {
	if fm.focusedIndex < 0 || fm.focusedIndex >= len(fm.components) {
		return fm, nil
	}

	newFM := fm
	updatedModel, cmd := newFM.components[newFM.focusedIndex].Update(msg)

	if updatedFocusable, ok := updatedModel.(Focusable); ok {
		newFM.components[newFM.focusedIndex] = updatedFocusable.Clone()
	}

	return newFM, cmd
}

func (fm FocusManager) UpdateAll(msg tea.Msg) (FocusManager, tea.Cmd) {
	var cmds []tea.Cmd
	newFM := fm

	for i, component := range newFM.components {
		updatedModel, cmd := component.Update(msg)
		if updatedFocusable, ok := updatedModel.(Focusable); ok {
			newFM.components[i] = updatedFocusable.Clone()
		}
		cmds = append(cmds, cmd)
	}

	return newFM, tea.Batch(cmds...)
}

func (fm FocusManager) Get(index int) (Focusable, error) {
	if index < 0 || index >= len(fm.components) {
		return nil, errors.New("index out of bounds")
	}
	return fm.components[index].Clone(), nil
}

func (fm FocusManager) GetAll() []Focusable {
	return cloneComponents(fm.components)
}

func (fm FocusManager) GetFocused() (Focusable, error) {
	if fm.focusedIndex < 0 || fm.focusedIndex >= len(fm.components) {
		return nil, errors.New("no component is focused")
	}
	return fm.components[fm.focusedIndex].Clone(), nil
}

func (fm FocusManager) Set(index int, component Focusable) (FocusManager, error) {
	if index < 0 || index >= len(fm.components) {
		return fm, errors.New("index out of bounds")
	}
	newFM := fm.Clone()
	newFM.components[index] = component.Clone()
	return newFM, nil
}

func (fm FocusManager) Clone() FocusManager {
	return FocusManager{
		components:   cloneComponents(fm.components),
		focusedIndex: fm.focusedIndex,
		wrap:         fm.wrap,
	}
}
