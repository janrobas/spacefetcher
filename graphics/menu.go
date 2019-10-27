package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

func DrawText(screen *ebiten.Image, x float32, y float32, message string, chosenColor color.Color) {
	text.Draw(screen, message, arcadeFont, int(x), int(y), chosenColor)
}
