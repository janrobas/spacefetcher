package graphics

import (
	"janrobas/spacefetcher/fonts"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	mainFont font.Face
)

func init() {
	tt, _ := truetype.Parse(fonts.MPlus1pRegular_ttf)

	mainFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}
