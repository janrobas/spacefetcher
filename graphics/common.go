package graphics

import (
	"janrobas/spacefetcher/fonts"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	arcadeFont font.Face
)

func init() {
	tt, _ := truetype.Parse(fonts.ArcadeN_ttf)

	const (
		arcadeFontSize = 16
		dpi            = 72
	)

	arcadeFont = truetype.NewFace(tt, &truetype.Options{
		Size:    arcadeFontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}
