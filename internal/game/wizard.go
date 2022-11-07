package game

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Wizard struct {
	Keyboard         Keyboard
	Status           string
	PrevStatus       string
	Type             string
	X                float64
	Y                float64
	Side             float64
	SpeedRun         float64
	MaxSpeedRun      float64
	BaseSpeedRun     float64
	AccelerationRun  float64
	SpeedJump        float64
	BaseSpeedJump    float64
	DecelerationJump float64
	Frame            int
	LastFrame        int
	Health           int
	Direction        uint8
	IsDead           bool
	IsHurted         bool
	IsAttacking      bool
	IsJumping        bool
	IsRunning        bool
	AttackAction     bool
	RunLeftAction    bool
	RunRightAction   bool
	Once             sync.Once
}

func NewWizard() *Wizard {
	return &Wizard{
		Keyboard:         NewEmulationKeyboard(),
		Type:             "LightBandit",
		X:                550,
		Y:                0,
		Status:           StatusIdle,
		PrevStatus:       StatusIdle,
		Side:             1.0,
		SpeedRun:         1.0,
		MaxSpeedRun:      2.0,
		BaseSpeedRun:     1.0,
		AccelerationRun:  0.2,
		SpeedJump:        6.0,
		BaseSpeedJump:    6.0,
		DecelerationJump: 0.2,
		Health:           100,
		LastFrame:        StatusFramesWizard[StatusCombatIdle].FramesNumber*StatusFramesWizard[StatusCombatIdle].FrameDuration - 1,
		Direction:        DirectionLeft,
	}
}

func GetFramesWizard() (*Unit, error) {
	var file *os.File
	var img image.Image
	var err error
	actionFrames := make(map[string][]*ebiten.Image)

	for status := range StatusFramesWizard {
		framesNumber := StatusFramesWizard[status].FramesNumber
		frms := make([]*ebiten.Image, 0, framesNumber)
		file, err = os.Open("_assets/Wizard/" + status + ".png")
		if err != nil {
			break
		}
		img, err = png.Decode(file)
		if err != nil {
			break
		}
		file.Close()
		ebitenImg := ebiten.NewImageFromImage(img)
		for i := 0; i < framesNumber; i++ {
			frms = append(frms, ebitenImg.SubImage(image.Rect(231*i, 0, 231*i+231, 150)).(*ebiten.Image))
		}
		actionFrames[status] = frms
	}
	unit := &Unit{
		ActionFrames: actionFrames,
		Width:        231,
		Height:       140,
	}
	return unit, err
}

func (w *Wizard) Death() {
	if w.Health <= 0 {
		w.IsDead = true
	}
}

func (w *Wizard) Hurt() {
	if !w.IsHurted {
		w.IsHurted = true
		w.Health -= 30
	}
}

func (w *Wizard) Attack() {
	if !w.IsAttacking && !w.IsDead {
		w.IsAttacking = true
	}
}

func (w *Wizard) Run() {
	w.IsRunning = true
}

func (w *Wizard) Stop() {
	w.IsRunning = false
}

func (w *Wizard) Update(heroKnight *HeroKnight) error {
	w.Keyboard[KeyAttack].KeyEmulation = w.AttackAction
	w.Keyboard[KeyRunLeft].KeyEmulation = w.RunLeftAction
	w.Keyboard[KeyRunRight].KeyEmulation = w.RunRightAction
	w.Keyboard.Update()
	w.Death()
	if w.Keyboard[KeyAttack].IsKeyJustPressed {
		w.Attack()
	}
	if w.Keyboard[KeyRunLeft].IsKeyPressed {
		w.Direction = DirectionLeft
	}
	if w.Keyboard[KeyRunRight].IsKeyPressed {
		w.Direction = DirectionRight
	}
	if w.Keyboard[KeyRunLeft].IsKeyPressed || w.Keyboard[KeyRunRight].IsKeyPressed {
		w.Run()
	} else {
		w.Stop()
	}
	switch {
	case w.IsDead:
		w.Status = StatusDeath
	case w.IsHurted:
		w.Status = StatusHurt
	case w.IsAttacking:
		w.Status = StatusAttack1
	case w.IsRunning:
		w.Status = StatusRun
	default:
		w.Status = StatusIdle
	}

	if w.IsRunning && !w.IsDead && !w.IsAttacking {
		switch w.Direction {
		case DirectionLeft:
			w.Side = 1.0
		case DirectionRight:
			w.Side = -1.0
		}
		w.X += w.SpeedRun * -w.Side
		if w.SpeedRun < w.MaxSpeedRun {
			w.SpeedRun += w.AccelerationRun
		} else {
			w.SpeedRun = w.MaxSpeedRun
		}
		w.SpeedRun += w.AccelerationRun
	} else {
		w.SpeedRun = w.BaseSpeedRun
	}

	if w.Status == StatusAttack && w.Frame == w.LastFrame/2 {
		if w.Side < 0 && ((w.X+float64(36*Scale))-(heroKnight.X+float64(35*Scale))) < float64(12*Scale) && ((w.X+float64(36*Scale))-(heroKnight.X+float64(35*Scale))) > -float64(12*Scale) {
			heroKnight.Hurt()
		}
		if w.Side > 0 && ((heroKnight.X+float64(65*Scale))-(w.X+float64(12*Scale))) < float64(12*Scale) && ((heroKnight.X+float64(65*Scale))-(w.X+float64(12*Scale))) > -float64(12*Scale) {
			heroKnight.Hurt()
		}
	}

	switch {
	case w.Status != w.PrevStatus:
		w.LastFrame = StatusFramesWizard[w.Status].FramesNumber*StatusFramesWizard[w.Status].FrameDuration - 1
		w.Frame = 0
	case w.Frame < w.LastFrame:
		w.Frame++
	}
	if w.Frame == w.LastFrame {
		w.IsAttacking = false
		w.IsHurted = false
		w.Frame = 0
	}
	fmt.Println(w.Frame)
	w.PrevStatus = w.Status
	return nil
}

func (w *Wizard) Draw(screen *ebiten.Image, unit *Unit, X, Y float64) {
	cameraX := X - float64(640*Scale-300)/2

	if cameraX < 0 {
		cameraX = 0
	}
	if cameraX > 355 {
		cameraX = 355
	}
	cameraY := -float64(360*Scale) / 2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-w.Side, 1.0)
	if w.Side > 0 {
		op.GeoM.Translate(unit.Width, 0.0)
	}

	op.GeoM.Translate(w.X-cameraX, -w.Y-cameraY-35)
	screen.DrawImage(unit.ActionFrames[w.Status][w.Frame/StatusFramesWizard[w.Status].FrameDuration], op)

	width := unit.Width
	if boxesShow {
		ebitenutil.DrawRect(screen, w.X, float64(TileSize)*9-unit.Height-w.Y, width, unit.Height, color.RGBA{0, 0, 255, 100})
		ebitenutil.DrawRect(screen, w.X+12, float64(TileSize)*9-unit.Height-w.Y, 24, unit.Height, color.RGBA{255, 0, 0, 100})
	}
}
