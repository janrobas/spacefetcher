package logic

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"janrobas/spaceship/constants"
	graphics "janrobas/spaceship/graphics"
	"janrobas/spaceship/models"

	"github.com/hajimehoshi/ebiten"
)

func initialize(state *models.GameState) {
	state.Countdown = 3
	state.Fuel = 100

	state.GameImages = models.GameImages{}

	state.GameImages.HexDanger, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	state.GameImages.HexRoad, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	state.GameImages.HexRoadFar, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	state.GameImages.HexSpace, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	state.GameImages.HexFuel, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	state.GameImages.EmptyImage, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)

	state.GameImages.HexRoad.Fill(color.RGBA{R: 23, G: 23, B: 200, A: 200})
	state.GameImages.HexRoadFar.Fill(color.RGBA{R: 23, G: 23, B: 200, A: 150})
	state.GameImages.HexSpace.Fill(color.RGBA{R: 8, B: 10, G: 10, A: 255})
	state.GameImages.HexDanger.Fill(color.RGBA{B: 23, G: 23, R: 200, A: 200})
	state.GameImages.HexFuel.Fill(color.RGBA{B: 200, G: 50, R: 200, A: 200})
	state.GameImages.EmptyImage.Fill(color.RGBA{B: 200, G: 200, R: 200, A: 200})

	state.ShipX = constants.ScreenWidth / 2
	state.ShipY = constants.ScreenHeight / 2

	state.ItemsLeft = 0
	for _, row := range state.Map {
		for _, val := range row {
			if val == 'X' {
				state.ItemsLeft = state.ItemsLeft + 1
			}
		}
	}
}

func shipIsClose(state *models.GameState, x float32, y float32, dist float64) bool {
	return math.Abs(float64(x)-float64(state.ShipX)) < dist/1.5 &&
		math.Abs(float64(y)-float64(state.ShipY)) < dist/1.5
}

func drawHexAndReturnIfCollision(screen *ebiten.Image, state *models.GameState, w float32, h float32, x float32, y float32, logicalX int, logicalY int) bool {
	isOnMap := logicalX >= 0 && logicalY >= 0 &&
		len(state.Map) > logicalX &&
		len(state.Map[logicalX]) > logicalY

	var img *ebiten.Image

	shipIsOnThisHex := shipIsClose(state, x, y, constants.HexSize)

	if isOnMap && state.Map[logicalX][logicalY] == 'X' {
		img = state.GameImages.HexFuel
	} else if isOnMap && state.Map[logicalX][logicalY] == '1' {
		img = state.GameImages.HexRoad
	} else if shipIsOnThisHex {
		img = state.GameImages.HexDanger
	} else {
		img = state.GameImages.HexSpace
	}

	graphics.DrawHex(screen, w, h, x, y, logicalX, logicalY, img)

	return shipIsOnThisHex && isOnMap
}

func UpdateGame(screen *ebiten.Image, state *models.GameState) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	currentCollisions := make([]models.IntCoordinates, 0)
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {

			offset := -4*constants.HexSize + 10
			if j%2 == 1 {
				offset = offset - constants.HexSize + 24
			}

			logicalX := i + state.IndexXOffset
			logicalY := 20 - j + state.IndexYOffset

			isCollision := drawHexAndReturnIfCollision(screen, state, constants.HexSize, constants.HexSize,
				float32(+offset+i*(constants.HexSize+4))+state.MoveXOffset,
				float32(-(2*constants.HexSize)+j*(constants.HexSize-9))+state.MoveYOffset,
				logicalX, logicalY)

			if isCollision {
				currentCollisions = append(currentCollisions, models.IntCoordinates{X: logicalX, Y: logicalY})
			}
		}
	}

	state.CurrentCollisions = currentCollisions

	graphics.DrawShip(screen, 30, 40, state.ShipX, state.ShipY, state.ShipRotation, state.GameImages.EmptyImage)
	graphics.DrawFuel(screen, constants.ScreenWidth, 20, float64(state.Fuel))

	if state.Countdown > 0 {
		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2,
			fmt.Sprintf("Starting in %d sec.", state.Countdown))
	}

	if state.Fuel <= 0 {
		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2,
			"NO FUEL.")
	}

	if state.ItemsLeft == 0 {
		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2,
			"Stage cleared! Press SPACE to continue.")
	}

	graphics.DisplayMessage(screen, 20,
		constants.ScreenHeight-20,
		fmt.Sprintf("%d left", state.ItemsLeft))

	return nil
}

func rotateShip(state *models.GameState) {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		state.ShipRotation = state.ShipRotation + 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		state.ShipRotation = state.ShipRotation - 0.025
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) && state.ItemsLeft == 0 && state.Countdown == 0 {
		initialize(state)
	}

	if state.ShipRotation > math.Pi*2 {
		state.ShipRotation = 0
	}
}

func moveTerrain(state *models.GameState) {
	state.MoveXOffset -= float32(math.Sin(float64(state.ShipRotation))) * 2
	state.MoveYOffset += float32(math.Cos(float64(state.ShipRotation))) * 2

	if math.Abs(float64(state.MoveXOffset)) > 160 {
		if state.MoveXOffset < 0 {
			state.IndexXOffset = state.IndexXOffset + 3
		} else {
			state.IndexXOffset = state.IndexXOffset - 3
		}
		state.MoveXOffset = 0
	}
	if math.Abs(float64(state.MoveYOffset)) > 80 {
		if state.MoveYOffset < 0 {
			state.IndexYOffset = state.IndexYOffset - 2
		} else {
			state.IndexYOffset = state.IndexYOffset + 2
		}
		state.MoveYOffset = 0
	}
}

func updateFuel(state *models.GameState) {
	for _, collision := range state.CurrentCollisions {
		collisionChar := state.Map[collision.X][collision.Y]
		if collisionChar == 'X' {
			state.Map[collision.X][collision.Y] = '0'
			state.ItemsLeft = state.ItemsLeft - 1
			state.Fuel += 10
		}

		if collisionChar == '0' {
			state.Fuel -= 0.1
		}
	}

	state.Fuel -= 0.01
}

func RunGame(state *models.GameState) *GameLoop {
	initialize(state)

	ts := time.Now().Second() + 1

	gameLoop := PrepareGameLoop(func(delta int64) {
		rotateShip(state)

		if state.Countdown != 0 {
			if time.Now().Second()-ts >= 1 {
				state.Countdown = state.Countdown - 1
				ts = time.Now().Second()
			}
			return
		}

		if state.Fuel <= 0 {
			return
		}

		if state.ItemsLeft == 0 {
			return
		}

		moveTerrain(state)

		updateFuel(state)
	}, func() {
	})

	go StartGameLoop(gameLoop)

	return gameLoop
}
