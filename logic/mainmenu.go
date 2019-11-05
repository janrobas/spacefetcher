package logic

import (
	"image/color"
	"janrobas/spacefetcher/constants"
	graphics "janrobas/spacefetcher/graphics"
	"janrobas/spacefetcher/models"

	"github.com/hajimehoshi/ebiten"
)

func UpdateMenu(screen *ebiten.Image, state *models.MenuState, gameImages *models.GameImages) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	graphics.DrawTitle(screen, 325, 80, "SPACE FETCHER", color.White)

	var yPos float32 = 130
	graphics.DrawShip(screen, 40, 50, 50, yPos, state.ShipRotation, gameImages.EmptyImage, constants.ShipColor)
	graphics.DrawText(screen, 40+constants.HexSize+25, yPos+constants.HexSize/2+5, "Use arrow keys to steer the ship.", color.White)

	yPos += 100
	graphics.DrawHex(screen, constants.HexSize, constants.HexSize, 40, yPos, gameImages.EmptyImage, constants.HexRoadColor)
	graphics.DrawText(screen, 40+constants.HexSize+25, yPos+constants.HexSize/2+5, "Use less fuel by staying on the blue path.", color.White)

	yPos += 100
	graphics.DrawHex(screen, constants.HexSize, constants.HexSize, 40, yPos, gameImages.EmptyImage, constants.HexFuelColor)
	graphics.DrawText(screen, 40+constants.HexSize+25, yPos+constants.HexSize/2+5, "Pick up all violet stuff. They also increase fuel.", color.White)

	yPos += 100
	graphics.DrawHex(screen, constants.HexSize, constants.HexSize, 40, yPos, gameImages.EmptyImage, constants.HexRoadFastColor)
	graphics.DrawText(screen, 40+constants.HexSize+25, yPos+constants.HexSize/2+5, "You can accelerate on green path (arrow key up).", color.White)

	yPos += 140
	graphics.DrawText(screen, 40, yPos, "You get more points by using less fuel. Ship flies by itself.", color.White)

	yPos += 100
	graphics.DrawText(screen, 40, yPos, "Press SPACE to start!", color.RGBA{B: 200, G: 50, R: 200, A: 200})

	return nil
}

func checkKeyPress(gl *GameLoop, state *models.MenuState) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		state.StartGame = true
		gl.Stop <- true
	}
}

func RunMenu(state *models.MenuState) *GameLoop {
	state.StartGame = false

	gameLoop := PrepareGameLoop(func(delta int64, gameLoop *GameLoop) {
		checkKeyPress(gameLoop, state)

		state.ShipRotation = state.ShipRotation + float32(delta)/599
	}, func() {
	})

	go StartGameLoop(gameLoop)

	return gameLoop
}
