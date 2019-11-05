package main

import (
	"image/color"
	"janrobas/spacefetcher/constants"
	"janrobas/spacefetcher/gameaudio"
	"janrobas/spacefetcher/logic"
	"janrobas/spacefetcher/models"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/vorbis"
)

func update(screen *ebiten.Image, menuState *models.MenuState, gameState *models.GameState, gameImages *models.GameImages) error {
	if ebiten.IsDrawingSkipped() {
		return nil
	}

	// different states
	if menuState.StartGame {
		logic.UpdateGame(screen, gameState, gameImages)
		return nil
	}

	logic.UpdateMenu(screen, menuState, gameImages)
	return nil
}

func getGameImages() *models.GameImages {
	gameImages := &models.GameImages{}

	gameImages.EmptyImage, _ = ebiten.NewImage(1, 1, ebiten.FilterDefault)
	gameImages.EmptyImage.Fill(color.RGBA{B: 255, G: 255, R: 255, A: 255})

	return gameImages
}

func getGameAudio() *models.GameAudio {
	const sampleRate = 44100

	audioContext, err := audio.NewContext(sampleRate)
	if err != nil {
		panic(err)
	}

	oggTheme, _ := vorbis.Decode(audioContext, audio.BytesReadSeekCloser(gameaudio.Theme))
	oggPick, _ := vorbis.Decode(audioContext, audio.BytesReadSeekCloser(gameaudio.Pick))

	themeLoop := audio.NewInfiniteLoop(oggTheme, (60+16.45)*4*sampleRate)
	themeAudioPlayer, _ := audio.NewPlayer(audioContext, themeLoop)
	pickAudioPlayer, _ := audio.NewPlayer(audioContext, oggPick)

	return &models.GameAudio{
		Theme: themeAudioPlayer,
		Pick:  pickAudioPlayer,
	}
}

func main() {
	lev1 := logic.StringToGameMap("000000000\n001111100\n00S000000\n001111100\n001X00000\n001111100", "First Step")
	lev2 := logic.StringToGameMap("0000X0000\n001111100\n00S000000\n001111100\n001X00X00\n002222200", "Giant Leap")
	lev3 := logic.StringToGameMap("1000000\n0011100\n0011100\n0011100\n0X000X0\n1110000\n0011100\n0011100\n0011100\n0X01000\n0000000", "Islands in Space")
	lev4 := logic.StringToGameMap("000001111111XX0000\n000000000111000000\n00SX01100010000000\n000001100010000011\n00000111111X000000\n000001100000000000\n000001111111110000", "Branches")
	lev5 := logic.StringToGameMap("2100000000222222222\n011000000SXX0000022\n110000X0112222222\n0002222220000000022\n0000000020000000022\n0000000001100000022\n1122222212000222222\n0110000000000220000\n0X11000000000220000\n0001222210000220000\n00000X0010000220000\n0000000002220220000\n0000000002200220000\n0000000002200220000\n0000000002222220000", "Connection")
	lev6 := logic.StringToGameMap("0000000000000S110000\n00000000000001122X00\n00000000000002200000\n00002220000002200000\n00002220X000022X0\n00002220000002200000\n00000222222222000000\n00000000220000000000\n00000000220000000000\n00000000220000000000\n000000002222\n00000000000002200000\n00000100000002200000\n00002220000002200000\n00002220000002200000\n0000222000000220000", "No Map")
	lev7 := logic.StringToGameMap("01000000000001000\n01022222222201000\n01020000000201000\n0102001S100201000\n0122XX111002X1000\n01220000000201000\n01022222222201000\n01000000000001000\n0100000000000XX00", "The Walls")
	lev8 := logic.StringToGameMap("0000000000000S220000\n00000000000001X22000\n00000000000002200000\n00002220000002200000\n000022200000022X0\n00002220X00002200000\n00000222222222000000\n0000000022X000000000\n0000000X22X000000000\n00000000220000000000\n000000002222\n00000000000002200000\n00000100000002200000\n00002220000002200000\n00002220000002200000\n0000222000000220000", "Nowhere")
	lev9 := logic.StringToGameMap("00000S000000000000000\n000001000000000000000\n000012000000000000000\n000002100000000000000\n000012000000000000000\n00000211X000000000000\n000002000000000000000\n000002000000000000000\n001112000000000000000\n000002000000000000000\n000002000000000000000\n00000200X000000000000\n000002111000000000000\n000002001000000000000\n0X111200X000000000000\n222222222222222222222", "A Day at the Races")
	lev10 := logic.StringToGameMap("0000000000000S110000\n00000000000001122X00\n00000000000002200000\n00002220000002200000\n00002220X000022X0\n00002220000002200000\n00000222222222000000\n0000000022X000000000\n00000000220000000000\n00000000220000000000\n000000002222\n00000000000002200000\n00000100000002200000\n00002220000002200000\n00002220000002200000\n0000222000000220000", "Back to Nowhere")
	lev11 := logic.StringToGameMap("00S100000000000000\n000111100000000000\n000022000000000000\n000022000000000000\n000022X00000000000\n000022001110111000\n00002200101X101000\n000X22001010101000\n000022111011101000\n00002200000000X000\n000022222222222000\n000022000000000000\n000022000000000000\n", "Cooling System")
	lev12 := logic.StringToGameMap("00S022200000000000\n002000200000000000\n002222222220X00000\n000002000001000000\n000002X00001000000\n000001111111000000\n00000X000000000000\n000002222222000000\n000002222222000000\n00000000000X000000\n", "Escape from Space")

	menuState := &models.MenuState{}
	gameImages := getGameImages()
	gameAudio := getGameAudio()

	gameState := &models.GameState{
		Maps: []models.GameMap{lev1, lev2, lev3, lev4, lev5, lev6, lev7, lev8, lev9, lev10, lev11, lev12},
	}

	gs := logic.RunMenu(menuState)

	go func() {
		<-gs.Stop
		logic.RunGame(gameState, gameAudio)
	}()

	screenProxy, _ := ebiten.NewImage(constants.ScreenWidth, constants.ScreenHeight, ebiten.FilterDefault)

	trueUpdate := func(screen *ebiten.Image) error {
		screenProxy.Fill(color.RGBA{B: 0, G: 0, R: 0, A: 255})
		result := update(screenProxy, menuState, gameState, gameImages)
		screen.DrawImage(screenProxy, &ebiten.DrawImageOptions{})
		return result
	}

	if err := ebiten.Run(trueUpdate, constants.ScreenWidth, constants.ScreenHeight, 1, "Space Fetcher"); err != nil {
		log.Fatal(err)
	}
}
