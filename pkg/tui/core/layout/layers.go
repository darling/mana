package layout

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/darling/mana/pkg/tui/core/components"
)

// FocusScope represents a component that can receive focus and provide help bindings
type FocusScope interface {
	components.Component
	Sizeable
	Help

	SetFocused(bool) (FocusScope, tea.Cmd)
	IsFocused() bool
	Clone() FocusScope
}

// Anchor defines where a layer should be positioned
type Anchor int

const (
	TopLeft Anchor = iota
	TopCenter
	TopRight
	MiddleLeft
	Center
	MiddleRight
	BottomLeft
	BottomCenter
	BottomRight
)

// Position defines where and how a layer should be placed
type Position struct {
	X, Y   int    // absolute offsets from anchor
	Anchor Anchor // positioning anchor
	Width  int    // desired bounding box (0 = auto)
	Height int    // desired bounding box (0 = auto)
}

// LayerMeta contains metadata about how a layer should behave
type LayerMeta struct {
	ID           string   // unique identifier
	Z            int      // z-order (higher = on top)
	Pos          Position // placement information
	Modal        bool     // captures all input
	CaptureMouse bool     // captures mouse events
	CaptureKeys  bool     // captures keyboard events
	DismissKeys  []string // keys that dismiss the layer
	Scrim        bool     // draw dim background under modal
}

// Layer represents a composable overlay component
type Layer interface {
	FocusScope
	LayerMeta() LayerMeta
}

// OpenLayerMsg is sent to request opening a new layer
type OpenLayerMsg struct {
	Layer Layer
}

// DismissLayerMsg is sent to request dismissing the top layer
type DismissLayerMsg struct{}

// DismissLayerByIDMsg is sent to request dismissing a specific layer
type DismissLayerByIDMsg struct {
	ID string
}

// ShowConfirmDialogMsg is a convenience message for showing a confirm dialog
type ShowConfirmDialogMsg struct {
	Text string
}
