package game

import (
	"fmt"
	"image/color"
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
	for range time.NewTicker(time.Millisecond * 700).C {
		if !e.IsDead {
			e.AttackAction = !e.AttackAction
		} else {
			break
		}
	}
}

func NewGame() (*Game, error) {
	frames := make(map[string]*Unit)
	enemies := make(map[string]*Enemy)
	var hk *HeroKnight
	if unit, err := GetFramesHeroKnight(); err != nil {
		return nil, err
	} else {
		frames[MainCharacter] = unit
		hk = NewHeroKnight(30.0, 0.0, unit.Width, unit.Height)
	}
	if unit, err := GetFramesBandit(LightBandit); err != nil {
		return nil, err
	} else {
		frames[LightBandit] = unit
		enemies[LightBandit] = NewEnemy(LightBandit, 400.0, 0.0, unit.Width, unit.Height)
	}
	if unit, err := GetFramesBandit(HeavyBandit); err != nil {
		return nil, err
	} else {
		frames[HeavyBandit] = unit
		enemies[HeavyBandit] = NewEnemy(HeavyBandit, 350.0, 0.0, unit.Width, unit.Height)
	}
	if unit, err := GetFramesWizard(); err != nil {
		return nil, err
	} else {
		frames[Wizard] = unit
		enemies[Wizard] = NewEnemy(Wizard, 450.0, 0.0, unit.Width, unit.Height)
	}

	game := &Game{
		Camera: &Camera{
			EnemyY: -(360 * float64(Scale)) / 2,
		},
		Frames:  frames,
		hk:      hk,
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
			enemy.RunLeftAction = (enemy.X+enemy.Indent)-(g.hk.X+g.hk.Width-g.hk.Indent) > 2
			enemy.RunRightAction = (g.hk.X+g.hk.Indent)-(enemy.X+enemy.Width-enemy.Indent) > 2
		}

		go func(enemy *Enemy) {
			enemy.Once.Do(func() {
				onceBody(enemy)
			})
		}(enemy)
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
	ebitenutil.DrawRect(screen, 9.0, 49.0, 100.0+2.0, 20.0+2.0, color.RGBA{255, 255, 255, 255})
	ebitenutil.DrawRect(screen, 10.0, 50.0, float64(g.hk.Health), 20.0, color.RGBA{255, 0, 0, 255})
}
