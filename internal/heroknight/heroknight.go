package heroknight

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

const (
	statusAttack1       string = "Attack1"
	statusAttack2       string = "Attack2"
	statusAttack3       string = "Attack3"
	statusBlock         string = "Block"
	statusBlockIdle     string = "BlockIdle"
	statusBlockNoEffect string = "BlockNoEffect"
	statusDeath         string = "Death"
	statusDeathNoBlood  string = "DeathNoBlood"
	statusFall          string = "Fall"
	statusHurt          string = "Hurt"
	statusIdle          string = "Idle"
	statusJump          string = "Jump"
	statusLedgeGrab     string = "LedgeGrab"
	statusRoll          string = "Roll"
	statusRun           string = "Run"
	statusWallSide      string = "WallSide"
)

type statusData struct {
	framesNumber  int
	frameDuration int
}

var statusFrames map[string]statusData = map[string]statusData{
	statusAttack1: statusData{
		framesNumber:  6,
		frameDuration: 4,
	},
	statusAttack2: statusData{
		framesNumber:  6,
		frameDuration: 4,
	},
	statusAttack3: statusData{
		framesNumber:  8,
		frameDuration: 4,
	},
	statusBlock: statusData{
		framesNumber:  5,
		frameDuration: 7,
	},
	statusBlockIdle: statusData{
		framesNumber:  8,
		frameDuration: 7,
	},
	statusBlockNoEffect: statusData{
		framesNumber:  5,
		frameDuration: 7,
	},
	statusDeath: statusData{
		framesNumber:  10,
		frameDuration: 7,
	},
	statusDeathNoBlood: statusData{
		framesNumber:  10,
		frameDuration: 7,
	},
	statusFall: statusData{
		framesNumber:  4,
		frameDuration: 7,
	},
	statusHurt: statusData{
		framesNumber:  3,
		frameDuration: 7,
	},
	statusIdle: statusData{
		framesNumber:  8,
		frameDuration: 7,
	},
	statusJump: statusData{
		framesNumber:  3,
		frameDuration: 2,
	},
	statusLedgeGrab: statusData{
		framesNumber:  5,
		frameDuration: 7,
	},
	statusRoll: statusData{
		framesNumber:  9,
		frameDuration: 4,
	},
	statusRun: statusData{
		framesNumber:  10,
		frameDuration: 7,
	},
	statusWallSide: statusData{
		framesNumber:  5,
		frameDuration: 7,
	},
}

type Frame struct {
	img    *ebiten.Image
	width  float64
	height float64
}

type atype uint8

const (
	attackType1 atype = iota + 1
	attackType2
	attackType3
)

const (
	directionLeft uint8 = iota + 1
	directionRight
)

type HeroKnight struct {
	frames           map[string][]Frame
	Status           string
	prevStatus       string
	X                float64
	y                float64
	Side             float64
	speedRun         float64
	maxSpeedRun      float64
	baseSpeedRun     float64
	accelerationRun  float64
	speedJump        float64
	baseSpeedJump    float64
	decelerationJump float64
	speedRoll        float64
	Frame            int
	LastFrame        int
	health           int
	stamina          int
	attackType       atype
	direction        uint8
	isDead           bool
	isHurted         bool
	IsAttacking      bool
	isBlocking       bool
	isJumping        bool
	isRolling        bool
	isRunning        bool
}

func GetFrames() (map[string][]Frame, error) {
	var file *os.File
	var img image.Image
	var cfg image.Config
	var err error
	frames := map[string][]Frame{}
	for status := range statusFrames {
		framesNumber := statusFrames[status].framesNumber
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
				img:    ebiten.NewImageFromImage(img),
				width:  float64(cfg.Width),
				height: float64(cfg.Height),
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
			img:    ebiten.NewImageFromImage(img),
			width:  float64(cfg.Width),
			height: float64(cfg.Height),
		})
	}
	frames["environment"] = frms

	return frames, err
}

func NewHeroKnight() *HeroKnight {
	baseSpeedRun := 1.0
	baseSpeedJump := 6.0
	frames, _ := GetFrames()
	return &HeroKnight{
		frames:           frames,
		X:                50,
		Status:           statusIdle,
		Side:             1.0,
		attackType:       attackType1,
		speedRun:         baseSpeedRun,
		maxSpeedRun:      2.5,
		baseSpeedRun:     baseSpeedRun,
		accelerationRun:  0.03,
		health:           100,
		stamina:          100,
		speedJump:        baseSpeedJump,
		baseSpeedJump:    baseSpeedJump,
		decelerationJump: 0.2,
		direction:        directionRight,
		speedRoll:        5.0,
	}
}

func (hk *HeroKnight) Death() {
	if hk.health <= 0 {
		hk.isDead = true
	}
}

func (hk *HeroKnight) Attack() {
	if !hk.IsAttacking && !hk.isRolling && !hk.isDead {
		hk.IsAttacking = true
		switch hk.attackType {
		case attackType1:
			hk.attackType = attackType2
		case attackType2:
			hk.attackType = attackType3
		case attackType3:
			hk.attackType = attackType1
		}
	}
}

func (hk *HeroKnight) Jump() {
	if !hk.isJumping && !hk.isRolling && !hk.IsAttacking {
		hk.isJumping = true
	}
}

func (hk *HeroKnight) Roll() {
	if !hk.isRolling && !hk.IsAttacking && !hk.isJumping {
		hk.isRolling = true
	}
}

func (hk *HeroKnight) Run() {
	hk.isRunning = true
}

func (hk *HeroKnight) Stop() {
	hk.isRunning = false
}

func (hk *HeroKnight) Update() error {
	hk.Death()
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		hk.health = 0
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
		hk.direction = directionLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		hk.direction = directionRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
		hk.Run()
	} else {
		hk.Stop()
	}

	switch {
	case hk.isDead:
		hk.Status = statusDeath
	case hk.IsAttacking:
		switch hk.attackType {
		case attackType1:
			hk.Status = statusAttack1
		case attackType2:
			hk.Status = statusAttack2
		case attackType3:
			hk.Status = statusAttack3
		}
	case hk.isRolling:
		hk.Status = statusRoll
	case hk.isJumping:
		hk.Status = statusJump
	case hk.isRunning:
		hk.Status = statusRun
		if hk.Frame == hk.LastFrame {
			hk.Frame = 0
		}
	default:
		hk.Status = statusIdle
		if hk.Frame == hk.LastFrame {
			hk.Frame = 0
		}
	}

	if (hk.isJumping && !hk.isDead) || (hk.isDead && hk.y > 0) {
		hk.y += hk.speedJump
		hk.speedJump -= hk.decelerationJump
		if hk.y <= 0 {
			hk.speedJump = hk.baseSpeedJump
			hk.y = 0
			hk.isJumping = false
		}
	}

	if hk.isRolling {
		hk.X += hk.speedRoll * hk.Side
	}

	if hk.isRunning && !hk.isDead && !hk.IsAttacking && !hk.isRolling {
		switch hk.direction {
		case directionLeft:
			hk.Side = -1.0
		case directionRight:
			hk.Side = 1.0
		}
		hk.X += hk.speedRun * hk.Side
		if hk.speedRun < hk.maxSpeedRun {
			hk.speedRun += hk.accelerationRun
		} else {
			hk.speedRun = hk.maxSpeedRun
		}
		hk.speedRun += hk.accelerationRun
	} else {
		hk.speedRun = hk.baseSpeedRun
	}

	switch {
	case hk.Status != hk.prevStatus:
		hk.LastFrame = statusFrames[hk.Status].framesNumber*statusFrames[hk.Status].frameDuration - 1
		hk.Frame = 0
	case hk.Frame < hk.LastFrame:
		hk.Frame++
	}

	if hk.Frame == hk.LastFrame {
		if hk.IsAttacking {
			hk.Frame = 0
			hk.IsAttacking = false
		}
		if hk.isRolling {
			hk.Frame = 0
			hk.isRolling = false
		}
	}

	hk.prevStatus = hk.Status
	return nil
}

func (hk *HeroKnight) Draw(screen *ebiten.Image) {
	tileSize := 32
	for i := 0; i < 20; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileSize*i), float64(tileSize)*9)
		screen.DrawImage(hk.frames["environment"][1].img, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(hk.Side, 1.0)
	if hk.Side < 0 {
		op.GeoM.Translate(hk.frames[hk.Status][hk.Frame/statusFrames[hk.Status].frameDuration].width, 0.0)
	}
	op.GeoM.Translate(hk.X, float64(tileSize)*9-hk.frames[hk.Status][hk.Frame/statusFrames[hk.Status].frameDuration].height-hk.y)
	screen.DrawImage(hk.frames[hk.Status][hk.Frame/statusFrames[hk.Status].frameDuration].img, op)
	w := hk.frames[hk.Status][hk.Frame/statusFrames[hk.Status].frameDuration].width - 35
	ebitenutil.DrawRect(screen, hk.X, float64(tileSize)*9-hk.frames[hk.Status][hk.Frame/statusFrames[hk.Status].frameDuration].height-hk.y, w, hk.frames[hk.Status][hk.Frame/statusFrames[hk.Status].frameDuration].height, color.RGBA{0, 0, 255, 20})
}
