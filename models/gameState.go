package models

import (
	"github.com/hajimehoshi/ebiten"
)

type GameState struct {
	MoveXOffset       float32
	MoveYOffset       float32
	IndexXOffset      int
	IndexYOffset      int
	ShipX             float32
	ShipY             float32
	Map               [][]rune
	CurrentCollisions []IntCoordinates
	Fuel              float32
	ShipRotation      float32
	GameImages        GameImages
}

type IntCoordinates struct {
	X int
	Y int
}

type GameImages struct {
	EmptyImage *ebiten.Image
	HexRoad    *ebiten.Image
	HexRoadFar *ebiten.Image
	HexSpace   *ebiten.Image
	HexDanger  *ebiten.Image
	HexFuel    *ebiten.Image
}
