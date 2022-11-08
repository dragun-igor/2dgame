package game

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	Keyboard         Keyboard
	Status           string
	PrevStatus       string
	Type             string
	X                float64
	Y                float64
	Scale            float64
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

func NewEnemy(t string) *Enemy {
	var enemy *Enemy
	switch t {
	case LightBandit:
		enemy = NewLightBandit()
	case HeavyBandit:
		enemy = NewHeavyBandit()
	case Wizard:
		enemy = NewWizard()
	}
	return enemy
}

func NewLightBandit() *Enemy {
	return &Enemy{
		Keyboard:         NewEmulationKeyboard(),
		Type:             LightBandit,
		X:                550,
		Y:                0,
		Scale:            1.0,
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
		LastFrame:        StatusFrames[LightBandit][StatusIdle].FramesNumber*StatusFrames[LightBandit][StatusIdle].FrameDuration - 1,
		Direction:        DirectionLeft,
	}
}

func NewHeavyBandit() *Enemy {
	return &Enemy{
		Keyboard:         NewEmulationKeyboard(),
		Type:             HeavyBandit,
		X:                300,
		Y:                0,
		Scale:            1.0,
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
		LastFrame:        StatusFrames[HeavyBandit][StatusIdle].FramesNumber*StatusFrames[HeavyBandit][StatusIdle].FrameDuration - 1,
		Direction:        DirectionLeft,
	}
}

func NewWizard() *Enemy {
	return &Enemy{
		Keyboard:         NewEmulationKeyboard(),
		Type:             Wizard,
		X:                250,
		Y:                0,
		Scale:            0.5,
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
		LastFrame:        StatusFrames[Wizard][StatusIdle].FramesNumber*StatusFrames[Wizard][StatusIdle].FrameDuration - 1,
		Direction:        DirectionLeft,
	}
}

func GetFramesBandit(strType string) (*Unit, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	actionFrames := make(map[string][]*ebiten.Image)
	for status := range StatusFramesLightBandit {
		framesNumber := StatusFramesLightBandit[status].FramesNumber
		frms := make([]*ebiten.Image, 0, framesNumber)
		for i := 0; i < framesNumber; i++ {
			file, err = os.Open("_assets/" + strType + "/" + status + "/" + strType + "_" + status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				break
			}

			size := img.Bounds().Size()
			r := image.NewRGBA(image.Rect(0, 0, size.X, size.Y))
			for i := 0; i < size.Y; i++ {
				for j := 0; j < size.X; j++ {
					r.Set(j, i, img.At(size.X-1-j, i))
				}
				// reverseImage = append(reverseImage, row)
			}

			file.Close()
			frms = append(frms, ebiten.NewImageFromImage(r))
		}
		actionFrames[status] = frms
	}
	file, _ = os.Open("_assets/" + strType + "/Idle/" + strType + "_Idle_0.png")
	cfg, _ = png.DecodeConfig(file)
	file.Close()
	unit := &Unit{
		ActionFrames: actionFrames,
		Width:        float64(cfg.Width),
		Height:       float64(cfg.Height),
	}
	return unit, err
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

func (e *Enemy) Death() {
	if e.Health <= 0 {
		e.IsDead = true
	}
}

func (e *Enemy) Hurt() {
	if !e.IsHurted {
		e.IsHurted = true
		e.Health -= 30
	}
}

func (e *Enemy) Attack() {
	if !e.IsAttacking && !e.IsDead {
		e.IsAttacking = true
	}
}

func (e *Enemy) Run() {
	e.IsRunning = true
}

func (e *Enemy) Stop() {
	e.IsRunning = false
}

func (e *Enemy) Update(heroKnight *HeroKnight) error {
	e.Keyboard[KeyAttack].KeyEmulation = e.AttackAction
	e.Keyboard[KeyRunLeft].KeyEmulation = e.RunLeftAction
	e.Keyboard[KeyRunRight].KeyEmulation = e.RunRightAction
	e.Keyboard.Update()
	e.Death()
	if e.Keyboard[KeyAttack].IsKeyJustPressed {
		e.Attack()
	}
	if e.Keyboard[KeyRunLeft].IsKeyPressed {
		e.Direction = DirectionLeft
	}
	if e.Keyboard[KeyRunRight].IsKeyPressed {
		e.Direction = DirectionRight
	}
	if e.Keyboard[KeyRunLeft].IsKeyPressed || e.Keyboard[KeyRunRight].IsKeyPressed {
		e.Run()
	} else {
		e.Stop()
	}
	switch {
	case e.IsDead:
		e.Status = StatusDeath
	case e.IsHurted:
		e.Status = StatusHurt
	case e.IsAttacking:
		switch e.Type {
		case LightBandit, HeavyBandit:
			e.Status = StatusAttack
		case Wizard:
			e.Status = StatusAttack1
		}
	case e.IsRunning:
		e.Status = StatusRun
	default:
		e.Status = StatusIdle
	}

	if e.IsRunning && !e.IsDead && !e.IsAttacking {
		switch e.Direction {
		case DirectionLeft:
			e.Side = -1.0
		case DirectionRight:
			e.Side = 1.0
		}
		e.X += e.SpeedRun * e.Side
		if e.SpeedRun < e.MaxSpeedRun {
			e.SpeedRun += e.AccelerationRun
		} else {
			e.SpeedRun = e.MaxSpeedRun
		}
		e.SpeedRun += e.AccelerationRun
	} else {
		e.SpeedRun = e.BaseSpeedRun
	}

	if (e.Status == StatusAttack || e.Status == StatusAttack1) && e.Frame == e.LastFrame/2 {
		if e.Side < 0 && ((e.X+float64(36*Scale))-(heroKnight.X+float64(35*Scale))) < float64(12*Scale) && ((e.X+float64(36*Scale))-(heroKnight.X+float64(35*Scale))) > -float64(12*Scale) {
			heroKnight.Hurt()
		}
		if e.Side > 0 && ((heroKnight.X+float64(65*Scale))-(e.X+float64(12*Scale))) < float64(12*Scale) && ((heroKnight.X+float64(65*Scale))-(e.X+float64(12*Scale))) > -float64(12*Scale) {
			heroKnight.Hurt()
		}
	}

	switch {
	case e.Status != e.PrevStatus:
		e.LastFrame = StatusFrames[e.Type][e.Status].FramesNumber*StatusFrames[e.Type][e.Status].FrameDuration - 1
		e.Frame = 0
	case e.Frame < e.LastFrame:
		e.Frame++
	}
	if e.Frame == e.LastFrame && !e.IsDead {
		e.IsAttacking = false
		e.IsHurted = false
		e.Frame = 0
	}

	e.PrevStatus = e.Status
	if e.Type == Wizard {
		fmt.Println(*e)
	}

	return nil
}

func (e *Enemy) Draw(screen *ebiten.Image, unit *Unit, camera *Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(e.Side*e.Scale*float64(Scale), e.Scale*float64(Scale))
	if e.Side < 0 {
		op.GeoM.Translate(unit.Width*float64(Scale)*e.Scale, 0.0)
	}
	offsetX := e.X - camera.X
	offsetY := -e.Y - camera.EnemyY - unit.Height*float64(Scale)*e.Scale
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(unit.ActionFrames[e.Status][e.Frame/StatusFrames[e.Type][e.Status].FrameDuration], op)
	if boxesShow {
		ebitenutil.DrawRect(screen, offsetX, offsetY, unit.Width*float64(Scale)*e.Scale, unit.Height*float64(Scale)*e.Scale, color.RGBA{0, 0, 255, 100})
		ebitenutil.DrawRect(screen, offsetX+12*float64(Scale)*e.Scale, offsetY, 24*float64(Scale)*e.Scale, unit.Height*float64(Scale)*e.Scale, color.RGBA{255, 0, 0, 100})
	}
}
