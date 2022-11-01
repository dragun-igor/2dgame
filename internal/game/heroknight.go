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

type HeroKnight struct {
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
	SpeedRoll        float64
	Frame            int
	LastFrame        int
	Health           int
	Stamina          int
	AttackType       atype
	Direction        uint8
	IsDead           bool
	IsHurted         bool
	IsAttacking      bool
	IsBlocking       bool
	IsJumping        bool
	IsRolling        bool
	IsRunning        bool
}

func GetFrames() (map[string][]Frame, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	Frames := map[string][]Frame{}
	for Status := range StatusFrames {
		FramesNumber := StatusFrames[Status].FramesNumber
		frms := make([]Frame, 0, FramesNumber)
		for i := 0; i < FramesNumber; i++ {
			file, err = os.Open("_assets/HeroKnight/" + Status + "/HeroKnight_" + Status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				break
			}
			file.Close()
			file, err = os.Open("_assets/HeroKnight/" + Status + "/HeroKnight_" + Status + "_" + strconv.Itoa(i) + ".png")
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

	frms := make([]Frame, 0, 4)
	for i := 0; i < 4; i++ {
		file, err = os.Open("_assets/EnvironmentTiles/Tile_" + strconv.Itoa(i) + ".png")
		if err != nil {
			break
		}
		img, err = png.Decode(file)
		if err != nil {
			break
		}
		file.Close()
		file, err = os.Open("_assets/EnvironmentTiles/Tile_" + strconv.Itoa(i) + ".png")
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
	Frames["environment"] = frms

	return Frames, err
}

func NewHeroKnight() *HeroKnight {
	baseSpeedRun := 1.0
	baseSpeedJump := 6.0
	Frames, _ := GetFrames()
	return &HeroKnight{
		Frames:           Frames,
		X:                50,
		Status:           StatusIdle,
		Side:             1.0,
		AttackType:       AttackType1,
		SpeedRun:         baseSpeedRun,
		MaxSpeedRun:      2.5,
		BaseSpeedRun:     baseSpeedRun,
		AccelerationRun:  0.03,
		Health:           100,
		Stamina:          100,
		SpeedJump:        baseSpeedJump,
		BaseSpeedJump:    baseSpeedJump,
		DecelerationJump: 0.2,
		Direction:        DirectionRight,
		SpeedRoll:        5.0,
	}
}

func (hk *HeroKnight) Death() {
	if hk.Health <= 0 {
		hk.IsDead = true
	}
}

func (hk *HeroKnight) Attack() {
	if !hk.IsAttacking && !hk.IsRolling && !hk.IsDead {
		hk.IsAttacking = true
		switch hk.AttackType {
		case AttackType1:
			hk.AttackType = AttackType2
		case AttackType2:
			hk.AttackType = AttackType3
		case AttackType3:
			hk.AttackType = AttackType1
		}
	}
}

func (hk *HeroKnight) Jump() {
	if !hk.IsJumping && !hk.IsRolling && !hk.IsAttacking {
		hk.IsJumping = true
	}
}

func (hk *HeroKnight) Roll() {
	if !hk.IsRolling && !hk.IsAttacking && !hk.IsJumping {
		hk.IsRolling = true
	}
}

func (hk *HeroKnight) Run() {
	hk.IsRunning = true
}

func (hk *HeroKnight) Stop() {
	hk.IsRunning = false
}

func (hk *HeroKnight) Update() error {
	hk.Death()
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		hk.Health = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		hk.Attack()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		hk.Jump()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyControlLeft) {
		hk.Roll()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		hk.Direction = DirectionLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		hk.Direction = DirectionRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
		hk.Run()
	} else {
		hk.Stop()
	}

	switch {
	case hk.IsDead:
		hk.Status = StatusDeath
	case hk.IsAttacking:
		switch hk.AttackType {
		case AttackType1:
			hk.Status = StatusAttack1
		case AttackType2:
			hk.Status = StatusAttack2
		case AttackType3:
			hk.Status = StatusAttack3
		}
	case hk.IsRolling:
		hk.Status = StatusRoll
	case hk.IsJumping:
		hk.Status = StatusJump
	case hk.IsRunning:
		hk.Status = StatusRun
		if hk.Frame == hk.LastFrame {
			hk.Frame = 0
		}
	default:
		hk.Status = StatusIdle
		if hk.Frame == hk.LastFrame {
			hk.Frame = 0
		}
	}

	if (hk.IsJumping && !hk.IsDead) || (hk.IsDead && hk.Y > 0) {
		hk.Y += hk.SpeedJump
		hk.SpeedJump -= hk.DecelerationJump
		if hk.Y <= 0 {
			hk.SpeedJump = hk.BaseSpeedJump
			hk.Y = 0
			hk.IsJumping = false
		}
	}

	if hk.IsRolling {
		hk.X += hk.SpeedRoll * hk.Side
	}

	if hk.IsRunning && !hk.IsDead && !hk.IsAttacking && !hk.IsRolling {
		switch hk.Direction {
		case DirectionLeft:
			hk.Side = -1.0
		case DirectionRight:
			hk.Side = 1.0
		}
		hk.X += hk.SpeedRun * hk.Side
		if hk.SpeedRun < hk.MaxSpeedRun {
			hk.SpeedRun += hk.AccelerationRun
		} else {
			hk.SpeedRun = hk.MaxSpeedRun
		}
		hk.SpeedRun += hk.AccelerationRun
	} else {
		hk.SpeedRun = hk.BaseSpeedRun
	}

	switch {
	case hk.Status != hk.PrevStatus:
		hk.LastFrame = StatusFrames[hk.Status].FramesNumber*StatusFrames[hk.Status].FrameDuration - 1
		hk.Frame = 0
	case hk.Frame < hk.LastFrame:
		hk.Frame++
	}

	if hk.Frame == hk.LastFrame {
		if hk.IsAttacking {
			hk.Frame = 0
			hk.IsAttacking = false
		}
		if hk.IsRolling {
			hk.Frame = 0
			hk.IsRolling = false
		}
	}

	hk.PrevStatus = hk.Status
	return nil
}

func (hk *HeroKnight) Draw(screen *ebiten.Image) {
	tileSize := 32
	for i := 0; i < 20; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileSize*i), float64(tileSize)*9)
		screen.DrawImage(hk.Frames["environment"][1].Img, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(hk.Side, 1.0)
	if hk.Side < 0 {
		op.GeoM.Translate(hk.Frames[hk.Status][hk.Frame/StatusFrames[hk.Status].FrameDuration].Width, 0.0)
	}
	op.GeoM.Translate(hk.X, float64(tileSize)*9-hk.Frames[hk.Status][hk.Frame/StatusFrames[hk.Status].FrameDuration].Height-hk.Y)
	screen.DrawImage(hk.Frames[hk.Status][hk.Frame/StatusFrames[hk.Status].FrameDuration].Img, op)
	w := hk.Frames[hk.Status][hk.Frame/StatusFrames[hk.Status].FrameDuration].Width - 35
	ebitenutil.DrawRect(screen, hk.X, float64(tileSize)*9-hk.Frames[hk.Status][hk.Frame/StatusFrames[hk.Status].FrameDuration].Height-hk.Y, w, hk.Frames[hk.Status][hk.Frame/StatusFrames[hk.Status].FrameDuration].Height, color.RGBA{0, 0, 255, 20})
}
