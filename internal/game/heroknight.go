package game

import (
	"fmt"
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
	Width            float64
	Height           float64
	Indent           float64
	Scale            float64
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
		frms = append(frms, ebiten.NewImageFromImage(img))
	}
	unit.ActionFrames["environment"] = frms

	return unit, err
}

func NewHeroKnight(x, y, width, height float64) *HeroKnight {
	baseSpeedRun := 3.0
	baseSpeedJump := 5.0
	return &HeroKnight{
		Keyboard:         NewDefaultKeyboard(),
		X:                x,
		Y:                y,
		Width:            width * float64(Scale),
		Height:           height * float64(Scale),
		Scale:            1.0,
		Indent:           35.0 * float64(Scale),
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
		hk.Health -= 20
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

func (hk *HeroKnight) Update(enemies map[string]*Enemy, environment []Environment) error {
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
		for _, env := range environment {
			fmt.Println(env.Y)
			fmt.Println(hk.Y)
			fmt.Println(hk.Height)
			if env.Y-(hk.Y+hk.Height) < 5 && env.Y-(hk.Y+hk.Height) > -5 {
				hk.SpeedJump = hk.BaseSpeedJump
				hk.Y = env.Y + hk.Height
				hk.IsJumping = false
				hk.IsFalling = false
			}
		}
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
			if hk.X+hk.Width-hk.Indent > unit.X+unit.Indent && hk.X+hk.Width-hk.Indent < unit.X+unit.Width-unit.Indent {
				hk.X -= unit.SpeedRun / 1.05
			}
			if hk.X+hk.Indent > unit.X+unit.Indent && hk.X+hk.Indent < unit.X+unit.Width-unit.Indent {
				hk.X += unit.SpeedRun / 1.05
			}
		}
		if hk.IsAttacking && hk.Frame == hk.LastFrame/2 {
			if hk.Side > 0 && ((hk.X+hk.Width-hk.Indent)-(unit.X+unit.Indent)) < hk.Indent && ((hk.X+hk.Width-hk.Indent)-(unit.X+unit.Indent)) > -hk.Indent {
				unit.Hurt()
			}
			if hk.Side < 0 && ((unit.X+unit.Width-unit.Indent)-(hk.X+hk.Indent)) < hk.Indent && ((unit.X+unit.Width-unit.Indent)-(hk.X+hk.Indent)) > -hk.Indent {
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

func (hk *HeroKnight) Draw(screen *ebiten.Image, unit *Unit, camera *Camera) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(hk.Side*hk.Scale*float64(Scale), hk.Scale*float64(Scale))
	if hk.Side < 0 {
		op.GeoM.Translate(hk.Width, 0.0)
	}
	op.GeoM.Translate(hk.X-camera.X, -hk.Y-camera.MainCharacterY-hk.Height+4.0)
	screen.DrawImage(unit.ActionFrames[hk.Status][hk.Frame/StatusFramesHeroKnight[hk.Status].FrameDuration], op)
	if boxesShow {
		ebitenutil.DrawRect(screen, hk.X-camera.X, -hk.Y-camera.MainCharacterY-hk.Height, hk.Width, hk.Height, color.RGBA{0, 0, 255, 100})
		ebitenutil.DrawRect(screen, hk.X-camera.X+hk.Indent, -hk.Y-camera.MainCharacterY-hk.Height, hk.Width-2*hk.Indent, hk.Height, color.RGBA{255, 0, 0, 100})
	}
}
