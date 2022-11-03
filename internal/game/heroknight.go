package game

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type HeroKnight struct {
	Keyboard         Keyboard
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

func GetFramesHeroKnight() (map[string][]Frame, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	frames := map[string][]Frame{}
	for status := range StatusFramesHeroKnight {
		framesNumber := StatusFramesHeroKnight[status].FramesNumber
		frms := make([]Frame, 0, framesNumber)
		for i := 0; i < framesNumber; i++ {
			file, err = os.Open("_assets/HeroKnight/" + status + "/HeroKnight_" + status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				break
			}
			file.Close()
			file, err = os.Open("_assets/HeroKnight/" + status + "/HeroKnight_" + status + "_" + strconv.Itoa(i) + ".png")
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
	frames["environment"] = frms

	return frames, err
}

func NewHeroKnight() *HeroKnight {
	baseSpeedRun := 3.0
	baseSpeedJump := 5.0
	Frames, _ := GetFramesHeroKnight()
	return &HeroKnight{
		Keyboard:         NewDefaultKeyboard(),
		Frames:           Frames,
		X:                50,
		Status:           StatusIdle,
		Side:             1.0,
		AttackType:       AttackType1,
		SpeedRun:         baseSpeedRun,
		MaxSpeedRun:      5.0,
		BaseSpeedRun:     baseSpeedRun,
		AccelerationRun:  0.03,
		Health:           100,
		Stamina:          100,
		SpeedJump:        baseSpeedJump,
		BaseSpeedJump:    baseSpeedJump,
		DecelerationJump: 0.3,
		Direction:        DirectionRight,
		SpeedRoll:        5.0,
		LastFrame:        StatusFramesHeroKnight[StatusIdle].FrameDuration * StatusFramesHeroKnight[StatusIdle].FramesNumber,
	}
}

func (hk *HeroKnight) Death() {
	if hk.Health <= 0 {
		hk.IsDead = true
	}
}

func (hk *HeroKnight) Hurt() {
	if !hk.IsHurted {
		hk.IsHurted = true
		hk.Health -= 50
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

func (hk *HeroKnight) Update(enemies map[string]*Bandit) error {
	hk.Keyboard.Update()
	hk.Death()
	if hk.Keyboard[KeyAttack].IsKeyJustPressed {
		hk.Attack()
	}
	if hk.Keyboard[KeyJump].IsKeyJustPressed {
		hk.Jump()
	}
	if hk.Keyboard[KeyRoll].IsKeyJustPressed {
		hk.Roll()
	}
	if hk.Keyboard[KeyRunLeft].IsKeyPressed {
		hk.Direction = DirectionLeft
	}
	if hk.Keyboard[KeyRunRight].IsKeyPressed {
		hk.Direction = DirectionRight
	}
	if hk.Keyboard[KeyRunLeft].IsKeyPressed || hk.Keyboard[KeyRunRight].IsKeyPressed {
		hk.Run()
	} else {
		hk.Stop()
	}

	switch {
	case hk.IsDead:
		hk.Status = StatusDeath
	case hk.IsHurted:
		hk.Status = StatusHurt
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

	for _, unit := range enemies {
		if !hk.IsJumping && !hk.IsRolling && !unit.IsDead {
			if hk.X+float64(65*Scale) > unit.X+float64(12*Scale) && hk.X+float64(65*Scale) < unit.X+float64(36*Scale) {
				hk.X -= unit.SpeedRun / 1.5
			}
			if hk.X+float64(35*Scale) > unit.X+float64(12*Scale) && hk.X+float64(35*Scale) < unit.X+float64(36*Scale) {
				hk.X += unit.SpeedRun / 1.5
			}
		}
		if hk.IsAttacking && hk.Frame == hk.LastFrame/2 {
			if hk.Side > 0 && ((hk.X+float64(65*Scale))-(unit.X+float64(12*Scale))) < float64(35*Scale) && ((hk.X+float64(65*Scale))-(unit.X+float64(12*Scale))) > -float64(35*Scale) {
				unit.Hurt()
			}
			if hk.Side < 0 && ((unit.X+float64(36*Scale))-(hk.X+float64(35*Scale))) < float64(35*Scale) && ((unit.X+float64(36*Scale))-(hk.X+float64(35*Scale))) > -float64(35*Scale) {
				unit.Hurt()
			}
		}
	}

	switch {
	case hk.Status != hk.PrevStatus:
		hk.LastFrame = StatusFramesHeroKnight[hk.Status].FramesNumber*StatusFramesHeroKnight[hk.Status].FrameDuration - 1
		hk.Frame = 0
	case hk.Frame < hk.LastFrame:
		hk.Frame++
	}

	if hk.Frame == hk.LastFrame && hk.Status != StatusDeath {
		if hk.IsAttacking {
			hk.Frame = 0
			hk.IsAttacking = false
		}
		if hk.IsRolling {
			hk.Frame = 0
			hk.IsRolling = false
		}
		if hk.IsHurted {
			hk.Frame = 0
			hk.IsHurted = false
		}
	}

	hk.PrevStatus = hk.Status
	return nil
}

func (hk *HeroKnight) Draw(screen *ebiten.Image) {
	cameraX := hk.X - float64(640*Scale-300)/2
	cameraY := hk.Y - float64(360*Scale)/2
	// for i := 0; i < 50; i++ {
	// 	op := &ebiten.DrawImageOptions{}
	// 	op.GeoM.Translate(float64(TileSize*i), float64(TileSize)*9)
	// 	op.GeoM.Scale(1.0*2, 1.0*2)
	// 	screen.DrawImage(hk.Frames["environment"][1].Img, op)
	// }
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(hk.Side*2, 1.0*2)
	if hk.Side < 0 {
		op.GeoM.Translate(hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Width*2, 0.0)
	}
	// op.GeoM.Translate(hk.X, float64(TileSize*2)*9-hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Height*2-hk.Y)
	op.GeoM.Translate(hk.X-cameraX, -hk.Y-cameraY)
	screen.DrawImage(hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Img, op)
	w := hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Width // - 35
	if boxesShow {
		ebitenutil.DrawRect(screen, hk.X, float64(TileSize)*9-hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Height-hk.Y, w, hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Height, color.RGBA{0, 0, 255, 100})
		ebitenutil.DrawRect(screen, hk.X+35, float64(TileSize)*9-hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Height-hk.Y, 30, hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Height, color.RGBA{255, 0, 0, 100})
	}
}
