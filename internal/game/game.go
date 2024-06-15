package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	ragdoll *Ragdoll
}

func NewGame() *Game {
	return &Game{
		ragdoll: NewRagdoll(),
	}
}

func (g *Game) Update() error {
	g.ragdoll.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	g.ragdoll.Draw(screen)
	ebitenutil.DebugPrint(screen, "WalkerAI Ragdoll Simulation")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
