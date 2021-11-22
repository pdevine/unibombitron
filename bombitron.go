package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	sprite "github.com/pdevine/go-asciisprite"
)

const (
	TILE_WIDTH               = 8
	TILE_HEIGHT              = 8
	HEADER_OFFSET            = 10
	EASY_BOMB_RATE   float64 = 0.12345
	MEDIUM_BOMB_RATE float64 = 0.15625
	HARD_BOMB_RATE   float64 = 0.20625
)

type TileType int

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
	TILE_QUESTION
)

type Tile struct {
	sprite.BaseSprite
	GridX        int
	GridY        int
	VX           int
	VY           int
	BombCount    int
	HaveBomb     bool
	HaveFlag     bool
	HaveQuestion bool
	Covered      bool
}

type Background struct {
	sprite.BaseSprite
}

type Grid struct {
	State          GameState
	BombRate       float64
	Width          int
	Height         int
	Tiles          []*Tile
	TotalBombs     int
	FlagsRemaining *FlagsRemainingText
	TimerElapsed   *TimerElapsedText
	Super          *SuperText
	Background     *Background
	Kaboom         *Kaboom
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

type SuperText struct {
	sprite.BaseSprite
	TargetY int
	VY      float64
}

type Kaboom struct {
	sprite.BaseSprite
	Timer      int
	TimeOut    int
	TimeOutVis int
}

func NewFlagsRemaining() *FlagsRemainingText {
	f := &FlagsRemainingText{BaseSprite: sprite.BaseSprite{
		X:       4,
		Y:       1,
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
		X:       20,
		Y:       1,
		Visible: false},
		font: sprite.NewPakuFont(),
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
		t.X = Width - surf.Width - 4
	}
}

func NewSuperText() *SuperText {
	s := &SuperText{BaseSprite: sprite.BaseSprite{
		Visible: false},
	}
	s.Init()

	surf := sprite.NewSurfaceFromPng("super.png", true)
	s.BlockCostumes = append(s.BlockCostumes, &surf)

	s.X = Width/2 - surf.Width/2
	s.Y = -s.Height
	s.TargetY = Height/2 - surf.Height/2

	s.RegisterEvent("resizeScreen", func() {
		s.X = Width/2 - surf.Width/2
		s.Y = -s.Height
		s.TargetY = Height/2 - surf.Height/2
	})

	s.RegisterEvent("GameWon", func() {
		allSprites.MoveToTop(s)
		s.Visible = true
	})

	return s
}

func (s *SuperText) Update() {
	if !s.Visible || s.Y == s.TargetY {
		return
	}

	s.VY = (float64(s.TargetY) - float64(s.Y)) * 0.3
	s.Y += int(math.Round(s.VY))
}

func NewKaboom() *Kaboom {
	k := &Kaboom{BaseSprite: sprite.BaseSprite{
		Visible: false},
		TimeOut:    10,
		TimeOutVis: 50,
	}
	k.Init()

	surf1 := sprite.NewSurfaceFromPng("ka.png", true)
	k.BlockCostumes = append(k.BlockCostumes, &surf1)
	surf2 := sprite.NewSurfaceFromPng("kaboom.png", true)
	k.BlockCostumes = append(k.BlockCostumes, &surf2)
	k.X = Width/2 - surf1.Width/2
	k.Y = Height/2 - surf1.Height/2

	k.RegisterEvent("Explode", func() {
		allSprites.MoveToTop(k)
		k.Visible = true
	})

	k.RegisterEvent("resizeScreen", func() {
		k.X = Width/2 - surf1.Width/2
		k.Y = Height/2 - surf1.Height/2
	})

	return k
}

func (k *Kaboom) Update() {
	if !k.Visible {
		return
	}
	k.Timer++

	if k.Timer > k.TimeOut {
		if k.CurrentCostume < len(k.BlockCostumes)-1 {
			k.NextCostume()
		}
	}
	if k.Timer > k.TimeOutVis {
		k.Visible = false
	}
}

func NewBackground() *Background {
	b := &Background{BaseSprite: sprite.BaseSprite{
		Visible: true},
	}
	b.Init()

	b.RegisterEvent("resizeScreen", func() {
		surf := sprite.NewSurface(Width, Height, true)
		x0 := gameGrid.Width * TILE_WIDTH
		y0 := HEADER_OFFSET
		x1 := gameGrid.Width * TILE_WIDTH
		y1 := gameGrid.Height*TILE_HEIGHT + HEADER_OFFSET
		surf.Line(x0, y0, x1, y1, 'X')
		surf.Line(0, y1, x1, y1, 'X')
		b.BlockCostumes = []*sprite.Surface{&surf}
	})

	return b
}

func NewTile() *Tile {
	t := &Tile{BaseSprite: sprite.BaseSprite{
		Visible: true},
		Covered: true,
	}
	t.Init()

	t.RegisterEvent("GameWon", func() {
		if t.HaveFlag == true {
			t.VX, t.VY = randVec()
		}
	})

	t.RegisterEvent("ReturnToGrid", func() {
		t.VX = 0
		t.VY = 0
		t.X = t.GridX
		t.Y = t.GridY
	})

	t.SetTile(TILE_COVERED)
	return t
}

func randVec() (int, int) {
	var x, y int
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(2)
	x = 1
	if n == 0 {
		x = -1
	}

	n = r.Intn(2)
	y = 1
	if n == 0 {
		y = -1
	}
	return x, y
}

func (t *Tile) RevealTile() {
	if t.HaveFlag || t.HaveQuestion {
		return
	} else if t.HaveBomb {
		t.SetTile(TILE_BOMB)
		allSprites.TriggerEvent("Explode")
		allSprites.TriggerEvent("GameOver")
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
		t.HaveQuestion = true
		t.SetTile(TILE_QUESTION)
		allSprites.TriggerEvent("ShowFlagsRemaining")
	} else if t.HaveQuestion {
		t.HaveQuestion = false
		t.SetTile(TILE_COVERED)
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
		tileQuestion,
	}
	surf := sprite.NewSurfaceFromString(tileImages[v], false)
	t.BlockCostumes = []*sprite.Surface{&surf}
	t.SetCostume(0)
}

func (t *Tile) Update() {
	if gameGrid.State != GAME_OVER {
		return
	}

	t.X += t.VX
	t.Y += t.VY

	if t.X < 0 {
		t.VX = 1
	} else if t.X > Width-t.Width {
		t.VX = -1
	}

	if t.Y < 0 {
		t.VY = 1
	} else if t.Y > Height-t.Height {
		t.VY = -1
	}
}

func NewGrid() *Grid {
	g := &Grid{
		State:          GAME_INIT,
		FlagsRemaining: NewFlagsRemaining(),
		TimerElapsed:   NewTimerElapsed(),
		Super:          NewSuperText(),
		Background:     NewBackground(),
		Kaboom:         NewKaboom(),
	}
	return g
}

func (g *Grid) SetReady() {
	if g.State != GAME_INIT {
		return
	}

	allSprites.Sprites = append(allSprites.Sprites, g.FlagsRemaining)
	allSprites.Sprites = append(allSprites.Sprites, g.TimerElapsed)
	allSprites.Sprites = append(allSprites.Sprites, g.Super)
	allSprites.Sprites = append(allSprites.Sprites, g.Kaboom)
	allSprites.Sprites = append(allSprites.Sprites, g.Background)
	g.State = GAME_READY
}

func (g *Grid) CheckGameOver() bool {
	for _, t := range g.Tiles {
		if t.Covered && !t.HaveFlag {
			return false
		}
	}
	g.State = GAME_OVER
	allSprites.TriggerEvent("GameWon")
	allSprites.TriggerEvent("GameOver")
	return true
}

func (g *Grid) SetSize(w, h int) {
	if g.State != GAME_READY {
		return
	}

	g.Width = w
	g.Height = h

	g.Tiles = make([]*Tile, 0, 0)

	// Add the tiles
	for cntY := 0; cntY < h; cntY++ {
		for cntX := 0; cntX < w; cntX++ {
			t := NewTile()
			t.X = cntX * 8
			t.Y = (cntY * 8) + HEADER_OFFSET
			t.GridX = t.X
			t.GridY = t.Y
			g.Tiles = append(g.Tiles, t)
			allSprites.Sprites = append(allSprites.Sprites, t)
		}
	}
}

func (g *Grid) FindTileClicked(x, y int) *Tile {
	if x >= (g.Width*8) || y-HEADER_OFFSET >= (g.Height*8) || y < HEADER_OFFSET {
		return nil
	}

	xPos := x / TILE_WIDTH
	yPos := (y - HEADER_OFFSET) / TILE_HEIGHT

	return g.Tiles[xPos+yPos*g.Width]
}

func (g *Grid) PlaceBombs(firstTile *Tile) {
	if g.State != GAME_STARTED {
		return
	}

	g.FlagsRemaining.Remaining = g.TotalBombs
	allSprites.TriggerEvent("ShowFlagsRemaining")

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

	return g.Tiles[r*g.Width+c]
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
	if pos-g.Width > -1 {
		return pos - g.Width
	}
	return -1
}

func (g *Grid) Down(pos int) int {
	if pos+g.Width >= len(g.Tiles) {
		return -1
	}
	return pos + g.Width
}

func (g *Grid) Left(pos int) int {
	if pos%g.Width == 0 {
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
