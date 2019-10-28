package logic

import (
	"image/color"
	graphics "janrobas/spacefetcher/graphics"
	"janrobas/spacefetcher/models"

	"github.com/hajimehoshi/ebiten"
)

func UpdateMenu(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	graphics.DrawText(screen, 50, 50, "Ship flies by itself.", color.White)
	graphics.DrawText(screen, 50, 100, "Left and right arrows rotate ship.", color.White)
	graphics.DrawText(screen, 50, 150, "Try to stay on blue path.", color.RGBA{R: 23, G: 23, B: 200, A: 200})
	graphics.DrawText(screen, 50, 200, "You must pick up violet stuff.", color.RGBA{B: 200, G: 50, R: 200, A: 200})
	graphics.DrawText(screen, 50, 250, "Press SPACE to start.", color.White)

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
	}, func() {
	})

	go StartGameLoop(gameLoop)

	return gameLoop
}
