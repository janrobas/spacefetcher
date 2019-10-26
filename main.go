package main

import (
	"janrobas/spaceship/constants"
	"janrobas/spaceship/logic"
	"janrobas/spaceship/models"
	"log"

	"github.com/hajimehoshi/ebiten"
)

var gameState *models.GameState

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// different states
	logic.UpdateGame(screen, gameState)

	return nil
}

func main() {
	// init game state

	gameState = &models.GameState{}

	gameState.Map = [][]rune{
		[]rune("00XXX01"),
		[]rune("0X11001"),
		[]rune("0010001X"),
		[]rune("0X10001"),
		[]rune("0X10011"),
		[]rune("0X11011"),
		[]rune("1111111"),
	}

	logic.RunGame(gameState)

	if err := ebiten.Run(update, constants.ScreenWidth, constants.ScreenHeight, 1, "Beng"); err != nil {
		log.Fatal(err)
	}

}
