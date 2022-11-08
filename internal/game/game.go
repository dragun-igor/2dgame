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
	Camera  *Camera
	Frames  map[string]*Unit
	enemies map[string]*Enemy
	hk      *HeroKnight
}

type Camera struct {
	X              float64
	MainCharacterY float64
	EnemyY         float64
}

func (c *Camera) RefreshOffset(x, y float64) {
	c.X = x - float64(640*Scale-300)/2

	if c.X < 0 {
		c.X = 0
	}
	if c.X > 355 {
		c.X = 355
	}
	c.MainCharacterY = y - float64(360*Scale)/2
}

var boxesShow bool
var onceBody = func(e *Enemy) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	ticker := time.NewTicker(time.Millisecond * 700)
LOOP:
	for {
		select {
		case <-ticker.C:
			if !e.IsDead {
				e.AttackAction = !e.AttackAction
			} else {
				break LOOP
			}
		}
	}
}

func NewGame() (*Game, error) {
	frames := make(map[string]*Unit)
	if unit, err := GetFramesHeroKnight(); err != nil {
		return nil, err
	} else {
		frames[MainCharacter] = unit
	}
	if unit, err := GetFramesBandit(LightBandit); err != nil {
		return nil, err
	} else {
		frames[LightBandit] = unit
	}
	if unit, err := GetFramesBandit(HeavyBandit); err != nil {
		return nil, err
	} else {
		frames[HeavyBandit] = unit
	}
	if unit, err := GetFramesWizard(); err != nil {
		return nil, err
	} else {
		frames[Wizard] = unit
	}

	enemies := make(map[string]*Enemy)
	enemies[HeavyBandit] = NewHeavyBandit()
	enemies[LightBandit] = NewLightBandit()
	enemies[Wizard] = NewWizard()
	game := &Game{
		Camera: &Camera{
			EnemyY: -(360 * float64(Scale)) / 2,
		},
		Frames:  frames,
		hk:      NewHeroKnight(),
		enemies: enemies,
	}

	return game, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		boxesShow = !boxesShow
	}
	rand.Seed(time.Now().UnixMicro())

	for _, enemy := range g.enemies {
		if !enemy.IsDead {
			enemy.RunLeftAction = (enemy.X+float64(12*Scale))-(g.hk.X+float64(65*Scale)) > 2
			enemy.RunRightAction = (g.hk.X+float64(35*Scale))-(enemy.X+float64(36*Scale)) > 2
		}

		go func(enemy *Enemy) {
			enemy.Once.Do(func() {
				onceBody(enemy)
			})
		}(enemy)
		if enemy.Type == Wizard {
			fmt.Println(enemy.Frame)
			fmt.Println(enemy.LastFrame)
		}
		if !enemy.IsDead || enemy.Frame != enemy.LastFrame {
			enemy.Update(g.hk)
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
	g.Camera.RefreshOffset(g.hk.X, g.hk.Y)
	keys := make([]string, 0, len(g.enemies))
	for key := range g.enemies {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, key := range keys {
		g.enemies[key].Draw(screen, g.Frames[key], g.Camera)
	}
	g.hk.Draw(screen, g.Frames[MainCharacter], g.Camera)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTPS: %.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}
