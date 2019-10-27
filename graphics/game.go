package graphics

import (
	"image/color"
	"math"

	"janrobas/spaceship/fonts"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var (
	arcadeFont font.Face
)

func init() {
	tt, _ := truetype.Parse(fonts.ArcadeN_ttf)

	const (
		arcadeFontSize = 16
		dpi            = 72
	)

	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    arcadeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
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

	screen.DrawTriangles(vs, indices, img, op)

	//text.Draw(screen, fmt.Sprintf("%d %d", logicalX, logicalY), arcadeFont, int(x), int(y+centerY), color.White)
}

func DrawShip(screen *ebiten.Image, w float32, h float32, x float32, y float32, rotation float32, img *ebiten.Image) {
	centerX := w / 2

	vs := []ebiten.Vertex{}

	vs = append(vs, makeShipVertex(centerX, 0, w, h, x, y, rotation))
	vs = append(vs, makeShipVertex(w, h, w, h, x, y, rotation))
	vs = append(vs, makeShipVertex(0, h, w, h, x, y, rotation))

	op := &ebiten.DrawTrianglesOptions{}
	indices := []uint16{}
	indices = append(indices, 0, 1, 2)
	screen.DrawTriangles(vs, indices, img, op)
}

func DisplayMessage(screen *ebiten.Image, x float32, y float32, message string) {
	text.Draw(screen, message, arcadeFont, int(x), int(y), color.White)
}

func DrawFuel(screen *ebiten.Image, w float64, h float64, value float64) {
	ebitenutil.DrawRect(screen, 0, 0, w*value/100, h, color.RGBA{B: 200, G: 50, R: 200, A: 200})
}
