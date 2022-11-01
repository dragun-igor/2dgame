package game

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type LightBandit struct {
	Frames           map[string][]Frame
	Status           string
	PrevStatus       string
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
}

func NewLightBandit() *LightBandit {
	Frames, _ := GetFrames()
	return &LightBandit{
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
		LastFrame:        (StatusFrames[StatusCombatIdle].FramesNumber - 1) * StatusFrames[StatusCombatIdle].FrameDuration,
		Direction:        DirectionLeft,
	}
}

func GetFramesLightBandit() (map[string][]Frame, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	Frames := map[string][]Frame{}
	for Status := range StatusFrames {
		FramesNumber := StatusFrames[Status].FramesNumber
		frms := make([]Frame, 0, FramesNumber)
		for i := 0; i < FramesNumber; i++ {
			file, err = os.Open("_assets/LightBandit/" + Status + "/LightBandit_" + Status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				break
			}
			file.Close()
			file, err = os.Open("_assets/LightBandit/" + Status + "/LightBandit_" + Status + "_" + strconv.Itoa(i) + ".png")
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
		Frames[Status] = frms
	}
	return Frames, err
}

func (lb *LightBandit) Death() {
	if !lb.IsDead {
		lb.IsDead = true
	}
}

func (lb *LightBandit) Hurt() {
	if !lb.IsHurted {
		lb.IsHurted = true
		lb.Health -= 2
	}
}

func (lb *LightBandit) Update() error {
	if lb.Health <= 0 {
		lb.Death()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		lb.Hurt()
	}

	switch {
	case lb.IsDead:
		lb.Status = StatusDeath
	case lb.IsHurted:
		lb.Status = StatusHurt
	default:
		lb.Status = StatusCombatIdle
	}

	switch {
	case lb.Status != lb.PrevStatus:
		lb.LastFrame = StatusFrames[lb.Status].FramesNumber*StatusFrames[lb.Status].FrameDuration - 1
		lb.Frame = 0
	case lb.Frame < lb.LastFrame:
		lb.Frame++
	}
	if lb.Frame == lb.LastFrame {
		lb.IsHurted = false
		lb.Frame = 0
	}
	lb.PrevStatus = lb.Status
	return nil
}

func (lb *LightBandit) Draw(screen *ebiten.Image) {
	tileSize := 32
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(lb.Side, 1.0)
	if lb.Side < 0 {
		op.GeoM.Translate(lb.Frames[lb.Status][lb.Frame/StatusFrames[lb.Status].FrameDuration].Width, 0.0)
	}
	op.GeoM.Translate(lb.X, float64(tileSize)*9-lb.Frames[lb.Status][lb.Frame/StatusFrames[lb.Status].FrameDuration].Height-lb.Y)
	screen.DrawImage(lb.Frames[lb.Status][lb.Frame/StatusFrames[lb.Status].FrameDuration].Img, op)
	w := lb.Frames[lb.Status][lb.Frame/StatusFrames[lb.Status].FrameDuration].Width
	ebitenutil.DrawRect(screen, lb.X, float64(tileSize)*9-lb.Frames[lb.Status][lb.Frame/StatusFrames[lb.Status].FrameDuration].Height-lb.Y, w, lb.Frames[lb.Status][lb.Frame/StatusFrames[lb.Status].FrameDuration].Height, color.RGBA{0, 0, 255, 20})

}
