package logic

import (
	"fmt"
	"math"
	"time"

	"janrobas/spacefetcher/constants"
	graphics "janrobas/spacefetcher/graphics"
	"janrobas/spacefetcher/models"

	"github.com/hajimehoshi/ebiten"
)

func initialize(state *models.GameState) {
	state.GameRunning = false

	state.Countdown = 3
	state.CountdownTs = time.Now().Unix() + 1
	state.Fuel = 100

	state.ShipX = constants.ScreenWidth / 2
	state.ShipY = constants.ScreenHeight / 2
	state.MoveXOffset = 0
	state.MoveYOffset = 0

	state.IndexXOffset = 0
	state.IndexYOffset = 0

	state.CurrentMap = models.GameMap{
		Name: state.Maps[state.CurrentMapIndex].Name,
		Map:  make([][]rune, len(state.Maps[state.CurrentMapIndex].Map)),
	}

	for i, row := range state.Maps[state.CurrentMapIndex].Map {
		state.CurrentMap.Map[i] = make([]rune, len(row))
		for j, val := range row {
			state.CurrentMap.Map[i][j] = val

			if val == 'S' {
				state.IndexXOffset = i - constants.HexesX/2 - 1
				state.IndexYOffset = j - constants.HexesY/2 - 3
			}
		}
	}

	state.ItemsLeft = 0
	for _, row := range state.CurrentMap.Map {
		for _, val := range row {
			if val == 'X' {
				state.ItemsLeft = state.ItemsLeft + 1
			}
		}
	}

	state.GameRunning = true
}

func shipIsClose(state *models.GameState, x float32, y float32, dist float64) bool {
	return math.Abs(float64(x)-float64(state.ShipX)) < dist/1.5 &&
		math.Abs(float64(y)-float64(state.ShipY)) < dist/1.5
}

func drawHexAndReturnIfCollision(screen *ebiten.Image, state *models.GameState, w float32, h float32, x float32, y float32, logicalX int, logicalY int, images *models.GameImages) bool {
	isOnMap := logicalX >= 0 && logicalY >= 0 &&
		len(state.CurrentMap.Map) > logicalX &&
		len(state.CurrentMap.Map[logicalX]) > logicalY

	var img *ebiten.Image

	shipIsOnThisHex := shipIsClose(state, x, y, constants.HexSize)

	if isOnMap && state.CurrentMap.Map[logicalX][logicalY] == 'X' {
		img = images.HexFuel
	} else if isOnMap && (state.CurrentMap.Map[logicalX][logicalY] == '1' || state.CurrentMap.Map[logicalX][logicalY] == 'S') {
		img = images.HexRoad
	} else if isOnMap && state.CurrentMap.Map[logicalX][logicalY] == '2' {
		img = images.HexRoadFast
	} else if shipIsOnThisHex {
		img = images.HexDanger
	} else {
		img = images.HexSpace
	}

	graphics.DrawHex(screen, w, h, x, y, img)

	return shipIsOnThisHex && isOnMap
}

func UpdateGame(screen *ebiten.Image, state *models.GameState, images *models.GameImages) error {
	if ebiten.IsDrawingSkipped() || !state.GameRunning {
		return nil
	}

	currentCollisions := make([]models.IntCoordinates, 0)
	for i := 0; i < constants.HexesX; i++ {
		for j := 0; j < constants.HexesY; j++ {

			offset := -4 * constants.HexSize
			if j%2 == 1 {
				offset = offset - constants.HexSize + constants.HexSize/2 - constants.HexMarginLeft + constants.HexMarginLeft/2
			}

			logicalX := i + state.IndexXOffset
			logicalY := constants.HexesY - j + state.IndexYOffset

			isCollision := drawHexAndReturnIfCollision(screen, state, constants.HexSize, constants.HexSize,
				float32(+offset+i*(constants.HexSize+constants.HexMarginLeft))+state.MoveXOffset,
				float32(-(2*constants.HexSize)+j*(constants.HexSize+constants.HexMarginTop))+state.MoveYOffset,
				logicalX, logicalY, images)

			if isCollision {
				currentCollisions = append(currentCollisions, models.IntCoordinates{X: logicalX, Y: logicalY})
			}
		}
	}

	state.CurrentCollisions = currentCollisions

	if state.ItemsLeft == 0 && state.CurrentMapIndex == len(state.Maps)-1 {
		graphics.DarkScreen(screen)

		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2,
			"Congratulations, space cadet,")

		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2+30,
			"you've cleared all stages!")

		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2+70,
			fmt.Sprintf("Your final score: %d", state.Score+int(math.Ceil(float64(state.Fuel)))))

		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2+110,
			"Always keep searching and never get lost!")

		return nil
	}

	graphics.DrawShip(screen, 40, 50, state.ShipX, state.ShipY, state.ShipRotation, images.EmptyImage)
	graphics.DrawFuel(screen, float64(state.Fuel))

	graphics.DisplayMessage(screen, 20, 22, fmt.Sprintf("Stage: %s", state.CurrentMap.Name))

	if state.Countdown > 0 {
		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2,
			fmt.Sprintf("Position the ship. Starting in %d sec.", state.Countdown))
	}

	if state.Fuel <= 0 {
		graphics.DarkScreen(screen)
		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2,
			"NO FUEL. Press SPACE to restart.")
	}

	if state.ItemsLeft == 0 {
		graphics.DarkScreen(screen)
		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight/2,
			"Stage cleared! Press SPACE to continue.")
	}

	graphics.DisplayMessage(screen, 20,
		constants.ScreenHeight-20,
		fmt.Sprintf("%d left to pick", state.ItemsLeft))

	if state.Score != 0 {
		graphics.DisplayMessage(screen, 20,
			constants.ScreenHeight-50,
			fmt.Sprintf("Score: %d", state.Score))
	}

	return nil
}

func processKeyboardActions(state *models.GameState, gameaudio *models.GameAudio, delta int64) {
	gameFinished := state.ItemsLeft == 0 && state.CurrentMapIndex == len(state.Maps)-1

	if gameFinished {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		state.ShipRotation = state.ShipRotation + float32(delta)/599
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		state.ShipRotation = state.ShipRotation - float32(delta)/599
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) &&
		(state.ItemsLeft == 0 && state.Countdown == 0 || state.Fuel <= 0) {
		if state.ItemsLeft == 0 && state.Countdown == 0 {
			state.CurrentMapIndex = state.CurrentMapIndex + 1
			state.Score += int(math.Ceil(float64(state.Fuel)))
		}

		initialize(state)
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) && state.Countdown == 0 {
		// restart
		initialize(state)
	}

	if ebiten.IsKeyPressed(ebiten.KeyM) {
		// this is checked on every tick in game loop
		// we want to wait a little before changing state, so we don't do mute-unmute-mute-unmute loop
		tMuteChange := time.NewTimer(time.Millisecond * 500)

		if gameaudio.Muted {
			gameaudio.Theme.SetVolume(constants.MusicVolume)
			gameaudio.Pick.SetVolume(constants.FxVolume)

			go func() {
				<-tMuteChange.C
				gameaudio.Muted = false
			}()
		} else {
			gameaudio.Theme.SetVolume(0)
			gameaudio.Pick.SetVolume(0)

			go func() {
				<-tMuteChange.C
				gameaudio.Muted = true
			}()
		}
	}

	if state.ShipRotation > math.Pi*2 {
		state.ShipRotation = 0
	}
}

func moveTerrain(state *models.GameState, delta int64) {
	terrainSpeedMultiply := float32(1)

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		for _, collision := range state.CurrentCollisions {
			collisionChar := state.CurrentMap.Map[collision.X][collision.Y]
			if collisionChar == '2' {
				terrainSpeedMultiply = float32(2)
				break
			}
		}
	}

	state.MoveXOffset -= float32(math.Sin(float64(state.ShipRotation))) * (float32(delta) / 7.5) * terrainSpeedMultiply
	state.MoveYOffset += float32(math.Cos(float64(state.ShipRotation))) * (float32(delta) / 7.5) * terrainSpeedMultiply

	if math.Abs(float64(state.MoveXOffset)) > 3*constants.HexSize+3*constants.HexMarginLeft {
		if state.MoveXOffset < 0 {
			state.IndexXOffset = state.IndexXOffset + 3
		} else {
			state.IndexXOffset = state.IndexXOffset - 3
		}
		state.MoveXOffset = 0
	}
	if math.Abs(float64(state.MoveYOffset)) > 2*constants.HexSize+2*constants.HexMarginTop {
		if state.MoveYOffset < 0 {
			state.IndexYOffset = state.IndexYOffset - 2
		} else {
			state.IndexYOffset = state.IndexYOffset + 2
		}
		state.MoveYOffset = 0
	}
}

func updateFuel(state *models.GameState, gameaudio *models.GameAudio, delta int64) {
	if len(state.CurrentCollisions) == 0 {
		state.Fuel -= float32(delta) / 166
	}

	for _, collision := range state.CurrentCollisions {
		collisionChar := state.CurrentMap.Map[collision.X][collision.Y]
		if collisionChar == 'X' {
			state.CurrentMap.Map[collision.X][collision.Y] = '0'
			state.ItemsLeft = state.ItemsLeft - 1
			state.Fuel += 10
			gameaudio.Pick.Rewind()
			gameaudio.Pick.Play()
		}

		if collisionChar == '0' {
			state.Fuel -= 0.09
		}
	}

	state.Fuel -= float32(delta) / 1444
}

func RunGame(state *models.GameState, gameaudio *models.GameAudio) *GameLoop {
	initialize(state)
	state.Score = 0

	gameaudio.Pick.SetVolume(constants.FxVolume)
	gameaudio.Theme.SetVolume(constants.MusicVolume)

	gameLoop := PrepareGameLoop(func(delta int64, gl *GameLoop) {
		if !gameaudio.Theme.IsPlaying() {
			gameaudio.Theme.Play()
		}

		if !state.GameRunning {
			return
		}

		processKeyboardActions(state, gameaudio, delta)

		if state.Countdown != 0 {
			if time.Now().Unix()-state.CountdownTs >= 1 {
				state.Countdown = state.Countdown - 1
				state.CountdownTs = time.Now().Unix()
			}
			return
		}

		if state.Fuel <= 0 {
			return
		}

		if state.ItemsLeft == 0 {
			return
		}

		moveTerrain(state, delta)

		updateFuel(state, gameaudio, delta)
	}, func() {
		gameaudio.Theme.Close()
	})

	go StartGameLoop(gameLoop)

	return gameLoop
}
