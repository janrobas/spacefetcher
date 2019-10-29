package main

import (
	"image/color"
	"janrobas/spacefetcher/constants"
	"janrobas/spacefetcher/logic"
	"janrobas/spacefetcher/models"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten"
)

func update(screen *ebiten.Image, menuState *models.MenuState, gameState *models.GameState, gameImages *models.GameImages) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// different states
	if menuState.StartGame {
		logic.UpdateGame(screen, gameState, gameImages)
		return nil
	}

	logic.UpdateMenu(screen, menuState, gameImages)
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

func getGameImages() *models.GameImages {
	gameImages := &models.GameImages{}

	gameImages.HexDanger, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	gameImages.HexRoad, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	gameImages.HexRoadFast, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	gameImages.HexSpace, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	gameImages.HexFuel, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	gameImages.EmptyImage, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)

	gameImages.HexRoad.Fill(color.RGBA{R: 23, G: 23, B: 200, A: 200})
	gameImages.HexRoadFast.Fill(color.RGBA{R: 23, G: 200, B: 23, A: 150})
	gameImages.HexSpace.Fill(color.RGBA{R: 8, B: 10, G: 10, A: 255})
	gameImages.HexDanger.Fill(color.RGBA{B: 23, G: 23, R: 200, A: 200})
	gameImages.HexFuel.Fill(color.RGBA{B: 200, G: 50, R: 200, A: 200})
	gameImages.EmptyImage.Fill(color.RGBA{B: 200, G: 200, R: 200, A: 200})

	return gameImages
}

func main() {
	// init game state

	gameState := &models.GameState{}
	menuState := &models.MenuState{}
	gameImages := getGameImages()

	lev1 := stringToGameMap("000000000\n002222200\n002000000\n002222200\n00SX00000\n001111100")
	//lev2 := stringToGameMap("0000X0000\n001111100\n001000000\n001111100\n001X00X00\n001111100")
	/*lev3 := stringToGameMap("1000000\n0011100\n0011100\n0011100\n0X000X0\n1110000\n0011100\n0011100\n0011100\n0X01000\n0000000")
	lev4 := stringToGameMap("000001111111XX0000\n000000000111000000\n000X01100010000000\n000001100010000011\n00000111111X000000\n000001100000000000\n000001111111110000")
	*/
	gameState.Maps = make([][][]rune, 0)
	gameState.Maps = append(gameState.Maps, lev1)
	//gameState.Maps = append(gameState.Maps, lev2)
	/*gameState.Maps = append(gameState.Maps, lev3)
	gameState.Maps = append(gameState.Maps, lev4)*/

	gs := logic.RunMenu(menuState)

	go func() {
		<-gs.Stop
		logic.RunGame(gameState)
	}()

	screenProxy, _ := ebiten.NewImage(constants.ScreenWidth, constants.ScreenHeight, ebiten.FilterLinear)

	trueUpdate := func(screen *ebiten.Image) error {
		screenProxy.Fill(color.Black)
		result := update(screenProxy, menuState, gameState, gameImages)
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(screenProxy, op)
		return result
	}

	if err := ebiten.Run(trueUpdate, constants.ScreenWidth, constants.ScreenHeight, 1, "Space Fetcher"); err != nil {
		log.Fatal(err)
	}
}
