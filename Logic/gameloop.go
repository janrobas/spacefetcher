package logic

import (
	"time"
)

type GameLoop struct {
	onUpdate func()
	onStop   func()
	Stop     chan bool
}

func PrepareGameLoop(onUpdate func(), onStop func()) *GameLoop {
	return &GameLoop{
		onUpdate: onUpdate,
		onStop:   onStop,
		Stop:     make(chan bool),
	}
}

func StartGameLoop(gameLoop *GameLoop) {
	ticker := time.NewTicker(16 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			gameLoop.onUpdate()

		case <-gameLoop.Stop:
			gameLoop.onStop()
			ticker.Stop()
		}
	}
}
