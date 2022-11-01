package lightbandit

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
	statusAttack     string = "Attack"
	statusCombatIdle string = "CombatIdle"
	statusDeath      string = "Death"
	statusHurt       string = "Hurt"
	statusIdle       string = "Idle"
	statusJump       string = "Jump"
	statusRecover    string = "Recover"
	statusRun        string = "Run"
)

type statusData struct {
	framesNumber  int
	frameDuration int
}

var statusFrames map[string]statusData = map[string]statusData{
	statusAttack: statusData{
		framesNumber:  8,
		frameDuration: 4,
	},
	statusCombatIdle: statusData{
		framesNumber:  4,
		frameDuration: 7,
	},
	statusDeath: statusData{
		framesNumber:  1,
		frameDuration: 4,
	},
	statusHurt: statusData{
		framesNumber:  2,
		frameDuration: 7,
	},
	statusIdle: statusData{
		framesNumber:  4,
		frameDuration: 7,
	},
	statusJump: statusData{
		framesNumber:  1,
		frameDuration: 7,
	},
	statusRecover: statusData{
		framesNumber:  8,
		frameDuration: 7,
	},
	statusRun: statusData{
		framesNumber:  8,
		frameDuration: 7,
	},
}

type Frame struct {
	img    *ebiten.Image
	width  float64
	height float64
}

const (
	directionLeft uint8 = iota + 1
	directionRight
)

type LightBandit struct {
	frames           map[string][]Frame
	status           string
	prevStatus       string
	X                float64
	y                float64
	side             float64
	speedRun         float64
	maxSpeedRun      float64
	baseSpeedRun     float64
	accelerationRun  float64
	speedJump        float64
	baseSpeedJump    float64
	decelerationJump float64
	frame            int
	lastFrame        int
	health           int
	direction        uint8
	isDead           bool
	isHurted         bool
	isAttacking      bool
	isJumping        bool
	isRunning        bool
}

func NewLightBandit() *LightBandit {
	frames, _ := GetFrames()
	return &LightBandit{
		frames:           frames,
		X:                550,
		y:                0,
		status:           statusCombatIdle,
		prevStatus:       statusCombatIdle,
		side:             1.0,
		speedRun:         1.0,
		maxSpeedRun:      2.0,
		baseSpeedRun:     1.0,
		accelerationRun:  0.2,
		speedJump:        6.0,
		baseSpeedJump:    6.0,
		decelerationJump: 0.2,
		health:           100,
		lastFrame:        (statusFrames[statusCombatIdle].framesNumber - 1) * statusFrames[statusCombatIdle].frameDuration,
		direction:        directionLeft,
	}
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
			file, err = os.Open("_assets/LightBandit/" + status + "/LightBandit_" + status + "_" + strconv.Itoa(i) + ".png")
			if err != nil {
				break
			}
			img, err = png.Decode(file)
			if err != nil {
				break
			}
			file.Close()
			file, err = os.Open("_assets/LightBandit/" + status + "/LightBandit_" + status + "_" + strconv.Itoa(i) + ".png")
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
	return frames, err
}

func (lb *LightBandit) Death() {
	if !lb.isDead {
		lb.isDead = true
	}
}

func (lb *LightBandit) Hurt() {
	if !lb.isHurted {
		lb.isHurted = true
		lb.health -= 2
	}
}

func (lb *LightBandit) Update() error {
	if lb.health <= 0 {
		lb.Death()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		lb.Hurt()
	}

	switch {
	case lb.isDead:
		lb.status = statusDeath
	case lb.isHurted:
		lb.status = statusHurt
	default:
		lb.status = statusCombatIdle
	}

	switch {
	case lb.status != lb.prevStatus:
		lb.lastFrame = statusFrames[lb.status].framesNumber*statusFrames[lb.status].frameDuration - 1
		lb.frame = 0
	case lb.frame < lb.lastFrame:
		lb.frame++
	}
	if lb.frame == lb.lastFrame {
		lb.isHurted = false
		lb.frame = 0
	}
	lb.prevStatus = lb.status
	return nil
}

func (lb *LightBandit) Draw(screen *ebiten.Image) {
	tileSize := 32
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(lb.side, 1.0)
	if lb.side < 0 {
		op.GeoM.Translate(lb.frames[lb.status][lb.frame/statusFrames[lb.status].frameDuration].width, 0.0)
	}
	op.GeoM.Translate(lb.X, float64(tileSize)*9-lb.frames[lb.status][lb.frame/statusFrames[lb.status].frameDuration].height-lb.y)
	screen.DrawImage(lb.frames[lb.status][lb.frame/statusFrames[lb.status].frameDuration].img, op)
	w := lb.frames[lb.status][lb.frame/statusFrames[lb.status].frameDuration].width
	ebitenutil.DrawRect(screen, lb.X, float64(tileSize)*9-lb.frames[lb.status][lb.frame/statusFrames[lb.status].frameDuration].height-lb.y, w, lb.frames[lb.status][lb.frame/statusFrames[lb.status].frameDuration].height, color.RGBA{0, 0, 255, 20})

}
