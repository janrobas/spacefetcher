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

	gameImages.HexRoad.Fill(color.RGBA{R: 40, G: 70, B: 230, A: 255})
	gameImages.HexRoadFast.Fill(color.RGBA{R: 23, G: 200, B: 23, A: 255})
	gameImages.HexSpace.Fill(color.RGBA{R: 20, B: 20, G: 20, A: 255})
	gameImages.HexDanger.Fill(color.RGBA{B: 23, G: 23, R: 200, A: 150})
	gameImages.HexFuel.Fill(color.RGBA{B: 200, G: 50, R: 200, A: 200})
	gameImages.EmptyImage.Fill(color.RGBA{B: 200, G: 200, R: 200, A: 200})

	return gameImages
}

func main() {
	gameState := &models.GameState{}
	menuState := &models.MenuState{}
	gameImages := getGameImages()

	lev1 := stringToGameMap("000000000\n001111100\n00S000000\n001111100\n001X00000\n001111100")
	lev2 := stringToGameMap("0000X0000\n001111100\n00S000000\n001111100\n001X00X00\n001111100")
	lev3 := stringToGameMap("1000000\n0011100\n0011100\n0011100\n0X000X0\n1110000\n0011100\n0011100\n0011100\n0X01000\n0000000")
	lev4 := stringToGameMap("000001111111XX0000\n000000000111000000\n00SX01100010000000\n000001100010000011\n00000111111X000000\n000001100000000000\n000001111111110000")
	lev5 := stringToGameMap("2100000000222222222\n011000000SXX0000022\n110000X0112222222\n0002222220000000022\n0000000020000000022\n0000000001100000022\n1122222212000222222\n0110000000000220000\n0X11000000000220000\n0001222210000220000\n00000X0010000220000\n0000000002220220000\n0000000002200220000\n0000000002200220000\n0000000002222220000")
	lev6 := stringToGameMap("0000000000000S110000\n00000000000001122X00\n00000000000002200000\n00002220000002200000\n00002220X000022X0\n00002220000002200000\n00000222222222000000\n00000000220000000000\n00000000220000000000\n00000000220000000000\n000000002222\n00000000000002200000\n00000100000002200000\n00002220000002200000\n00002220000002200000\n0000222000000220000")
	lev7 := stringToGameMap("0000000000000S220000\n00000000000001X22000\n00000000000002200000\n00002220000002200000\n000022200000022X0\n00002220X00002200000\n00000222222222000000\n0000000022X000000000\n0000000X22X000000000\n00000000220000000000\n000000002222\n00000000000002200000\n00000100000002200000\n00002220000002200000\n00002220000002200000\n0000222000000220000")
	lev8 := stringToGameMap("0000000000000S110000\n00000000000001122X00\n00000000000002200000\n00002220000002200000\n00002220X000022X0\n00002220000002200000\n00000222222222000000\n0000000022X000000000\n00000000220000000000\n00000000220000000000\n000000002222\n00000000000002200000\n00000100000002200000\n00002220000002200000\n00002220000002200000\n0000222000000220000")

	gameState.Maps = make([][][]rune, 0)

	gameState.Maps = append(gameState.Maps, lev1)
	gameState.Maps = append(gameState.Maps, lev2)
	gameState.Maps = append(gameState.Maps, lev3)
	gameState.Maps = append(gameState.Maps, lev4)
	gameState.Maps = append(gameState.Maps, lev5)
	gameState.Maps = append(gameState.Maps, lev6)
	gameState.Maps = append(gameState.Maps, lev7)
	gameState.Maps = append(gameState.Maps, lev8)

	gs := logic.RunMenu(menuState)

	go func() {
		<-gs.Stop
		logic.RunGame(gameState)
	}()

	screenProxy, _ := ebiten.NewImage(constants.ScreenWidth, constants.ScreenHeight, ebiten.FilterLinear)

	trueUpdate := func(screen *ebiten.Image) error {
		screenProxy.Fill(color.RGBA{B: 10, G: 10, R: 0, A: 255})
		result := update(screenProxy, menuState, gameState, gameImages)
		screen.DrawImage(screenProxy, &ebiten.DrawImageOptions{})
		return result
	}

	if err := ebiten.Run(trueUpdate, constants.ScreenWidth, constants.ScreenHeight, 1, "Space Fetcher"); err != nil {
		log.Fatal(err)
	}
}
