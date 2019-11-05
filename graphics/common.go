package graphics

import (
	"janrobas/spacefetcher/fonts"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	titleFont font.Face
	mainFont  font.Face
)

func init() {
	tt, _ := truetype.Parse(fonts.MPlus1pRegular_ttf)

	mainFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	tt, _ = truetype.Parse(fonts.ArcadeN_ttf)

	titleFont = truetype.NewFace(tt, &truetype.Options{
		Size:    34,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
