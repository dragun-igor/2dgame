package game

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Bandit struct {
	Keyboard         Keyboard
	Frames           map[string][]Frame
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

func NewLightBandit() *Bandit {
	Frames, _ := GetFramesBandit("LightBandit")
	return &Bandit{
		Keyboard:         NewEmulationKeyboard(),
		Type:             "LightBandit",
		Frames:           Frames,
		X:                550,
		Y:                0,
		Status:           StatusCombatIdle,
		PrevStatus:       StatusCombatIdle,
		Side:             1.0,
		SpeedRun:         1.0,
		MaxSpeedRun:      2.0,
		BaseSpeedRun:     1.0,
		AccelerationRun:  0.2,
		SpeedJump:        6.0,
		BaseSpeedJump:    6.0,
		DecelerationJump: 0.2,
		Health:           100,
		LastFrame:        StatusFramesLightBandit[StatusCombatIdle].FramesNumber*StatusFramesLightBandit[StatusCombatIdle].FrameDuration - 1,
		Direction:        DirectionLeft,
	}
}

func NewHeavyBandit() *Bandit {
	Frames, _ := GetFramesBandit("HeavyBandit")
	return &Bandit{
		Keyboard:         NewEmulationKeyboard(),
		Type:             "HeavyBandit",
		Frames:           Frames,
		X:                300,
		Y:                0,
		Status:           StatusCombatIdle,
		PrevStatus:       StatusCombatIdle,
		Side:             1.0,
		SpeedRun:         1.0,
		MaxSpeedRun:      2.0,
		BaseSpeedRun:     1.0,
		AccelerationRun:  0.2,
		SpeedJump:        6.0,
		BaseSpeedJump:    6.0,
		DecelerationJump: 0.2,
		Health:           100,
		LastFrame:        StatusFramesLightBandit[StatusCombatIdle].FramesNumber*StatusFramesLightBandit[StatusCombatIdle].FrameDuration - 1,
		Direction:        DirectionLeft,
	}
}

func GetFramesBandit(strType string) (map[string][]Frame, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	frames := map[string][]Frame{}
	for status := range StatusFramesLightBandit {
		framesNumber := StatusFramesLightBandit[status].FramesNumber
		frms := make([]Frame, 0, framesNumber)
		for i := 0; i < framesNumber; i++ {
			file, err = os.Open("_assets/" + strType + "/" + status + "/" + strType + "_" + status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				break
			}
			file.Close()
			file, err = os.Open("_assets/" + strType + "/" + status + "/" + strType + "_" + status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			cfg, err = png.DecodeConfig(file)
			if err != nil {
				break
			}
			file.Close()
			frms = append(frms, Frame{
				Img:    ebiten.NewImageFromImage(img),
				Width:  float64(cfg.Width),
				Height: float64(cfg.Height),
			})
		}
		frames[status] = frms
	}
	return frames, err
}

func (b *Bandit) Death() {
	if b.Health <= 0 {
		b.IsDead = true
	}
}

func (b *Bandit) Hurt() {
	if !b.IsHurted {
		b.IsHurted = true
		b.Health -= 30
	}
}

func (b *Bandit) Attack() {
	if !b.IsAttacking && !b.IsDead {
		b.IsAttacking = true
	}
}

func (b *Bandit) Run() {
	b.IsRunning = true
}

func (b *Bandit) Stop() {
	b.IsRunning = false
}

func (b *Bandit) Update(heroKnight *HeroKnight) error {
	b.Keyboard[KeyAttack].KeyEmulation = b.AttackAction
	b.Keyboard[KeyRunLeft].KeyEmulation = b.RunLeftAction
	b.Keyboard[KeyRunRight].KeyEmulation = b.RunRightAction
	b.Keyboard.Update()
	b.Death()
	if b.Keyboard[KeyAttack].IsKeyJustPressed {
		b.Attack()
	}
	if b.Keyboard[KeyRunLeft].IsKeyPressed {
		b.Direction = DirectionLeft
	}
	if b.Keyboard[KeyRunRight].IsKeyPressed {
		b.Direction = DirectionRight
	}
	if b.Keyboard[KeyRunLeft].IsKeyPressed || b.Keyboard[KeyRunRight].IsKeyPressed {
		b.Run()
	} else {
		b.Stop()
	}
	switch {
	case b.IsDead:
		b.Status = StatusDeath
	case b.IsHurted:
		b.Status = StatusHurt
	case b.IsAttacking:
		b.Status = StatusAttack
	case b.IsRunning:
		b.Status = StatusRun
	default:
		b.Status = StatusCombatIdle
	}

	if b.IsRunning && !b.IsDead && !b.IsAttacking {
		switch b.Direction {
		case DirectionLeft:
			b.Side = 1.0
		case DirectionRight:
			b.Side = -1.0
		}
		b.X += b.SpeedRun * -b.Side
		if b.SpeedRun < b.MaxSpeedRun {
			b.SpeedRun += b.AccelerationRun
		} else {
			b.SpeedRun = b.MaxSpeedRun
		}
		b.SpeedRun += b.AccelerationRun
	} else {
		b.SpeedRun = b.BaseSpeedRun
	}

	if b.Status == StatusAttack && b.Frame == b.LastFrame/2 {
		if b.Side < 0 && ((b.X+float64(36*Scale))-(heroKnight.X+float64(35*Scale))) < float64(12*Scale) && ((b.X+float64(36*Scale))-(heroKnight.X+float64(35*Scale))) > -float64(12*Scale) {
			heroKnight.Hurt()
		}
		if b.Side > 0 && ((heroKnight.X+float64(65*Scale))-(b.X+float64(12*Scale))) < float64(12*Scale) && ((heroKnight.X+float64(65*Scale))-(b.X+float64(12*Scale))) > -float64(12*Scale) {
			heroKnight.Hurt()
		}
	}

	switch {
	case b.Status != b.PrevStatus:
		b.LastFrame = StatusFramesLightBandit[b.Status].FramesNumber*StatusFramesLightBandit[b.Status].FrameDuration - 1
		b.Frame = 0
	case b.Frame < b.LastFrame:
		b.Frame++
	}
	if b.Frame == b.LastFrame {
		b.IsAttacking = false
		b.IsHurted = false
		b.Frame = 0
	}

	b.PrevStatus = b.Status
	return nil
}

func (b *Bandit) Draw(screen *ebiten.Image, X, Y float64) {
	cameraX := X - float64(640*Scale-300)/2
	cameraY := Y - float64(360*Scale)/2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(b.Side*2, 1.0*2)
	if b.Side < 0 {
		op.GeoM.Translate(b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Width*2, 0.0)
	}
	// op.GeoM.Translate(b.X, float64(TileSize*2)*9-b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Height*2-b.Y)
	op.GeoM.Translate(b.X-cameraX, -b.Y-cameraY)
	screen.DrawImage(b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Img, op)

	w := b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Width
	if boxesShow {
		ebitenutil.DrawRect(screen, b.X, float64(TileSize)*9-b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Height-b.Y, w, b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Height, color.RGBA{0, 0, 255, 100})
		ebitenutil.DrawRect(screen, b.X+12, float64(TileSize)*9-b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Height-b.Y, 24, b.Frames[b.Status][b.Frame/StatusFramesLightBandit[b.Status].FrameDuration].Height, color.RGBA{255, 0, 0, 100})
	}
}
