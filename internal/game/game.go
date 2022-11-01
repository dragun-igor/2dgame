package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	hk *HeroKnight
	lb *LightBandit
}

func NewGame() *Game {
	return &Game{
		hk: NewHeroKnight(),
		lb: NewLightBandit(),
	}
}

func (g *Game) Update() error {
	g.lb.Update()
	g.hk.Update()
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
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.hk.Draw(screen)
	g.lb.Draw(screen)
}
