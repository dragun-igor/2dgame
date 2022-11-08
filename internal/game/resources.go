package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

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
	Scale    int = 2
)

type StatusData struct {
	FramesNumber  int
	FrameDuration int
}

type Unit struct {
	ActionFrames map[string][]*ebiten.Image
	Width        float64
	Height       float64
}

const (
	MainCharacter string = "HeroKnight"
	LightBandit   string = "LightBandit"
	HeavyBandit   string = "HeavyBandit"
	Wizard        string = "Wizard"
)

var StatusFrames = map[string]map[string]StatusData{
	MainCharacter: StatusFramesHeroKnight,
	LightBandit:   StatusFramesLightBandit,
	HeavyBandit:   StatusFramesLightBandit,
	Wizard:        StatusFramesWizard,
}

var StatusFramesHeroKnight = map[string]StatusData{
	StatusAttack1: {
		FramesNumber:  6,
		FrameDuration: 4,
	},
	StatusAttack2: {
		FramesNumber:  6,
		FrameDuration: 4,
	},
	StatusAttack3: {
		FramesNumber:  8,
		FrameDuration: 4,
	},
	StatusBlock: {
		FramesNumber:  5,
		FrameDuration: 7,
	},
	StatusBlockIdle: {
		FramesNumber:  8,
		FrameDuration: 7,
	},
	StatusBlockNoEffect: {
		FramesNumber:  5,
		FrameDuration: 7,
	},
	StatusDeath: {
		FramesNumber:  10,
		FrameDuration: 7,
	},
	StatusDeathNoBlood: {
		FramesNumber:  10,
		FrameDuration: 7,
	},
	StatusFall: {
		FramesNumber:  4,
		FrameDuration: 7,
	},
	StatusHurt: {
		FramesNumber:  3,
		FrameDuration: 7,
	},
	StatusIdle: {
		FramesNumber:  8,
		FrameDuration: 7,
	},
	StatusJump: {
		FramesNumber:  3,
		FrameDuration: 2,
	},
	StatusLedgeGrab: {
		FramesNumber:  5,
		FrameDuration: 7,
	},
	StatusRoll: {
		FramesNumber:  9,
		FrameDuration: 4,
	},
	StatusRun: {
		FramesNumber:  10,
		FrameDuration: 7,
	},
	StatusWallSide: {
		FramesNumber:  5,
		FrameDuration: 7,
	},
}

var StatusFramesLightBandit = map[string]StatusData{
	StatusAttack: {
		FramesNumber:  8,
		FrameDuration: 4,
	},
	StatusCombatIdle: {
		FramesNumber:  4,
		FrameDuration: 7,
	},
	StatusDeath: {
		FramesNumber:  1,
		FrameDuration: 4,
	},
	StatusHurt: {
		FramesNumber:  2,
		FrameDuration: 7,
	},
	StatusIdle: {
		FramesNumber:  4,
		FrameDuration: 7,
	},
	StatusJump: {
		FramesNumber:  1,
		FrameDuration: 7,
	},
	StatusRecover: {
		FramesNumber:  8,
		FrameDuration: 7,
	},
	StatusRun: {
		FramesNumber:  8,
		FrameDuration: 7,
	},
}

var StatusFramesWizard = map[string]StatusData{
	StatusAttack1: {
		FramesNumber:  8,
		FrameDuration: 4,
	},
	StatusAttack2: {
		FramesNumber:  8,
		FrameDuration: 7,
	},
	StatusDeath: {
		FramesNumber:  7,
		FrameDuration: 4,
	},
	StatusFall: {
		FramesNumber:  2,
		FrameDuration: 7,
	},
	StatusHurt: {
		FramesNumber:  4,
		FrameDuration: 7,
	},
	StatusIdle: {
		FramesNumber:  6,
		FrameDuration: 7,
	},
	StatusJump: {
		FramesNumber:  2,
		FrameDuration: 7,
	},
	StatusRun: {
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
