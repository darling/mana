package layout

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea/v2"

	"github.com/darling/mana/pkg/tui/core/components"
)

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

func (fm FocusManager) FocusNext() (FocusManager, tea.Cmd) {
	if len(fm.components) == 0 {
		return fm, nil
	}

	if fm.focusedIndex == -1 {
		newFM := fm
		newFM.focusedIndex = 0
		var focusCmd tea.Cmd
		newFM.components[newFM.focusedIndex], focusCmd = newFM.components[newFM.focusedIndex].SetFocused(true)
		return newFM, focusCmd
	}

	if fm.focusedIndex == len(fm.components)-1 && !fm.wrap {
		return fm, nil
	}

	newFM, blurCmd := fm.blurCurrent()
	newFM.focusedIndex++

	if newFM.focusedIndex >= len(newFM.components) {
		newFM.focusedIndex = 0
	}

	var focusCmd tea.Cmd
	newFM.components[newFM.focusedIndex], focusCmd = newFM.components[newFM.focusedIndex].SetFocused(true)

	return newFM, tea.Batch(blurCmd, focusCmd)
}

func (fm FocusManager) FocusPrev() (FocusManager, tea.Cmd) {
	if len(fm.components) == 0 {
		return fm, nil
	}

	if fm.focusedIndex == -1 {
		newFM := fm
		newFM.focusedIndex = len(newFM.components) - 1
		var focusCmd tea.Cmd
		newFM.components[newFM.focusedIndex], focusCmd = newFM.components[newFM.focusedIndex].SetFocused(true)
		return newFM, focusCmd
	}

	if fm.focusedIndex == 0 && !fm.wrap {
		return fm, nil
	}

	newFM, blurCmd := fm.blurCurrent()
	newFM.focusedIndex--

	if newFM.focusedIndex < 0 {
		newFM.focusedIndex = len(newFM.components) - 1
	}

	var focusCmd tea.Cmd
	newFM.components[newFM.focusedIndex], focusCmd = newFM.components[newFM.focusedIndex].SetFocused(true)

	return newFM, tea.Batch(blurCmd, focusCmd)
}

func (fm FocusManager) Focus(index int) (FocusManager, tea.Cmd, error) {
	if index < 0 || index >= len(fm.components) {
		return fm, nil, errors.New("index out of bounds")
	}

	newFM, blurCmd := fm.blurCurrent()
	newFM.focusedIndex = index

	var focusCmd tea.Cmd
	newFM.components[newFM.focusedIndex], focusCmd = newFM.components[newFM.focusedIndex].SetFocused(true)

	return newFM, tea.Batch(blurCmd, focusCmd), nil
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

func (fm FocusManager) Clone() FocusManager {
	return FocusManager{
		components:   cloneComponents(fm.components),
		focusedIndex: fm.focusedIndex,
		wrap:         fm.wrap,
	}
}
