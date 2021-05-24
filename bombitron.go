package main

import (
	"fmt"
	"time"
	"math/rand"

	sprite "github.com/pdevine/go-asciisprite"
	tm "github.com/pdevine/go-asciisprite/termbox"
)

var allSprites sprite.SpriteGroup
var Width int
var Height int
var MouseX int
var MouseY int
var gameGrid *Grid

const (
	TILE_WIDTH = 8
	TILE_HEIGHT = 8
	HEADER_OFFSET = 10
)

const tileEmpty = `BBBBBBBB
B
B
B
B
B
B
B`

const tileCovered = `wwwwwwww
w
w      G
w      G
w      G
w      G
w      G
w GGGGGG`

const tileCoveredReverse = `GGGGGGGG
G
G      w
G      w
G      w
G      w
G      w
G wwwwww`

const tileFlag = `wwwwwwww
w       
w    r G
w  rrr G
w rrrr G
w   rB G
w    B G
w      G`

const tile1 = `BBBBBBBB
B    
B   b
B  bb
B   b
B   b
B  bbb
B`

const tile2 = `BBBBBBBB
B
B  gg
B    g
B  ggg
B  g
B  ggg
B`

const tile3 = `BBBBBBBB
B
B  rr
B    r
B  rr
B    r
B  rrr
B`

const tile4 = `BBBBBBBB
B
B  B B
B  B B
B  BBB
B    B
B    B
B`

const tile5 = `BBBBBBBB
B
B  BBB
B  B
B  BB
B    B
B  BB
B`

const tile6 = `BBBBBBBB
B
B   BB
B  B
B  BBB
B  B B
B  BB
B`

const tile7 = `BBBBBBBB
B
B  BBB
B    B
B   B
B  B
B  B
B`

const tile8 = `BBBBBBBB
B
B   BB
B  B B
B  BBB
B  B B
B  BB
B`

const tileBomb = `BBBBBBBB
BRRRRRRR
BRRRBRRR
BRRBBBRR
BRBBBBBR
BRBBBBBR
BRRBBBRR
BRRRBRRR
BRRRRRRR`

type TileType int
type GameState int

const (
	TILE_EMPTY = iota
	TILE_1
	TILE_2
	TILE_3
	TILE_4
	TILE_5
	TILE_6
	TILE_7
	TILE_8
	TILE_COVERED
	TILE_COVERED_REVERSE
	TILE_FLAG
	TILE_BOMB
)

const (
	GAME_INIT = iota
	GAME_STARTED
	GAME_RUNNING
	GAME_OVER
)

type Tile struct {
	sprite.BaseSprite
	GridX     int
	GridY     int
	BombCount int
	HaveBomb  bool
	HaveFlag  bool
	Covered   bool
}

type Grid struct {
	State          GameState
	Width          int
	Height         int
	Tiles          []*Tile
	TotalBombs     int
	FlagsRemaining *FlagsRemainingText
	TimerElapsed   *TimerElapsedText
}

type FlagsRemainingText struct {
	sprite.BaseSprite
	font      *sprite.Font
	Remaining int
}

type TimerElapsedText struct {
	sprite.BaseSprite
	font      *sprite.Font
	Started   bool
	StartTime time.Time
}

func NewFlagsRemaining() *FlagsRemainingText {
	f := &FlagsRemainingText{BaseSprite: sprite.BaseSprite{
		X: 4,
		Y: 1,
		Visible: false},
		font:      sprite.NewPakuFont(),
		Remaining: 0,
	}
	f.Init()

	f.RegisterEvent("ShowFlagsRemaining", func() {
		f.Visible = true
		f.UpdateText()
	})

	return f
}

func (f *FlagsRemainingText) UpdateText() {
	s := fmt.Sprintf("%d", f.Remaining)
	surf := sprite.NewSurfaceFromString(f.font.BuildString(s), true)
	f.BlockCostumes = []*sprite.Surface{&surf}
}

func NewTimerElapsed() *TimerElapsedText {
	t := &TimerElapsedText{BaseSprite: sprite.BaseSprite{
		X: 20,
		Y: 1,
		Visible: false},
		font:  sprite.NewPakuFont(),
	}
	t.Init()

	t.RegisterEvent("StartTimer", func() {
		t.Visible = true
		t.Started = true
		t.StartTime = time.Now()
	})

	t.RegisterEvent("UpdateTimer", func() {
		t.UpdateText()
	})

	t.RegisterEvent("Explode", func() {
		t.Started = false
	})

	t.RegisterEvent("GameOver", func() {
		t.Started = false
	})

	return t
}

func (t *TimerElapsedText) UpdateText() {
	if t.Started {
		s := fmt.Sprintf("%d", int(time.Since(t.StartTime).Seconds()))
		surf := sprite.NewSurfaceFromString(t.font.BuildString(s), true)
		t.BlockCostumes = []*sprite.Surface{&surf}
	}
}

func NewTile() *Tile {
	t := &Tile{BaseSprite: sprite.BaseSprite{
		Visible: true},
		Covered: true,
	}
	t.Init()

	t.SetTile(TILE_COVERED)
	return t
}

func (t *Tile) RevealTile() {
	if t.HaveFlag {
		return
	} else if t.HaveBomb {
		t.SetTile(TILE_BOMB)
		allSprites.TriggerEvent("Explode")
		gameGrid.State = GAME_OVER
		return
	}
	t.Covered = false
	t.SetTile(TileType(t.BombCount))
}

func (t *Tile) SetFlag() {
	if gameGrid.State != GAME_RUNNING {
		return
	}

	if t.HaveFlag {
		gameGrid.FlagsRemaining.Remaining += 1
		t.HaveFlag = false
		t.SetTile(TILE_COVERED)
		allSprites.TriggerEvent("ShowFlagsRemaining")
	} else {
		if gameGrid.FlagsRemaining.Remaining > 0 {
			gameGrid.FlagsRemaining.Remaining -= 1
			t.HaveFlag = true
			t.SetTile(TILE_FLAG)
			allSprites.TriggerEvent("ShowFlagsRemaining")
		}
	}
}

func (t *Tile) SetTile(v TileType) {
	tileImages := []string{
		tileEmpty,
		tile1,
		tile2,
		tile3,
		tile4,
		tile5,
		tile6,
		tile7,
		tile8,
		tileCovered,
		tileCoveredReverse,
		tileFlag,
		tileBomb,
	}
	surf := sprite.NewSurfaceFromString(tileImages[v], false)
	t.BlockCostumes = []*sprite.Surface{&surf}
}

func NewGrid() *Grid {
	g := &Grid{
		State:          GAME_INIT,
		FlagsRemaining: NewFlagsRemaining(),
		TimerElapsed:   NewTimerElapsed(),
	}
	allSprites.Sprites = append(allSprites.Sprites, g.FlagsRemaining)
	allSprites.Sprites = append(allSprites.Sprites, g.TimerElapsed)
	return g
}

func (g *Grid) CheckGameOver() bool {

	for _, t := range g.Tiles {
		if t.Covered && !t.HaveFlag {
			return false
		}
	}
	allSprites.TriggerEvent("GameOver")
	return true
}

func (g *Grid) SetSize(w, h, totalBombs int) {
	if g.State != GAME_INIT {
		return
	}

	g.Width = w
	g.Height = h
	g.TotalBombs = totalBombs
	g.FlagsRemaining.Remaining = totalBombs
	allSprites.TriggerEvent("ShowFlagsRemaining")

	g.Tiles = make([]*Tile, 0, 0)

	for cntY := 0; cntY < h; cntY++ {
		for cntX := 0; cntX < w; cntX++ {
			t := NewTile()
			t.X = cntX * 8
			t.Y = (cntY * 8) + HEADER_OFFSET
			g.Tiles = append(g.Tiles, t)
			allSprites.Sprites = append(allSprites.Sprites, t)
		}
	}
	g.State = GAME_STARTED
}

func (g *Grid) FindTileClicked(x, y int) *Tile {
	if (x*2) >= (g.Width*8) || (y*2)-HEADER_OFFSET >= (g.Height*8) {
		return nil
	}

	xPos := x*2 / TILE_WIDTH
	yPos := ((y*2)-HEADER_OFFSET) / TILE_HEIGHT

	return g.Tiles[xPos + yPos*g.Width]
}

func (g *Grid) PlaceBombs(firstTile *Tile) {
	if g.State != GAME_STARTED {
		return
	}

	var cnt int

	for cnt < g.TotalBombs {
		n := rand.Intn(len(g.Tiles))
		if g.Tiles[n].HaveBomb || g.Tiles[n] == firstTile {
			continue
		}

		g.Tiles[n].HaveBomb = true
		cnt += 1
	}

	for cnt, _ := range g.Tiles {
		g.FindSurroundingBombs(cnt)
	}

	g.State = GAME_RUNNING
	allSprites.TriggerEvent("StartTimer")
}

func (g *Grid) FindSurroundingBombs(pos int) {
	r := pos / g.Width
	c := pos % g.Width

	for _, rowCnt := range []int{-1, 0, 1} {
		for _, colCnt := range []int{-1, 0, 1} {
			t := g.GetTileAtPos(r+rowCnt, c+colCnt)
			if t != nil && t.HaveBomb {
				g.Tiles[pos].BombCount += 1
			}
		}
	}
}

func (g *Grid) GetTilePos(t *Tile) int {
	for cnt := 0; cnt < len(g.Tiles); cnt++ {
		if g.Tiles[cnt] == t {
			return cnt
		}
	}
	return -1
}

func (g *Grid) GetTileAtPos(r, c int) *Tile {
	if r < 0 || c < 0 {
		return nil
	}

	if r >= g.Height || c >= g.Width {
		return nil
	}

	return g.Tiles[r * g.Width + c]
}

func (g *Grid) RevealTileAtPos(pos int) {
	if pos == -1 || !g.Tiles[pos].Covered || g.State == GAME_OVER {
		return
	}

	g.Tiles[pos].RevealTile()
	if g.Tiles[pos].BombCount > 0 {
		return
	}

	g.RevealTileAtPos(g.Left(pos))
	g.RevealTileAtPos(g.Up(pos))
	g.RevealTileAtPos(g.Right(pos))
	g.RevealTileAtPos(g.Down(pos))

	g.RevealTileAtPos(g.UpLeft(pos))
	g.RevealTileAtPos(g.UpRight(pos))
	g.RevealTileAtPos(g.DownLeft(pos))
	g.RevealTileAtPos(g.DownRight(pos))

}

func (g *Grid) Up(pos int) int {
	if pos - g.Width > -1 {
		return pos - g.Width
	}
	return -1
}

func (g *Grid) Down(pos int) int {
	if pos + g.Width >= len(g.Tiles) {
		return -1
	}
	return pos + g.Width
}

func (g *Grid) Left(pos int) int {
	if pos % g.Width == 0 {
		return -1
	}
	return pos - 1
}

func (g *Grid) Right(pos int) int {
	if ((pos + 1) % g.Width) == 0 {
		return -1
	}
	return pos + 1
}

func (g *Grid) UpLeft(pos int) int {
	if g.Up(pos) == -1 || g.Left(pos) == -1 {
		return -1
	}
	return pos - g.Width - 1
}

func (g *Grid) UpRight(pos int) int {
	if g.Up(pos) == -1 || g.Right(pos) == -1 {
		return -1
	}
	return pos - g.Width + 1
}

func (g *Grid) DownLeft(pos int) int {
	if g.Down(pos) == -1 || g.Left(pos) == -1 {
		return -1
	}
	return pos + g.Width - 1
}

func (g *Grid) DownRight(pos int) int {
	if g.Down(pos) == -1 || g.Right(pos) == -1 {
		return -1
	}
	return pos + g.Width + 1
}

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
}

func main() {
	err := tm.Init()
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	w, h := tm.Size()
	Width = w*2
	Height = h*2

	setPalette()

	allSprites.Init(Width, Height, true)
	allSprites.Background = tm.Color187

	rand.Seed(time.Now().UnixNano())

	gameGrid = NewGrid()
	//t := NewTitle()
	//allSprites.Sprites = append(allSprites.Sprites, t)

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
			case <- done:
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
				if ev.Key == tm.KeyArrowUp {
				} else if ev.Key == tm.KeyArrowDown {
				}
			} else if ev.Type == tm.EventMouse {
				if ev.Key == tm.MouseLeft {
					MouseX = ev.MouseX
					MouseY = ev.MouseY
					t := gameGrid.FindTileClicked(ev.MouseX, ev.MouseY)
					if t != nil {
						pos := gameGrid.GetTilePos(t)
						gameGrid.PlaceBombs(t)
						gameGrid.RevealTileAtPos(pos)
						gameGrid.CheckGameOver()
					}
				} else if ev.Key == tm.MouseRight {
					if gameGrid.State != GAME_RUNNING {
						continue
					}
					MouseX = ev.MouseX
					MouseY = ev.MouseY
					t := gameGrid.FindTileClicked(ev.MouseX, ev.MouseY)
					if t != nil && t.Covered {
						t.SetFlag()
						gameGrid.CheckGameOver()
					}
				}

			} else if ev.Type == tm.EventResize {
				Width = ev.Width*2
				Height = ev.Height*2
				allSprites.Init(Width, Height, true)
				allSprites.Background = tm.Color187
				allSprites.TriggerEvent("resizeScreen")

				totalBombs := Width/8 * (Height-HEADER_OFFSET)/8 / 10
				gameGrid.SetSize(Width/8, (Height-HEADER_OFFSET)/8, totalBombs)
			}
		default:
			allSprites.Update()
			allSprites.Render()
			time.Sleep(60 * time.Millisecond)
		}
	}
}

