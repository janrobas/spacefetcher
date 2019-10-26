package models

type GameState struct {
	MoveXOffset  float32
	MoveYOffset  float32
	IndexXOffset int
	IndexYOffset int
	ShipX        float32
	ShipY        float32
	Map          [][]rune
}
