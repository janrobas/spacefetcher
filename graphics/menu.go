package graphics

import (
	"image/color"
	"janrobas/spacefetcher/constants"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

func DrawText(screen *ebiten.Image, x float32, y float32, message string, chosenColor color.Color) {
	text.Draw(screen, message, mainFont, int(x), int(y), chosenColor)
}

func DrawTitle(screen *ebiten.Image, x float32, y float32, message string, chosenColor color.Color) {
	text.Draw(screen, message, titleFont, int(x), int(y), constants.HexFuelColor)
	text.Draw(screen, message, titleFont, int(x)+1, int(y)+1, constants.HexRoadColor)
	text.Draw(screen, message, titleFont, int(x)+2, int(y)+2, constants.HexRoadFastColor)
	text.Draw(screen, message, titleFont, int(x)+3, int(y)+3, color.White)
}
