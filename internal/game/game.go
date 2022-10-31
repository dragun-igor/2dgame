package game

import (
	"fmt"
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

type Unit struct {
	status           string
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
	speedJump        float64
	baseSpeedJump    float64
	decelerationJump float64
	isAttacking      bool
	isRunning        bool
	isJumping        bool
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
	}
}

func NewGame() (*Game, error) {
	frames, err := GetFrames()
	return &Game{
		frames: frames,
		unit:   NewUnit(),
	}, err
}

func (u *Unit) RestoreStamina() {
	if u.stamina < 100 && u.status == statusIdle {
		u.stamina++
	}
}

func (u *Unit) WasteStamina(num int) {
	if num >= u.stamina {
		u.stamina = 0
	} else {
		u.stamina -= num
	}
}

func (u *Unit) Idle() {
	if !u.isAttacking && !u.isRunning && u.isJumping {
		u.status = statusIdle
	}
}

func (u *Unit) Hurt() {
	u.status = statusHurt
	u.frame = 0
	u.health -= 20
}

func (u *Unit) Run() {
	if !u.isAttacking {
		u.status = statusRun
	}
	u.isRunning = true
}

func (u *Unit) Death() {
	u.status = statusDeath
	u.frame = 0
}

func (u *Unit) Roll() {
	u.status = statusRoll
	u.frame = 0
}

func (u *Unit) Jump() {
	u.status = statusJump
	u.frame = 0
	u.isJumping = true
}

func (u *Unit) Block() {
	u.status = statusBlockIdle
	u.frame = 0
}

func (u *Unit) Attack() {
	if !u.isAttacking {
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
		u.frame = 0
		u.isAttacking = true
	}
}

// func (u *Unit) StatusSwitch(framesDrop int) {
// 	u.lastFrame = statusFrames[u.status].frameDuration * (statusFrames[u.status].framesNumber - 1)
// 	switch u.status {
// 	case statusHurt:
// 		if u.frame >= statusFrames[statusHurt].frameDuration*2 {
// 			u.status = statusIdle
// 			u.frame = 0
// 		}
// 	case statusJump:
// 		u.y += u.speedJump
// 		u.speedJump -= u.decelerationJump
// 		if u.frame == u.lastFrame && u.y <= 0 {
// 			u.y = 0
// 			u.status = statusIdle
// 			u.speedJump = u.baseSpeedJump
// 		}

// 	case statusAttack1, statusAttack2:
// 		u.x += 0.1 * u.side
// 		if u.frame >= 5*statusFrames[statusAttack1].frameDuration {
// 			u.status = statusIdle
// 			u.frame = 0
// 		}
// 	case statusAttack3:
// 		u.x += 0.1 * u.side
// 		if u.frame >= 7*statusFrames[statusAttack3].frameDuration {
// 			u.status = statusIdle
// 			u.frame = 0
// 		}
// 	case statusIdle, statusBlockIdle:
// 		if u.frame >= 7*statusFrames[statusIdle].frameDuration {
// 			u.frame = 0
// 		}
// 	case statusRoll:
// 		u.x += 3.0 * u.side
// 		if u.frame >= 8*statusFrames[statusRoll].frameDuration {
// 			u.status = statusIdle
// 			u.frame = 0
// 		}
// 	case statusRun:
// 		if u.frame >= 9*statusFrames[statusRun].frameDuration {
// 			u.frame = 0
// 		}
// 	case statusDeath:
// 		if u.frame >= 9*statusFrames[statusDeath].frameDuration {
// 			// flag = true
// 		}
// 	}
// 	if u.frame < u.lastFrame {
// 		u.frame++
// 	}
// }

func (u *Unit) StatusSwitch() {
	u.lastFrame = statusFrames[u.status].frameDuration * (statusFrames[u.status].framesNumber - 1)

	if u.isAttacking && u.frame >= u.lastFrame {
		u.frame = 0
		u.isAttacking = false
	}
	if (u.isRunning && !u.isAttacking) || (u.isRunning && u.isAttacking && u.isJumping) {
		u.x += u.speedRun * u.side
		if u.speedRun <= u.maxSpeedRun {
			u.speedRun += u.accelerationRun
		} else {
			u.speedRun = u.maxSpeedRun
		}
		if u.status == statusRun && u.frame >= u.lastFrame {
			u.frame = 0
		}
	} else {
		u.speedRun = u.baseSpeedRun
	}
	if u.isJumping {
		u.y += u.speedJump
		u.speedJump -= u.decelerationJump
		if u.y <= 0 {
			u.y = 0
			u.isJumping = false
		}
	} else {
		u.speedJump = u.baseSpeedJump
	}

	if u.isAttacking && u.frame >= u.lastFrame {
		u.frame = 0
		u.isAttacking = false
	}
	if !u.isAttacking && !u.isJumping && !u.isRunning {
		if u.frame >= u.lastFrame {
			u.frame = 0
		}
	}
	if u.frame < u.lastFrame {
		u.frame++
	}
}

func (g *Game) Update() error {
	if g.unit.health > 0 {
		if inpututil.IsKeyJustPressed(ebiten.KeyA) {
			g.unit.side = -1.0
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.unit.side = 1.0
		}
		g.unit.isRunning = ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD)
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) && !g.unit.isJumping && !g.unit.isAttacking {
			g.unit.frame = 0
			g.unit.isJumping = true
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyE) && !g.unit.isAttacking {
			g.unit.frame = 0
			g.unit.isAttacking = true
			switch g.unit.attackType {
			case attackType1:
				g.unit.attackType = attackType2
			case attackType2:
				g.unit.attackType = attackType3
			case attackType3:
				g.unit.attackType = attackType1
			}
		}
		switch {
		case g.unit.isAttacking:
			switch g.unit.attackType {
			case attackType1:
				g.unit.status = statusAttack1
			case attackType2:
				g.unit.status = statusAttack2
			case attackType3:
				g.unit.status = statusAttack3
			default:
				g.unit.status = statusAttack1
			}
		case g.unit.isJumping:
			g.unit.status = statusJump
		case g.unit.isRunning:
			g.unit.status = statusRun
		default:
			g.unit.status = statusIdle
		}
	} else {
		g.unit.Death()
	}
	g.unit.StatusSwitch()
	// g.unit.RestoreStamina()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(*g.unit)
		}
	}()
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
