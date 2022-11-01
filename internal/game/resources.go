package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	StatusAttack1       string = "Attack1"
	StatusAttack2       string = "Attack2"
	StatusAttack3       string = "Attack3"
	StatusBlock         string = "Block"
	StatusBlockIdle     string = "BlockIdle"
	StatusBlockNoEffect string = "BlockNoEffect"
	StatusDeath         string = "Death"
	StatusDeathNoBlood  string = "DeathNoBlood"
	StatusFall          string = "Fall"
	StatusHurt          string = "Hurt"
	StatusIdle          string = "Idle"
	StatusJump          string = "Jump"
	StatusLedgeGrab     string = "LedgeGrab"
	StatusRoll          string = "Roll"
	StatusRun           string = "Run"
	StatusWallSide      string = "WallSide"
	StatusAttack        string = "Attack"
	StatusCombatIdle    string = "CombatIdle"
	StatusRecover       string = "Recover"

	TileSize int = 32
)

type StatusData struct {
	FramesNumber  int
	FrameDuration int
}

type Frame struct {
	Img    *ebiten.Image
	Width  float64
	Height float64
}

var StatusFramesHeroKnight map[string]StatusData = map[string]StatusData{
	StatusAttack1: StatusData{
		FramesNumber:  6,
		FrameDuration: 4,
	},
	StatusAttack2: StatusData{
		FramesNumber:  6,
		FrameDuration: 4,
	},
	StatusAttack3: StatusData{
		FramesNumber:  8,
		FrameDuration: 4,
	},
	StatusBlock: StatusData{
		FramesNumber:  5,
		FrameDuration: 7,
	},
	StatusBlockIdle: StatusData{
		FramesNumber:  8,
		FrameDuration: 7,
	},
	StatusBlockNoEffect: StatusData{
		FramesNumber:  5,
		FrameDuration: 7,
	},
	StatusDeath: StatusData{
		FramesNumber:  10,
		FrameDuration: 7,
	},
	StatusDeathNoBlood: StatusData{
		FramesNumber:  10,
		FrameDuration: 7,
	},
	StatusFall: StatusData{
		FramesNumber:  4,
		FrameDuration: 7,
	},
	StatusHurt: StatusData{
		FramesNumber:  3,
		FrameDuration: 7,
	},
	StatusIdle: StatusData{
		FramesNumber:  8,
		FrameDuration: 7,
	},
	StatusJump: StatusData{
		FramesNumber:  3,
		FrameDuration: 2,
	},
	StatusLedgeGrab: StatusData{
		FramesNumber:  5,
		FrameDuration: 7,
	},
	StatusRoll: StatusData{
		FramesNumber:  9,
		FrameDuration: 2,
	},
	StatusRun: StatusData{
		FramesNumber:  10,
		FrameDuration: 7,
	},
	StatusWallSide: StatusData{
		FramesNumber:  5,
		FrameDuration: 7,
	},
}

var StatusFramesLightBandit map[string]StatusData = map[string]StatusData{
	StatusAttack: StatusData{
		FramesNumber:  8,
		FrameDuration: 4,
	},
	StatusCombatIdle: StatusData{
		FramesNumber:  4,
		FrameDuration: 7,
	},
	StatusDeath: StatusData{
		FramesNumber:  1,
		FrameDuration: 4,
	},
	StatusHurt: StatusData{
		FramesNumber:  2,
		FrameDuration: 7,
	},
	StatusIdle: StatusData{
		FramesNumber:  4,
		FrameDuration: 7,
	},
	StatusJump: StatusData{
		FramesNumber:  1,
		FrameDuration: 7,
	},
	StatusRecover: StatusData{
		FramesNumber:  8,
		FrameDuration: 7,
	},
	StatusRun: StatusData{
		FramesNumber:  8,
		FrameDuration: 7,
	},
}

const (
	DirectionLeft uint8 = iota + 1
	DirectionRight
)

type atype uint8

const (
	AttackType1 atype = iota + 1
	AttackType2
	AttackType3
)
