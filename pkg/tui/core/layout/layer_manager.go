package layout

import (
	"sort"

	"github.com/charmbracelet/bubbles/v2/key"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// LayerManager manages a stack of layer overlays
type LayerManager struct {
	layers []Layer
	width  int
	height int
}

// NewLayerManager creates a new layer manager
func NewLayerManager() *LayerManager {
	return &LayerManager{
		layers: make([]Layer, 0),
	}
}

// SetSize updates the size of the layer manager and all layers
func (lm *LayerManager) SetSize(w, h int) {
	lm.width, lm.height = w, h
	for i := range lm.layers {
		lm.layers[i].SetSize(w, h)
	}
}

// Push adds a new layer to the stack
func (lm *LayerManager) Push(l Layer) tea.Cmd {
	// Defocus current top layer
	if top := lm.Top(); top != nil {
		top.SetFocused(false)
	}

	// Focus and size the new layer
	l.SetFocused(true)
	l.SetSize(lm.width, lm.height)

	// Add to stack and sort by Z-order
	lm.layers = append(lm.layers, l)
	lm.sortLayers()

	return nil
}

// Pop removes the top layer from the stack
func (lm *LayerManager) Pop() tea.Cmd {
	if len(lm.layers) == 0 {
		return nil
	}

	lm.layers = lm.layers[:len(lm.layers)-1]

	// Focus the new top layer if it exists
	if top := lm.Top(); top != nil {
		top.SetFocused(true)
	}

	return nil
}

// PopByID removes a specific layer by ID
func (lm *LayerManager) PopByID(id string) tea.Cmd {
	for i, layer := range lm.layers {
		if layer.LayerMeta().ID == id {
			// If we're removing the top layer, handle focus transfer
			wasTop := i == len(lm.layers)-1

			// Remove the layer
			lm.layers = append(lm.layers[:i], lm.layers[i+1:]...)

			// If we removed the top layer, focus the new top
			if wasTop {
				if top := lm.Top(); top != nil {
					top.SetFocused(true)
				}
			}
			break
		}
	}
	return nil
}

// Top returns the topmost layer (highest Z-order)
func (lm *LayerManager) Top() Layer {
	if len(lm.layers) == 0 {
		return nil
	}
	return lm.layers[len(lm.layers)-1]
}

// HasLayers returns true if there are any active layers
func (lm *LayerManager) HasLayers() bool {
	return len(lm.layers) > 0
}

// sortLayers sorts layers by Z-order (lowest to highest)
func (lm *LayerManager) sortLayers() {
	sort.Slice(lm.layers, func(i, j int) bool {
		return lm.layers[i].LayerMeta().Z < lm.layers[j].LayerMeta().Z
	})
}

// Update routes input to the appropriate layer
func (lm *LayerManager) Update(msg tea.Msg) (*LayerManager, tea.Cmd, bool) {
	if top := lm.Top(); top != nil {
		meta := top.LayerMeta()

		// Check for dismiss keys
		if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
			for _, dismissKey := range meta.DismissKeys {
				if keyMsg.String() == dismissKey {
					lm.Pop()
					return lm, nil, true
				}
			}
		}

		// Update the top layer
		newTop, cmd := top.Update(msg)
		if nt, ok := newTop.(Layer); ok {
			lm.layers[len(lm.layers)-1] = nt
		}

		// Return whether input was handled based on layer settings
		handled := meta.Modal || meta.CaptureKeys
		return lm, cmd, handled
	}

	return lm, nil, false
}

// HelpBindings returns the help bindings for the current top layer
func (lm *LayerManager) HelpBindings() []key.Binding {
	if top := lm.Top(); top != nil {
		return top.Bindings()
	}
	return nil
}

// RenderOver renders all layers over the base content using lipgloss Canvas
func (lm *LayerManager) RenderOver(base string) string {
	if len(lm.layers) == 0 {
		return base
	}

	canvas := lipgloss.NewCanvas(lipgloss.NewLayer(base))

	// Add all layers
	for _, layer := range lm.layers {
		meta := layer.LayerMeta()
		content := layer.View()

		if content == "" {
			continue
		}

		x, y := lm.calculatePosition(meta.Pos, content)
		canvas.AddLayers(lipgloss.NewLayer(content).X(x).Y(y).Z(meta.Z))
	}

	return canvas.Render()
}

// calculatePosition determines the X,Y coordinates for layer placement
func (lm *LayerManager) calculatePosition(pos Position, content string) (int, int) {
	contentWidth := lipgloss.Width(content)
	contentHeight := lipgloss.Height(content)

	var baseX, baseY int

	// Calculate base position based on anchor
	switch pos.Anchor {
	case TopLeft:
		baseX, baseY = 0, 0
	case TopCenter:
		baseX, baseY = (lm.width-contentWidth)/2, 0
	case TopRight:
		baseX, baseY = lm.width-contentWidth, 0
	case MiddleLeft:
		baseX, baseY = 0, (lm.height-contentHeight)/2
	case Center:
		baseX, baseY = (lm.width-contentWidth)/2, (lm.height-contentHeight)/2
	case MiddleRight:
		baseX, baseY = lm.width-contentWidth, (lm.height-contentHeight)/2
	case BottomLeft:
		baseX, baseY = 0, lm.height-contentHeight
	case BottomCenter:
		baseX, baseY = (lm.width-contentWidth)/2, lm.height-contentHeight
	case BottomRight:
		baseX, baseY = lm.width-contentWidth, lm.height-contentHeight
	default:
		baseX, baseY = 0, 0
	}

	// Apply offsets
	finalX := baseX + pos.X
	finalY := baseY + pos.Y

	// Ensure the layer stays within bounds
	if finalX < 0 {
		finalX = 0
	}
	if finalY < 0 {
		finalY = 0
	}
	if finalX+contentWidth > lm.width {
		finalX = lm.width - contentWidth
	}
	if finalY+contentHeight > lm.height {
		finalY = lm.height - contentHeight
	}

	return finalX, finalY
}
