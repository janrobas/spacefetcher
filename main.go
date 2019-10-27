package main

import (
	"janrobas/spaceship/constants"
	"janrobas/spaceship/logic"
	"janrobas/spaceship/models"
	"log"
	"strings"

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

func stringToGameMap(minimized string) [][]rune {
	result := make([][]rune, 0, 0)

	for _, row := range strings.Split(minimized, "\n") {

		rowRes := make([][]rune, 4)

		for n := 0; n < 4; n++ {
			rowRes[n] = make([]rune, len(row)*4)
			for i, char := range row {
				if char == 'X' {
					if n%4 != 2 {
						char = '0'
					}

					rowRes[n][4*i] = '0'
					rowRes[n][4*i+1] = '0'
					rowRes[n][4*i+2] = char
					rowRes[n][4*i+3] = '0'
				} else {
					rowRes[n][4*i] = char
					rowRes[n][4*i+1] = char
					rowRes[n][4*i+2] = char
					rowRes[n][4*i+3] = char
				}

			}
		}

		result = append(result, rowRes...)

	}

	return result
}

func main() {
	// init game state

	gameState = &models.GameState{}

	lev1 := stringToGameMap("000000000\n001111100\n001000000\n001111100\n001X00000\n001111100")
	lev2 := stringToGameMap("0000X0000\n001111100\n001000000\n001111100\n001X00X00\n001111100")

	gameState.Maps = make([][][]rune, 0)
	gameState.Maps = append(gameState.Maps, lev1)
	gameState.Maps = append(gameState.Maps, lev2)

	gs := logic.RunGame(gameState)
	gs.Stop <- true

	gs = logic.RunGame(gameState)

	if err := ebiten.Run(update, constants.ScreenWidth, constants.ScreenHeight, 1, "Beng"); err != nil {
		log.Fatal(err)
	}
}
