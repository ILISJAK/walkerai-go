package game

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "log"
)

type Game struct{}

func (g *Game) Update() error {
    // Update game logic
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Draw the game world
    ebitenutil.DebugPrint(screen, "Hello, Walker AI!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return outsideWidth, outsideHeight
}

func Run() {
    game := &Game{}
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
