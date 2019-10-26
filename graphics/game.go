package graphics

import (
	"image/color"
	"math"

	"janrobas/spaceship/fonts"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

var (
	arcadeFont font.Face
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

func init() {
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

func DrawHex(screen *ebiten.Image, w float32, h float32, x float32, y float32, logicalX int, logicalY int, img *ebiten.Image) {

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

	//text.Draw(screen, fmt.Sprintf("%d %d", logicalX, logicalY), arcadeFont, int(x), int(y+centerY), color.White)
	screen.DrawTriangles(vs, indices, img, op)
}

func DrawShip(screen *ebiten.Image, w float32, h float32, x float32, y float32, rotation float32) {
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
