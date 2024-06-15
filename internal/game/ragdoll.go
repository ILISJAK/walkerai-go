package game

import "github.com/hajimehoshi/ebiten/v2"

// Ragdoll represents a simple ragdoll
type Ragdoll struct {
	X, Y float64
}

// NewRagdoll creates a new ragdoll
func NewRagdoll(x, y float64) *Ragdoll {
	return &Ragdoll{X: x, Y: y}
}

// Update updates the ragdoll's state
func (r *Ragdoll) Update() {
	// Update ragdoll physics
}

// Draw draws the ragdoll on the screen
func (r *Ragdoll) Draw(screen *ebiten.Image) {
	// Draw ragdoll
}
