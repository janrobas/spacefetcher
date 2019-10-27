package logic

import (
	"image/color"
	"math"

	"janrobas/spaceship/constants"
	"janrobas/spaceship/fonts"
	graphics "janrobas/spaceship/graphics"
	"janrobas/spaceship/models"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

var (
	arcadeFont    font.Face
	emptyImage, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexRoad, _    = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexRoadFar, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexSpace, _   = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexDanger, _  = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexFuel, _    = ebiten.NewImage(16, 16, ebiten.FilterDefault)
)

func initialize(state *models.GameState) {
	tt, _ := truetype.Parse(fonts.ArcadeN_ttf)

	const (
		arcadeFontSize = 8
		dpi            = 72
	)

	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    arcadeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	hexRoad.Fill(color.RGBA{R: 23, G: 23, B: 200, A: 200})
	hexRoadFar.Fill(color.RGBA{R: 23, G: 23, B: 200, A: 150})
	hexSpace.Fill(color.RGBA{R: 8, B: 10, G: 10, A: 255})
	hexDanger.Fill(color.RGBA{B: 23, G: 23, R: 200, A: 200})
	hexFuel.Fill(color.RGBA{B: 200, G: 50, R: 200, A: 200})
	emptyImage.Fill(color.RGBA{B: 200, G: 200, R: 200, A: 200})

	state.ShipX = constants.ScreenWidth / 2
	state.ShipY = constants.ScreenHeight / 2
}

func shipIsClose(state *models.GameState, x float32, y float32, dist float64) bool {
	return math.Abs(float64(x)-float64(state.ShipX)) < dist/1.5 &&
		math.Abs(float64(y)-float64(state.ShipY)) < dist/1.5
}

func drawHex(screen *ebiten.Image, state *models.GameState, w float32, h float32, x float32, y float32, logicalX int, logicalY int) {
	isOnMap := logicalX/4 >= 0 && logicalY/4 >= 0 &&
		len(state.Map) > logicalX/4 &&
		len(state.Map[logicalX/4]) > logicalY/4

	var img *ebiten.Image

	if isOnMap && state.Map[logicalX/4][logicalY/4] == 'X' && logicalX%4 == 2 && logicalY%4 == 2 {
		img = hexFuel
		if shipIsClose(state, x, y, constants.HexSize/1.5) {
			state.Map[logicalX/4][logicalY/4] = '0'
		}
	} else if isOnMap && state.Map[logicalX/4][logicalY/4] == '1' {
		if shipIsClose(state, x, y, constants.HexSize*5) {
			img = hexRoad
		} else {
			img = hexRoadFar
		}
	} else if shipIsClose(state, x, y, constants.HexSize/1.5) {
		// ship here
		img = hexDanger
	} else {
		img = hexSpace
	}

	graphics.DrawHex(screen, w, h, x, y, logicalX, logicalY, img)
}

func UpdateGame(screen *ebiten.Image, state *models.GameState) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {

			offset := -4*constants.HexSize + 10
			if j%2 == 1 {
				offset = offset - constants.HexSize + 24
			}

			drawHex(screen, state, constants.HexSize, constants.HexSize,
				float32(+offset+i*(constants.HexSize+4))+state.MoveXOffset,
				float32(-(2*constants.HexSize)+j*(constants.HexSize-9))+state.MoveYOffset,
				20-j+state.IndexYOffset, i+state.IndexXOffset)

		}
	}

	graphics.DrawShip(screen, 30, 40, state.ShipX, state.ShipY, state.ShipRotation)

	return nil
}

func rotateShip(state *models.GameState) {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		state.ShipRotation = state.ShipRotation + 0.03
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		state.ShipRotation = state.ShipRotation - 0.03
	}

	if state.ShipRotation > math.Pi*2 {
		state.ShipRotation = math.Pi * 2
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

func RunGame(state *models.GameState) *GameLoop {
	initialize(state)

	gameLoop := PrepareGameLoop(func() {
		rotateShip(state)
		moveTerrain(state)
	}, func() {
	})

	go StartGameLoop(gameLoop)

	return gameLoop
}
