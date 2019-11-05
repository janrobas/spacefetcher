package constants

import "image/color"

const (
	ScreenWidth   = 1080
	ScreenHeight  = 720
	HexesY        = 20
	HexesX        = 25
	HexSize       = 80
	HexMarginTop  = -10
	HexMarginLeft = 5
	MusicVolume   = 0.25
	FxVolume      = 0.55
)

var (
	HexRoadColor     = &color.RGBA{R: 50, G: 90, B: 240, A: 255}
	HexRoadFastColor = &color.RGBA{R: 23, G: 200, B: 23, A: 255}
	HexSpaceColor    = &color.RGBA{R: 20, B: 20, G: 20, A: 255}
	HexFuelColor     = &color.RGBA{B: 200, G: 50, R: 200, A: 255}
	HexDangerColor   = &color.RGBA{B: 15, G: 15, R: 120, A: 255}
	ShipColor        = &color.RGBA{B: 255, G: 255, R: 255, A: 220}
)
