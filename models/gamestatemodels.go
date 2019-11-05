package models

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
)

type GameState struct {
	MoveXOffset       float32
	MoveYOffset       float32
	IndexXOffset      int
	IndexYOffset      int
	ShipX             float32
	ShipY             float32
	CurrentMap        GameMap
	CurrentCollisions []IntCoordinates
	Fuel              float64
	ShipRotation      float32
	Countdown         int
	CountdownTs       int64
	ItemsLeft         int
	CurrentMapIndex   int
	Maps              []GameMap
	GameRunning       bool
	Score             int
}

type GameImages struct {
	EmptyImage *ebiten.Image
}

type GameAudio struct {
	Theme *audio.Player
	Pick  *audio.Player
	Muted bool
}

type GameMap struct {
	Map  [][]rune
	Name string
}

type IntCoordinates struct {
	X int
	Y int
}
