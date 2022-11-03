package game

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	hk      *HeroKnight
	enemies map[string]*Bandit
}

var boxesShow bool
var onceBody = func(b *Bandit) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	ticker := time.NewTicker(time.Millisecond * 700)
LOOP:
	for {
		select {
		case <-ticker.C:
			if !b.IsDead {
				b.AttackAction = !b.AttackAction
			} else {
				break LOOP
			}
		}
	}
}

func NewGame() *Game {
	enemies := make(map[string]*Bandit)
	enemies["bandit1"] = NewHeavyBandit()
	enemies["bandit2"] = NewLightBandit()
	return &Game{
		hk:      NewHeroKnight(),
		enemies: enemies,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		boxesShow = !boxesShow
	}
	rand.Seed(time.Now().UnixMicro())

	for _, unit := range g.enemies {
		if !unit.IsDead {
			unit.RunLeftAction = (unit.X+float64(12*Scale))-(g.hk.X+float64(65*Scale)) > 2
			unit.RunRightAction = (g.hk.X+float64(35*Scale))-(unit.X+float64(36*Scale)) > 2
		}

		go func(unit *Bandit) {
			unit.Once.Do(func() {
				onceBody(unit)
			})
		}(unit)

		if !unit.IsDead {
			unit.Update(g.hk)
		}
	}

	if !g.hk.IsDead || g.hk.Frame != g.hk.LastFrame {
		g.hk.Update(g.enemies)
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	var keys []string
	for key := range g.enemies {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, key := range keys {
		g.enemies[key].Draw(screen, g.hk.X, g.hk.Y)
	}
	g.hk.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTPS: %.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}
