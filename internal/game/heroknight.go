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
	IsFalling        bool
	IsRolling        bool
	IsRunning        bool
}

func GetFramesHeroKnight() (*Unit, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	actionFrames := make(map[string][]*ebiten.Image)
	for status := range StatusFramesHeroKnight {
		framesNumber := StatusFramesHeroKnight[status].FramesNumber
		frms := make([]*ebiten.Image, 0, framesNumber)
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
			frms = append(frms, ebiten.NewImageFromImage(img))
		}
		actionFrames[status] = frms
	}
	file, _ = os.Open("_assets/HeroKnight/Idle/HeroKnight_Idle_0.png")
	cfg, _ = png.DecodeConfig(file)
	file.Close()
	unit := &Unit{
		ActionFrames: actionFrames,
		Width:        float64(cfg.Width),
		Height:       float64(cfg.Height),
	}

	frms := make([]*ebiten.Image, 0, 4)
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
		// 	file, err = os.Open("_assets/EnvironmentTiles/Tile_" + strconv.Itoa(i) + ".png")
		// 	if err != nil {
		// 		break
		// 	}
		// 	cfg, err = png.DecodeConfig(file)
		// 	if err != nil {
		// 		break
		// 	}
		// 	file.Close()
		frms = append(frms, ebiten.NewImageFromImage(img))
	}
	unit.ActionFrames["environment"] = frms

	return unit, err
}

func NewHeroKnight() *HeroKnight {
	baseSpeedRun := 3.0
	baseSpeedJump := 5.0
	return &HeroKnight{
		Keyboard:         NewDefaultKeyboard(),
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

func (hk *HeroKnight) Fall() {
	if hk.IsJumping && !hk.IsFalling {
		hk.IsFalling = true
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
	case hk.IsFalling:
		hk.Status = StatusFall
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
		if hk.SpeedJump < 0 {
			hk.Fall()
		}
		hk.Y += hk.SpeedJump
		hk.SpeedJump -= hk.DecelerationJump
		if hk.Y <= 0 {
			hk.SpeedJump = hk.BaseSpeedJump
			hk.Y = 0
			hk.IsJumping = false
			hk.IsFalling = false
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

	if hk.X < 0 {
		hk.X = 0
	}

	if hk.X+100 > 49*32 {
		hk.X = 49*32 - 100
	}

	hk.PrevStatus = hk.Status
	return nil
}

func (hk *HeroKnight) Draw(screen *ebiten.Image, unit *Unit) {

	cameraX := hk.X - float64(640*Scale-300)/2

	if cameraX < 0 {
		cameraX = 0
	}
	if cameraX > 355 {
		cameraX = 355
	}
	cameraY := hk.Y - float64(360*Scale)/2
	for i := 0; i < 50; i++ {
		op := &ebiten.DrawImageOptions{}

		op.GeoM.Scale(1.0*2, 1.0*2)
		op.GeoM.Translate(float64(TileSize*i*Scale)-cameraX, hk.Y-cameraY+106)
		screen.DrawImage(unit.ActionFrames["environment"][1], op)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(hk.Side*2, 1.0*2)
	if hk.Side < 0 {
		op.GeoM.Translate(unit.Width*2, 0.0)
	}
	// op.GeoM.Translate(hk.X, float64(TileSize*2)*9-hk.Frames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration].Height*2-hk.Y)
	op.GeoM.Translate(hk.X-cameraX, -hk.Y-cameraY)
	screen.DrawImage(unit.ActionFrames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration], op)
	w := unit.Width // - 35
	if boxesShow {
		ebitenutil.DrawRect(screen, hk.X, float64(TileSize)*9-unit.Height-hk.Y, w, unit.Height, color.RGBA{0, 0, 255, 100})
		ebitenutil.DrawRect(screen, hk.X+35, float64(TileSize)*9-unit.Height-hk.Y, 30, unit.Height, color.RGBA{255, 0, 0, 100})
	}
}
