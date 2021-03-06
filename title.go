package main

import (
	"image"
	"math"
	"math/rand"

	sprite "github.com/pdevine/go-asciisprite"
)

type TitleOverlay struct {
	Selectors []*Selector
	Logo      *TitleLogo
	Bomb      *TitleBomb
	Uni       *UniLogo
}

type TitleLogo struct {
	sprite.BaseSprite
}

type UniLogo struct {
	sprite.BaseSprite
	Timer   int
	TimeOut int
}

type AdjustText struct {
	sprite.BaseSprite
}

type Selector struct {
	sprite.BaseSprite
	Type     string
	TargetX  int
	TargetY  int
	VX       float64
	VY       float64
	BombRate float64
}

type Spark struct {
	sprite.BaseSprite
	Yoffset  int
	Dead     bool
	Lifetime int
	VX       int
	VY       int
}

type TitleBomb struct {
	sprite.BaseSprite
	TargetY int
	VY      float64
	Sparks  []*Spark
}

func NewTitleOverlay() *TitleOverlay {
	t := &TitleOverlay{}

	adjText := NewAdjustText()
	allSprites.Sprites = append(allSprites.Sprites, adjText)

	return t
}

func (t *TitleOverlay) SetGameReady() {
	t.Selectors = []*Selector{
		NewSelector("easy"),
		NewSelector("med."),
		NewSelector("hard"),
	}
	t.Logo = NewTitleLogo()
	t.Bomb = NewTitleBomb()
	t.Uni = NewUniLogo()

	allSprites.Sprites = append(allSprites.Sprites, t.Logo)
	allSprites.Sprites = append(allSprites.Sprites, t.Bomb)
	allSprites.Sprites = append(allSprites.Sprites, t.Uni)

	for _, s := range t.Selectors {
		allSprites.Sprites = append(allSprites.Sprites, s)
	}
}

func (t *TitleOverlay) MoveToTop() {
	allSprites.MoveToTop(t.Logo)
	allSprites.MoveToTop(t.Uni)
	allSprites.MoveToTop(t.Bomb)
	for _, s := range t.Selectors {
		allSprites.MoveToTop(s)
	}
	for _, s := range t.Bomb.Sparks {
		allSprites.MoveToTop(s)
	}
}

func (t *TitleOverlay) CheckSelectorClicked(x, y int) *Selector {
	for _, s := range t.Selectors {
		if s.HitAtPointSurface(x, y) {
			return s
		}
	}
	return nil
}

func NewSelector(n string) *Selector {
	s := &Selector{BaseSprite: sprite.BaseSprite{
		Y:       Height - 20,
		Visible: true},
		Type: n,
	}
	s.Init()

	f := sprite.NewPakuFont()
	w := sprite.NewSurfaceFromString(f.BuildString(n), false)

	surf1 := sprite.NewSurface(40, 10, false)
	surf2 := sprite.NewSurface(40, 10, false)
	for rcnt, r := range surf1.Blocks {
		for ccnt, _ := range r {
			surf1.Blocks[rcnt][ccnt] = 'l'
			surf2.Blocks[rcnt][ccnt] = 'w'
		}
	}
	surf1.Rectangle(0, 0, 39, 9, 'X')
	surf2.Rectangle(0, 0, 39, 9, 'X')
	surf1.Blit(w, surf1.Width/2-w.Width/2, 2)
	surf2.Blit(w, surf2.Width/2-w.Width/2, 2)
	s.BlockCostumes = []*sprite.Surface{&surf1, &surf2}
	s.SetCostume(0)

	if n == "easy" {
		s.TargetX = 10
		s.X = -surf1.Width
		s.Y = Height - 20
		s.BombRate = EASY_BOMB_RATE
	} else if n == "med." {
		s.X = Width/2 - surf1.Width/2
		s.Y = Height + 10
		s.TargetY = Height - 21
		s.BombRate = MEDIUM_BOMB_RATE
	} else if n == "hard" {
		s.TargetX = Width - surf1.Width - 10
		s.X = Width
		s.Y = Height - 20
		s.BombRate = HARD_BOMB_RATE
	}

	s.RegisterEvent("SelectorClicked", func() {
		s.Visible = false
	})

	s.RegisterEvent("MouseMove", func() {
		if MouseX >= s.X && MouseX < s.X+surf1.Width && MouseY >= s.Y && MouseY < s.Y+surf1.Height {
			s.SetCostume(1)
		} else {
			s.SetCostume(0)
		}
	})

	return s
}

func (s *Selector) Update() {
	if !s.Visible {
		return
	}

	if s.Type == "easy" || s.Type == "hard" {
		if s.TargetX == s.X {
			return
		}
		s.VX = (float64(s.TargetX) - float64(s.X)) * 0.3
		s.X += int(math.Round(s.VX))
	} else {
		if s.TargetY == s.Y {
			return
		}
		s.VY = (float64(s.TargetY) - float64(s.Y)) * 0.3
		s.Y += int(math.Round(s.VY))
	}

}

func NewTitleLogo() *TitleLogo {
	t := &TitleLogo{BaseSprite: sprite.BaseSprite{
		Visible: true},
	}
	t.Init()

	surf := sprite.NewSurfaceFromPng("title.png", true)

	f := sprite.NewJRSMFont()
	cSurf := sprite.NewSurfaceFromString(f.BuildString("(c) 2021 Patrick Devine"), true)
	surf.Blit(cSurf, 22, 25)

	t.BlockCostumes = append(t.BlockCostumes, &surf)

	t.X = Width/2 - surf.Width/2
	t.Y = 16

	t.RegisterEvent("SelectorClicked", func() {
		t.Visible = false
	})

	return t
}

func NewUniLogo() *UniLogo {
	u := &UniLogo{BaseSprite: sprite.BaseSprite{
		Y:       10,
		Visible: true},
		TimeOut: 2,
	}
	u.X = Width/2 - 84
	u.Init()

	r := image.Rect(0, 0, 28, 15)
	surfs := sprite.NewSurfacesFromPngSheet("uni.png", r, true)

	for cnt := 0; cnt < len(surfs); cnt++ {
		u.BlockCostumes = append(u.BlockCostumes, &surfs[cnt])
	}
	u.SetCostume(0)

	u.RegisterEvent("SelectorClicked", func() {
		u.Visible = false
	})
	return u
}

func (u *UniLogo) Update() {
	if u.CurrentCostume == len(u.BlockCostumes)-1 {
		return
	}
	if u.Timer >= u.TimeOut {
		u.NextCostume()
		u.Timer = 0
	}
	u.Timer++
}

func NewSpark() *Spark {
	s := &Spark{BaseSprite: sprite.BaseSprite{
		Visible: true},
	}
	s.Init()
	s.Reset()

	colors := []string{"o", "y", "r"}
	c := rand.Intn(len(colors))

	surf := sprite.NewSurfaceFromString(colors[c], false)
	s.BlockCostumes = []*sprite.Surface{&surf}
	s.SetCostume(0)

	s.RegisterEvent("SelectorClicked", func() {
		s.Visible = false
	})

	return s
}

func (s *Spark) Update() {
	s.Lifetime -= 1
	if s.Lifetime <= 0 {
		s.Reset()
	}
	s.X += s.VX
	s.Y += s.VY
}

func (s *Spark) Reset() {
	s.X = Width/2 - 35
	s.Y = s.Yoffset - 1
	s.VX = rand.Intn(4) - 2
	s.VY = rand.Intn(4) - 3
	s.Lifetime = rand.Intn(5) + 2
}

func NewTitleBomb() *TitleBomb {
	b := &TitleBomb{BaseSprite: sprite.BaseSprite{
		Y:       -30,
		Visible: true},
		TargetY: 19,
		Sparks:  []*Spark{},
	}
	b.Init()

	surf := sprite.NewSurfaceFromPng("bomb.png", true)
	b.BlockCostumes = append(b.BlockCostumes, &surf)
	b.X = Width/2 - surf.Width/2 - 44

	for cnt := 0; cnt < 15; cnt++ {
		s := NewSpark()
		b.Sparks = append(b.Sparks, s)
		allSprites.Sprites = append(allSprites.Sprites, s)
	}

	b.RegisterEvent("SelectorClicked", func() {
		b.Visible = false
	})
	return b
}

func (b *TitleBomb) Update() {
	if b.TargetY == b.Y {
		return
	}
	b.VY = (float64(b.TargetY) - float64(b.Y)) * 0.3
	b.Y += int(math.Round(b.VY))

	for _, s := range b.Sparks {
		s.Yoffset = b.Y
	}
}

func NewAdjustText() *AdjustText {
	a := &AdjustText{BaseSprite: sprite.BaseSprite{
		Visible: false},
	}
	a.Init()

	f := sprite.NewPakuFont()
	surf := sprite.NewSurfaceFromString(f.BuildString("your terminal, too small"), false)
	a.BlockCostumes = append(a.BlockCostumes, &surf)
	a.SetCostume(0)

	a.RegisterEvent("resizeScreen", func() {
		a.X = Width/2 - surf.Width/2
		a.Y = Height/2 - surf.Height/2
		if Width < 40 || Height < 20 {
			a.Visible = true
		} else {
			a.Visible = false
		}
	})
	return a
}
