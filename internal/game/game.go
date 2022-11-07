package game

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Frames  map[string]*Unit
	enemies map[string]*Bandit
	hk      *HeroKnight
	w       *Wizard
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

var onceBody1 = func(b *Wizard) {
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
	frames := make(map[string]*Unit)
	if unit, err := GetFramesHeroKnight(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println()
		fmt.Println(unit)
		frames["HeroKnight"] = unit
	}
	if unit, err := GetFramesBandit(LightBandit); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println()
		fmt.Println(unit)
		frames["LightBandit"] = unit
	}
	if unit, err := GetFramesBandit(HeavyBandit); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println()
		fmt.Println(unit)
		frames["HeavyBandit"] = unit
	}
	if unit, err := GetFramesWizard(); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println()
		fmt.Println(unit)
		frames["Wizard"] = unit
	}

	enemies := make(map[string]*Bandit)
	enemies[HeavyBandit] = NewHeavyBandit()
	enemies[LightBandit] = NewLightBandit()
	wizard := NewWizard()
	return &Game{
		Frames:  frames,
		hk:      NewHeroKnight(),
		enemies: enemies,
		w:       wizard,
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

	if !g.w.IsDead {
		g.w.RunLeftAction = (g.w.X+float64(12*Scale))-(g.hk.X+float64(65*Scale)) > 2
		g.w.RunRightAction = (g.hk.X+float64(35*Scale))-(g.w.X+float64(36*Scale)) > 2
	}

	go func(unit *Wizard) {
		unit.Once.Do(func() {
			onceBody1(unit)
		})
	}(g.w)

	if !g.w.IsDead {
		g.w.Update(g.hk)
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
		g.enemies[key].Draw(screen, g.Frames[key], g.hk.X, g.hk.Y)
	}
	// g.w.Draw(screen, g.Frames["Wizard"], g.hk.X, g.hk.Y)
	g.hk.Draw(screen, g.Frames["HeroKnight"])
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTPS: %.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}
