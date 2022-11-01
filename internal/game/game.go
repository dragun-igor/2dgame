package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	hk *HeroKnight
	lb *LightBandit
}

var boxesShow bool

func NewGame() *Game {
	return &Game{
		hk: NewHeroKnight(),
		lb: NewLightBandit(),
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		boxesShow = !boxesShow
	}
	if !g.lb.IsDead {
		g.lb.Update()
	}
	if !g.hk.IsDead || g.hk.Frame != g.hk.LastFrame {
		g.hk.Update()
	}

	if g.hk.IsAttacking && g.hk.Frame == g.hk.LastFrame/2 {
		if (g.hk.X+100 > g.lb.X+20 &&
			g.hk.X+100 < g.lb.X+70 &&
			g.hk.Side > 0) ||
			(g.hk.X > g.lb.X-70 &&
				g.hk.X < g.lb.X+32 &&
				g.hk.Side < 0) {
			g.lb.Hurt()
		}
	}
	if g.hk.IsRunning && !g.lb.IsDead {
		if g.hk.Side > 0 {
			if g.hk.X+63 > g.lb.X+12 && g.hk.X+40 < g.lb.X+38 {
				g.hk.X = g.lb.X + 12 - 63
			}
		} else {
			if g.hk.X+40 < g.lb.X+38 && g.hk.X+63 > g.lb.X+12 {
				g.hk.X = g.lb.X + 38 - 40
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
	g.hk.Draw(screen)
}
