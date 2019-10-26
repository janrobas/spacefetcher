package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"janrobas/spaceship/fonts"
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

	moveXOffset  = float32(0)
	moveYOffset  = float32(0)
	indexXOffset = 0
	indexYOffset = 0
	shipX        = float32(0)
	shipY        = float32(0)
)

func init() {

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
	hexFuel.Fill(color.RGBA{B: 23, G: 23, R: 200, A: 200})
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

	shipX = screenWidth / 2
	shipY = screenHeight / 2
}

var (
	vertices []ebiten.Vertex

	ngon     = 10
	prevNgon = 0

	level = []string{"0000001",
		"0011001",
		"0010001X",
		"0010001",
		"0010011",
		"0X11011",
		"1111111"}
)

func makeShipVertex(x float32, y float32, w float32, h float32, ox float32, oy float32, rotation float32) ebiten.Vertex {

	x = x - w/2
	y = y - h/2

	newX := x*float32(math.Cos(float64(rotation))) - y*float32(math.Sin(float64(rotation)))
	newY := x*float32(math.Sin(float64(rotation))) + y*float32(math.Cos(float64(rotation)))

	x = newX + w/2 + ox
	y = newY + h/2 + oy

	return ebiten.Vertex{
		DstX:   x,
		DstY:   y,
		SrcX:   0,
		SrcY:   0,
		ColorR: float32(1),
		ColorG: float32(1),
		ColorB: float32(1),
		ColorA: 1,
	}

}

func makeHexVertex(x float32, y float32) ebiten.Vertex {
	return ebiten.Vertex{
		DstX:   x,
		DstY:   y,
		SrcX:   0,
		SrcY:   0,
		ColorR: float32(1),
		ColorG: float32(1),
		ColorB: float32(1),
		ColorA: 1,
	}
}

func shipIsClose(x float32, y float32, dist float64) bool {
	return math.Abs(float64(x)-float64(shipX)) < dist/1.5 &&
		math.Abs(float64(y)-float64(shipY)) < dist/1.5
}

func drawHex(screen *ebiten.Image, w float32, h float32, x float32, y float32, logicalX int, logicalY int) {

	centerX := w / 2
	centerY := h / 2

	vs := []ebiten.Vertex{}

	vs = append(vs, makeHexVertex(x, y+centerY/2))
	vs = append(vs, makeHexVertex(x+centerX, y))
	vs = append(vs, makeHexVertex(x+w, y+centerY/2))
	vs = append(vs, makeHexVertex(x+w, y+centerY+centerY/2))
	vs = append(vs, makeHexVertex(x+centerX, y+h))
	vs = append(vs, makeHexVertex(x, y+h-centerY/2))

	op := &ebiten.DrawTrianglesOptions{}
	indices := []uint16{}
	indices = append(indices, 0, 1, 2, 0, 5, 3, 3, 2, 0, 5, 4, 3)

	isOnMap := logicalX/4 >= 0 && logicalY/4 >= 0 &&
		len(level) > logicalX/4 &&
		len(level[logicalX/4]) > logicalY/4

	if isOnMap && level[logicalX/4][logicalY/4] == 'X' && logicalX%4 == 2 && logicalY%4 == 2 {
		screen.DrawTriangles(vs, indices, hexFuel, op)

		if shipIsClose(x, y, hexSize/1.5) {
			level[logicalX/4] = level[logicalX/4][:logicalY/4] + "0" + level[logicalX/4][logicalY+1/4:]
		}
	} else if isOnMap && level[logicalX/4][logicalY/4] == '1' {
		if shipIsClose(x, y, hexSize*5) {
			screen.DrawTriangles(vs, indices, hexRoad, op)
		} else {
			screen.DrawTriangles(vs, indices, hexRoadFar, op)
		}
	} else if shipIsClose(x, y, hexSize/1.5) {
		// ship here
		screen.DrawTriangles(vs, indices, hexDanger, op)
	} else {
		screen.DrawTriangles(vs, indices, hexSpace, op)
	}

	//text.Draw(screen, fmt.Sprintf("%d %d", logicalX, logicalY), arcadeFont, int(x), int(y+centerY), color.White)
}

func drawShip(screen *ebiten.Image, w float32, h float32, x float32, y float32, rotation float32) {
	centerX := w / 2

	vs := []ebiten.Vertex{}

	vs = append(vs, makeShipVertex(centerX, 0, w, h, x, y, rotation))
	vs = append(vs, makeShipVertex(w, h, w, h, x, y, rotation))
	vs = append(vs, makeShipVertex(0, h, w, h, x, y, rotation))

	op := &ebiten.DrawTrianglesOptions{}
	indices := []uint16{}
	indices = append(indices, 0, 1, 2)
	screen.DrawTriangles(vs, indices, emptyImage, op)
}

func update(screen *ebiten.Image) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {

			offset := -4*hexSize + 10
			if j%2 == 1 {
				offset = offset - hexSize + 24
			}

			drawHex(screen, hexSize, hexSize,
				float32(+offset+i*(hexSize+4))+moveXOffset,
				float32(-(2*hexSize)+j*(hexSize-9))+moveYOffset,
				20-j+indexYOffset, i+indexXOffset)
		}
	}

	drawShip(screen, 30, 40, shipX, shipY, shipRotation)

	return nil
}

func main() {
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

			moveXOffset -= float32(math.Sin(float64(shipRotation)))
			moveYOffset += float32(math.Cos(float64(shipRotation)))

			if math.Abs(float64(moveXOffset)) > 160 {
				if moveXOffset < 0 {
					indexXOffset = indexXOffset + 3
				} else {
					indexXOffset = indexXOffset - 3
				}
				moveXOffset = 0
			}
			if math.Abs(float64(moveYOffset)) > 80 {
				if moveYOffset < 0 {
					indexYOffset = indexYOffset - 2
				} else {
					indexYOffset = indexYOffset + 2
				}
				moveYOffset = 0
			}
		}
	}()

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Beng"); err != nil {
		log.Fatal(err)
	}

}
