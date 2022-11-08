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
	Camera        *Camera
	Frames        map[string]*Unit
	Enemies       map[string]*Enemy
	MainCharacter *HeroKnight
	Environment   []Environment
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

	environment := make([]Environment, 0, 200)
	// for i := 0; i < 50; i++ {
	// 	op := &ebiten.DrawImageOptions{}
	// 	op.GeoM.Scale(1.0*float64(Scale), 1.0*float64(Scale))
	// 	op.GeoM.Translate(float64(TileSize*i*Scale)-camera.X, hk.Y-camera.MainCharacterY)
	// 	screen.DrawImage(unit.ActionFrames["environment"][1], op)
	// }
	for i := 0; i < 50; i++ {
		environment = append(environment, Environment{
			X:      float64(TileSize * i * Scale),
			Y:      0.0,
			Width:  float64(TileSize * Scale),
			Height: float64(TileSize * Scale),
			Scale:  1.0,
		})
	}

	for i := 15; i < 18; i++ {
		environment = append(environment, Environment{
			X:      float64(TileSize * i * Scale),
			Y:      -float64(TileSize * Scale),
			Width:  float64(TileSize * Scale),
			Height: float64(TileSize * Scale),
			Scale:  1.0,
		})
	}

	game := &Game{
		Camera: &Camera{
			EnemyY: -(360 * float64(Scale)) / 2,
		},
		Frames:        frames,
		MainCharacter: hk,
		Enemies:       enemies,
		Environment:   environment,
	}

	return game, nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		boxesShow = !boxesShow
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		Scale++
		if Scale > 3 {
			Scale = 1
		}
	}
	rand.Seed(time.Now().UnixMicro())

	for _, enemy := range g.Enemies {
		if !enemy.IsDead {
			enemy.RunLeftAction = (enemy.X+enemy.Indent)-(g.MainCharacter.X+g.MainCharacter.Width-g.MainCharacter.Indent) > 2
			enemy.RunRightAction = (g.MainCharacter.X+g.MainCharacter.Indent)-(enemy.X+enemy.Width-enemy.Indent) > 2
		}

		go func(enemy *Enemy) {
			enemy.Once.Do(func() {
				onceBody(enemy)
			})
		}(enemy)
		if !enemy.IsDead || enemy.Frame != enemy.LastFrame {
			enemy.Update(g.MainCharacter)
		}
	}

	if !g.MainCharacter.IsDead || g.MainCharacter.Frame != g.MainCharacter.LastFrame {
		g.MainCharacter.Update(g.Enemies, g.Environment)
	}

	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Camera.RefreshOffset(g.MainCharacter.X, g.MainCharacter.Y)
	for _, env := range g.Environment {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(Scale), float64(Scale))
		op.GeoM.Translate(env.X-g.Camera.X, env.Y-g.Camera.EnemyY)
		screen.DrawImage(g.Frames[MainCharacter].ActionFrames["environment"][1], op)
	}
	keys := make([]string, 0, len(g.Enemies))
	for key := range g.Enemies {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	for _, key := range keys {
		g.Enemies[key].Draw(screen, g.Frames[key], g.Camera)
	}
	g.MainCharacter.Draw(screen, g.Frames[MainCharacter], g.Camera)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTPS: %.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
	ebitenutil.DrawRect(screen, 9.0, 49.0, 100.0+2.0, 20.0+2.0, color.RGBA{255, 255, 255, 255})
	ebitenutil.DrawRect(screen, 10.0, 50.0, float64(g.MainCharacter.Health), 20.0, color.RGBA{255, 0, 0, 255})
}
