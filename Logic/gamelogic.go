package logic

import (
	"image/color"
	"log"
	"math"
	"time"

	"janrobas/spaceship/constants"
	"janrobas/spaceship/fonts"
	graphics "janrobas/spaceship/graphics"
	"janrobas/spaceship/models"
	sound "janrobas/spaceship/sound"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
	"golang.org/x/image/font"
)

var (
	arcadeFont   font.Face
	audioContext *audio.Context
	audioPlayer  *audio.Player
)

const (
	screenWidth  = 640
	screenHeight = 480
	hexSize      = 50
	sampleRate   = 44100
)

var (
	emptyImage, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexRoad, _    = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexRoadFar, _ = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexSpace, _   = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexDanger, _  = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	hexFuel, _    = ebiten.NewImage(16, 16, ebiten.FilterDefault)
	shipRotation  = float32(0)
)

func initialize() {

	var err error
	// Initialize audio context.
	audioContext, err = audio.NewContext(sampleRate)
	if err != nil {
		log.Fatal(err)
	}

	audioData := sound.GetTheme()

	if err != nil {
		log.Fatal(err)
	}

	s, err := vorbis.Decode(audioContext, audio.BytesReadSeekCloser(audioData))

	if err != nil {
		log.Fatal(err)
	}

	audioPlayer, err = audio.NewPlayer(audioContext, s)

	tt, _ := truetype.Parse(fonts.ArcadeN_ttf)

	hexRoad.Fill(color.RGBA{R: 23, G: 23, B: 200, A: 200})
	hexRoadFar.Fill(color.RGBA{R: 23, G: 23, B: 200, A: 150})
	hexSpace.Fill(color.RGBA{R: 8, B: 10, G: 10, A: 255})
	hexDanger.Fill(color.RGBA{B: 23, G: 23, R: 200, A: 200})
	hexFuel.Fill(color.RGBA{B: 200, G: 50, R: 200, A: 200})
	emptyImage.Fill(color.RGBA{B: 200, G: 200, R: 200, A: 200})

	const (
		arcadeFontSize = 8
		dpi            = 72
	)

	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    arcadeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

}

var (
	vertices []ebiten.Vertex

	ngon     = 10
	prevNgon = 0
)

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
		//screen.DrawTriangles(vs, indices, hexFuel, op)
		img = hexFuel
		if shipIsClose(state, x, y, constants.HexSize/1.5) {
			//state.Map[logicalX/4] = state.Map[logicalX/4][:logicalY/4] + "0" + state.Map[logicalX/4][logicalY/4+1:]
			state.Map[logicalX/4][logicalY/4] = '0'
		}
	} else if isOnMap && state.Map[logicalX/4][logicalY/4] == '1' {
		if shipIsClose(state, x, y, constants.HexSize*5) {
			//screen.DrawTriangles(vs, indices, hexRoad, op)
			img = hexRoad
		} else {
			//screen.DrawTriangles(vs, indices, hexRoadFar, op)
			img = hexRoadFar
		}
	} else if shipIsClose(state, x, y, constants.HexSize/1.5) {
		// ship here
		//screen.DrawTriangles(vs, indices, hexDanger, op)
		img = hexDanger
	} else {
		//screen.DrawTriangles(vs, indices, hexSpace, op)
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

	graphics.DrawShip(screen, 30, 40, state.ShipX, state.ShipY, shipRotation)

	return nil
}

func RunGame(state *models.GameState) {
	initialize()
	state.ShipX = screenWidth / 2
	state.ShipY = screenHeight / 2

	audioPlayer.Play()

	go func() {
		for {
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				shipRotation = shipRotation + 0.009
			}
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				shipRotation = shipRotation - 0.009
			}

			<-time.NewTimer(8 * time.Millisecond).C

			if shipRotation > math.Pi*2 {
				shipRotation = math.Pi * 2
			}
		}
	}()

	moveTicker := time.NewTicker(10 * time.Millisecond)
	go func() {
		for {
			<-moveTicker.C

			state.MoveXOffset -= float32(math.Sin(float64(shipRotation)))
			state.MoveYOffset += float32(math.Cos(float64(shipRotation)))

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
	}()
}
