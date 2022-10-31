package game

import (
	"image"
	"image/png"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
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

type Game struct {
	frames map[string][]Frame
	unit   *Unit
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

type Unit struct {
	status           string
	prevStatus       string
	x                float64
	y                float64
	side             float64
	frame            int
	lastFrame        int
	attackType       atype
	speedRun         float64
	maxSpeedRun      float64
	baseSpeedRun     float64
	accelerationRun  float64
	health           int
	stamina          int
	direction        uint8
	speedJump        float64
	baseSpeedJump    float64
	decelerationJump float64
	isDead           bool
	isHurted         bool
	isAttacking      bool
	isBlocking       bool
	isJumping        bool
	isRolling        bool
	isRunning        bool
	priority         int
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
			file, err = os.Open("_assets/" + status + "/HeroKnight_" + status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				break
			}
			file.Close()
			file, err = os.Open("_assets/" + status + "/HeroKnight_" + status + "_" + strconv.Itoa(i) + ".png")
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

func NewUnit() *Unit {
	baseSpeedRun := 1.0
	baseSpeedJump := 6.0
	return &Unit{
		status:           statusIdle,
		side:             1.0,
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
	}
}

func NewGame() (*Game, error) {
	frames, err := GetFrames()
	return &Game{
		frames: frames,
		unit:   NewUnit(),
	}, err
}

func (u *Unit) Death() {
	if u.health <= 0 {
		u.isDead = true
	}
}

func (u *Unit) Attack() {
	if !u.isAttacking {
		u.isAttacking = true
		switch u.attackType {
		case attackType1:
			u.status = statusAttack1
			u.attackType = attackType2
		case attackType2:
			u.status = statusAttack2
			u.attackType = attackType3
		case attackType3:
			u.status = statusAttack3
			u.attackType = attackType1
		}
	}
}

func (u *Unit) Jump() {
	if !u.isJumping {
		u.isJumping = true
	}
}

func (u *Unit) Run() {
	u.isRunning = true
}

func (u *Unit) Stop() {
	u.isRunning = false
	u.speedRun = u.baseSpeedRun
}

func (g *Game) Update() error {
	u := g.unit
	u.Death()
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		u.health = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		u.Attack()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		u.Jump()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		u.direction = directionLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		u.direction = directionRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
		u.Run()
	} else {
		u.Stop()
	}
	u.lastFrame = (statusFrames[u.status].framesNumber - 1) * statusFrames[u.status].frameDuration

	switch {
	case u.isDead:
		u.status = statusDeath
	case u.isAttacking:
		if u.frame == u.lastFrame {
			u.isAttacking = false
		}
	case u.isJumping:
		u.status = statusJump
		u.y += u.speedJump
		u.speedJump -= u.decelerationJump
		if u.y <= 0 {
			u.y = 0
			u.speedJump = u.baseSpeedJump
			u.isJumping = false
		}
	case u.isRunning:
		switch u.direction {
		case directionLeft:
			u.side = -1.0
		case directionRight:
			u.side = 1.0
		}
		u.status = statusRun
		u.x += u.speedRun * u.side
		if u.speedRun < u.maxSpeedRun {
			u.speedRun += u.accelerationRun
		} else {
			u.speedRun = u.maxSpeedRun
		}
		u.speedRun += u.accelerationRun
		if u.frame == u.lastFrame {
			u.frame = 0
		}
	default:
		u.status = statusIdle
		if u.frame == u.lastFrame {
			u.frame = 0
		}
	}
	if u.frame < u.lastFrame {
		u.frame++
	}
	if u.status != u.prevStatus {
		u.frame = 0
	}
	u.prevStatus = u.status
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	tileSize := 32
	for i := 0; i < 50; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(tileSize*i), float64(tileSize)*5+g.frames[g.unit.status][g.unit.frame/statusFrames[g.unit.status].frameDuration].height)
		screen.DrawImage(g.frames["environment"][1].img, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.unit.side, 1.0)
	if g.unit.side < 0 {
		op.GeoM.Translate(g.frames[g.unit.status][g.unit.frame/7].width, 0.0)
	}
	op.GeoM.Translate(g.unit.x, float64(tileSize)*5-g.unit.y)
	screen.DrawImage(g.frames[g.unit.status][g.unit.frame/statusFrames[g.unit.status].frameDuration].img, op)
}
