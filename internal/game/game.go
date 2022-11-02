package game

import (
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	hk *HeroKnight
	lb *LightBandit
	hb *LightBandit
}

var boxesShow bool
var once sync.Once
var attack bool
var runLeft bool
var runRight bool

func NewGame() *Game {
	return &Game{
		hk: NewHeroKnight(),
		lb: NewLightBandit(),
		hb: NewHeavyBandit(),
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		boxesShow = !boxesShow
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyZ) {
		g.lb.Health = 100
		g.lb.IsDead = false
	}
	var onceBody = func() {
		ticker := time.NewTicker(time.Millisecond * 700)
	LOOP:
		for {
			select {
			case <-ticker.C:
				if !g.lb.IsDead {
					attack = !attack
				} else {
					break LOOP
				}
			}
		}
	}
	go func() {
		once.Do(onceBody)
	}()

	if !g.lb.IsDead {
		runLeft = (g.lb.X+12)-(g.hk.X+65) > 2
		runRight = (g.hk.X+35)-(g.lb.X+36) > 2
	}

	if !g.lb.IsDead {
		g.lb.Update(attack, runLeft, runRight)
	}

	if !g.hb.IsDead {
		g.hb.Update(attack, runLeft, runRight)
	}

	if !g.hk.IsDead || g.hk.Frame != g.hk.LastFrame {
		g.hk.Update()
	}

	if (g.hk.Status == StatusAttack1 || g.hk.Status == StatusAttack2 || g.hk.Status == StatusAttack3) && g.hk.Frame == g.hk.LastFrame/2 && g.hk.Y < 48 {
		if ((g.lb.X+12)-(g.hk.X+65) < 35 &&
			((g.lb.X+12)-(g.hk.X+65)) > -1 &&
			g.hk.Side > 0) ||
			((g.hk.X+35)-(g.lb.X+36) < 35 &&
				(g.hk.X+35)-(g.lb.X+12) > -1 &&
				g.hk.Side < 0) {
			g.lb.Hurt()
		}
	}

	if g.lb.Status == StatusAttack && g.lb.Frame == g.lb.LastFrame/2 && !g.hk.IsRolling && g.hk.Y < 48 {
		if ((g.lb.X+12)-(g.hk.X+65) < 12 &&
			(g.lb.X+12)-(g.hk.X+65) > -1 &&
			g.lb.Side > 0) ||
			((g.hk.X+35)-(g.lb.X+36) < 12 &&
				(g.hk.X+35)-(g.lb.X+36) > -1 &&
				g.lb.Side < 0) {
			g.hk.Hurt()
		}
	}

	if g.hk.IsRunning && !g.hk.IsRolling && !g.lb.IsDead {
		if g.hk.Side > 0 {
			if g.hk.X+65 > g.lb.X+12 && g.hk.X+35 < g.lb.X+36 && g.hk.Y < 48 {
				g.hk.X = g.lb.X + 12 - 65
			}
		} else {
			if g.hk.X+35 < g.lb.X+36 && g.hk.X+65 > g.lb.X+12 && g.hk.Y < 48 {
				g.hk.X = g.lb.X + 36 - 35
			}
		}
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.lb.Draw(screen)
	g.hb.Draw(screen)
	g.hk.Draw(screen)
}
