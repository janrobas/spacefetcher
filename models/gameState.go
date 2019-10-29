package models

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
	Countdown         int
	CountdownTs       int64
	ItemsLeft         int
	CurrentMapIndex   int
	Maps              [][][]rune
	GameRunning       bool
}

type IntCoordinates struct {
	X int
	Y int
}
