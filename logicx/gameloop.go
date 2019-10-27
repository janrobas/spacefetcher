package logic

import (
	"time"
)

type GameLoop struct {
	onUpdate func(int64)
	onStop   func()
	Stop     chan bool
}

func PrepareGameLoop(onUpdate func(int64), onStop func()) *GameLoop {
	return &GameLoop{
		onUpdate: onUpdate,
		onStop:   onStop,
		Stop:     make(chan bool),
	}
}

func StartGameLoop(gameLoop *GameLoop) {
	ticker := time.NewTicker(15 * time.Millisecond)

	tsStart := time.Now().UnixNano() / 1000000

	for {
		select {
		case <-ticker.C:
			delta := int64(time.Now().UnixNano()/1000000 - tsStart)
			tsStart = time.Now().UnixNano() / 1000000
			gameLoop.onUpdate(delta)
		case <-gameLoop.Stop:
			gameLoop.onStop()
			ticker.Stop()
		}
	}
}
