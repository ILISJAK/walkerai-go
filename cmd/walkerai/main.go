// main.go
package main

import (
	"log"
	"walkerai-go/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("WalkerAI Ragdoll")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
