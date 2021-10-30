package main

import (
	"math"
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var (
	allSprites sprite.SpriteGroup
	Width      int
	Height     int
	MouseX     int
	MouseY     int
	gameGrid   *Grid
)

type GameState int

const (
	GAME_INIT = iota
	GAME_READY
	GAME_STARTED
	GAME_RUNNING
	GAME_OVER
)

func setPalette() {
	sprite.ColorMap['o'] = tm.Color214
	sprite.ColorMap['y'] = tm.Color228
	sprite.ColorMap['r'] = tm.Color197
	sprite.ColorMap['d'] = tm.Color52
	sprite.ColorMap['B'] = tm.ColorBlack
	sprite.ColorMap['w'] = tm.ColorWhite
	sprite.ColorMap['t'] = tm.Color173
	sprite.ColorMap['T'] = tm.Color130
	sprite.ColorMap['g'] = tm.ColorSilver
	sprite.ColorMap['G'] = tm.ColorGray
	sprite.ColorMap['X'] = tm.ColorBlack
	sprite.ColorMap['b'] = tm.ColorBlue
	sprite.ColorMap['x'] = tm.Color187
}

func main() {
	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	w, h := tm.Size()
	Width = w * 2
	Height = h * 2

	setPalette()

	allSprites.Init(Width, Height, true)
	allSprites.Background = tm.Color187

	rand.Seed(time.Now().UnixNano())

	gameGrid = NewGrid()
	titleOverlay := NewTitleOverlay()

	eventQueue := make(chan tm.Event)
	go func() {
		for {
			eventQueue <- tm.PollEvent()
		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				allSprites.TriggerEvent("UpdateTimer")
			}
		}
	}()

mainloop:
	for {
		tm.Clear(tm.Color187, tm.Color187)

		select {
		case ev := <-eventQueue:
			if ev.Type == tm.EventKey {
				if ev.Key == tm.KeyCtrlC || ev.Key == tm.KeyEsc || ev.Ch == 'q' {
					break mainloop
				}
			} else if ev.Type == tm.EventMouse {
				MouseX = ev.MouseX * 2
				MouseY = ev.MouseY * 2
				if ev.Key == tm.MouseLeft {
					if gameGrid.State == GAME_READY {
						s := titleOverlay.CheckSelectorClicked(MouseX, MouseY)
						if s != nil {
							gameGrid.TotalBombs = int(math.Round(float64(gameGrid.Width) * float64(gameGrid.Height) * s.BombRate))
							gameGrid.State = GAME_STARTED
							allSprites.TriggerEvent("SelectorClicked")
						}
					} else if gameGrid.State == GAME_RUNNING || gameGrid.State == GAME_STARTED {
						t := gameGrid.FindTileClicked(MouseX, MouseY)
						if t != nil {
							pos := gameGrid.GetTilePos(t)
							if gameGrid.State == GAME_STARTED {
								gameGrid.PlaceBombs(t)
								allSprites.MoveToTop(gameGrid.Super)
							}
							gameGrid.RevealTileAtPos(pos)
							gameGrid.CheckGameOver()
						}
					} else if gameGrid.State == GAME_OVER {
						allSprites.TriggerEvent("ReturnToGrid")
					}
				} else if ev.Key == tm.MouseRight {
					if gameGrid.State == GAME_RUNNING {
						t := gameGrid.FindTileClicked(MouseX, MouseY)
						if t != nil && t.Covered {
							t.SetFlag()
							gameGrid.CheckGameOver()
						}
					}
				}
			} else if ev.Type == tm.EventResize {
				if ev.Width == 0 || ev.Height == 0 {
					continue
				}
				Width = ev.Width * 2
				Height = ev.Height * 2
				allSprites.Init(Width, Height, true)
				allSprites.Background = tm.Color187
				allSprites.TriggerEvent("resizeScreen")

				if gameGrid.State == GAME_INIT && Width > 80 && Height > 40 {
					gameGrid.SetReady()
					gameGrid.SetSize(Width/8, (Height-HEADER_OFFSET)/8)
					titleOverlay.SetGameReady()
					titleOverlay.MoveToTop()
				}
			}
		default:
			allSprites.Update()
			allSprites.Render()
			time.Sleep(60 * time.Millisecond)
		}
	}
}
