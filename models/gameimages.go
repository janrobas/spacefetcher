package models

import "github.com/hajimehoshi/ebiten"

type GameImages struct {
	EmptyImage  *ebiten.Image
	HexRoad     *ebiten.Image
	HexRoadFast *ebiten.Image
	HexSpace    *ebiten.Image
	HexDanger   *ebiten.Image
	HexFuel     *ebiten.Image
}
