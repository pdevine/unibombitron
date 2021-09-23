package main

import (
	"math"

	sprite "github.com/pdevine/go-asciisprite"
)

type TitleOverlay struct {
	Selectors []*Selector
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

func NewTitleOverlay() *TitleOverlay {
	t := &TitleOverlay{
		Selectors: []*Selector{
			NewSelector("easy"),
			NewSelector("med."),
			NewSelector("hard"),
		},
	}

	for _, s := range t.Selectors {
		allSprites.Sprites = append(allSprites.Sprites, s)
	}

	return t
}

func (t *TitleOverlay) MoveToTop() {
	for _, s := range t.Selectors {
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

	surf := sprite.NewSurface(40, 10, false)
	for rcnt, r := range surf.Blocks {
		for ccnt, _ := range r {
			surf.Blocks[rcnt][ccnt] = 'w'
		}
	}
	surf.Rectangle(1, 1, 39, 9, 'X')
	surf.Blit(w, surf.Width/2-w.Width/2, 2)
	s.BlockCostumes = []*sprite.Surface{&surf}
	s.SetCostume(0)

	s.RegisterEvent("resizeScreen", func() {
		if n == "easy" {
			s.TargetX = 10
			s.X = -surf.Width
			s.Y = Height - 20
			s.BombRate = EASY_BOMB_RATE
		} else if n == "med." {
			s.X = Width/2 - surf.Width/2
			s.Y = Height + 10
			s.TargetY = Height - 21
			s.BombRate = MEDIUM_BOMB_RATE
		} else if n == "hard" {
			s.TargetX = Width - surf.Width - 10
			s.X = Width
			s.Y = Height - 20
			s.BombRate = HARD_BOMB_RATE
		}
	})

	s.RegisterEvent("SelectorClicked", func() {
		s.Visible = false
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
