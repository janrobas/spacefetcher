package main

import (
	"fmt"
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

	gameState.Map = [][]rune{
		/*[]rune("00XXX01"),
		[]rune("0X11001"),
		[]rune("0010001X"),
		[]rune("0X10001"),
		[]rune("0X10011"),
		[]rune("0X11011"),
		[]rune("1111111"),*/

		[]rune("000000000"),
		[]rune("00000XX00"),
		[]rune("000001100"),
		[]rune("000001100"),
		[]rune("011001100"),
		[]rune("001111000"),
		[]rune("000000000"),
	}

	gameState.Map = stringToGameMap("00X000000\n001111100\n001000000\n001111100\n0010X0X00\n001111100")

	fmt.Println(gameState.Map)
	//panic("nc")
	gs := logic.RunGame(gameState)
	gs.Stop <- true

	gs = logic.RunGame(gameState)

	if err := ebiten.Run(update, constants.ScreenWidth, constants.ScreenHeight, 1, "Beng"); err != nil {
		log.Fatal(err)
	}
}
