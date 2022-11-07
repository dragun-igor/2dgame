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

func NewWizard() *Wizard {

	Frames, _ := GetFramesWizard()
	return &Wizard{
		Keyboard:         NewEmulationKeyboard(),
		Type:             "LightBandit",
		Frames:           Frames,
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
		LastFrame:        StatusFramesLightBandit[StatusCombatIdle].FramesNumber*StatusFramesLightBandit[StatusCombatIdle].FrameDuration - 1,
		Direction:        DirectionLeft,
	}
}

func GetFramesWizard() (map[string][]Frame, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	frames := map[string][]Frame{}
	frms := make([]Frame, 0, 8)
	file, err = os.Open("_assets/Wizard/Attack1.png")
	if err != nil {
		return map[string][]Frame{}, err
	}
	img, err = png.Decode(file)
	if err != nil {
		return map[string][]Frame{}, err
	}
	file.Close()
	fmt.Println("I'm here")
	for i := 0; i < 8; i++ {
		frms = append(frms, Frame{
			Img:    ebiten.NewImageFromImage(img).SubImage(image.Rect(231*i, 0, 231*i+231, 190)).(*ebiten.Image),
			Width:  float64(cfg.Width),
			Height: float64(cfg.Height),
		})
	}
	fmt.Println(frms)
	frames["Attack1"] = frms
	return frames, err
}

func (w *Wizard) Update() error {
	return nil
}

func (w *Wizard) Draw(screen *ebiten.Image) {
	screen.DrawImage(w.Frames["Attack1"][w.Frame/7%6].Img, nil)
	if boxesShow {
		ebitenutil.DrawRect(screen, w.X, w.Y, 231, 190, color.RGBA{0, 0, 255, 100})
	}
	w.Frame++
}

// screen.DrawImage(img.SubImage(image.Rect(0, 0, 231, 190)).(*ebiten.Image), nil)
// screen.DrawImage(img.SubImage(image.Rect(231, 0, 231*2, 190)).(*ebiten.Image), nil)
// screen.DrawImage(img.SubImage(image.Rect(231*2, 0, 231*3, 190)).(*ebiten.Image), nil)
// screen.DrawImage(img.SubImage(image.Rect(231*3, 0, 231*4, 190)).(*ebiten.Image), nil)
// screen.DrawImage(img.SubImage(image.Rect(231*4, 0, 231*5, 190)).(*ebiten.Image), nil)
// screen.DrawImage(img.SubImage(image.Rect(231*5, 0, 231*6, 190)).(*ebiten.Image), nil)
// screen.DrawImage(img.SubImage(image.Rect(231*6, 0, 231*7, 190)).(*ebiten.Image), nil)
// screen.DrawImage(img.SubImage(image.Rect(231*7, 0, 231*8, 190)).(*ebiten.Image), nil)
