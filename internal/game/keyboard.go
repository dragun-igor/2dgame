package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	KeyAttack uint8 = iota + 1
	KeyRunLeft
	KeyRunRight
	KeyJump
	KeyRoll
)

type Keyboard map[uint8]*Key

type Key struct {
	IsKeyPressed     bool
	IsKeyJustPressed bool
	Key1             ebiten.Key
	Key2             ebiten.Key
	KeyEmulation     bool
	prevStatus       bool
}

func (k *Key) update() {
	k.IsKeyPressed = ebiten.IsKeyPressed(k.Key1) || ebiten.IsKeyPressed(k.Key2) || k.KeyEmulation
	k.IsKeyJustPressed = k.IsKeyPressed && !k.prevStatus
	k.prevStatus = k.IsKeyPressed
}

func NewKey(key1, key2 ebiten.Key) *Key {
	return &Key{
		Key1: key1,
		Key2: key2,
	}
}

func NewDefaultKeyboard() Keyboard {
	return Keyboard{
		KeyAttack:   NewKey(ebiten.KeyE, -1),
		KeyRunLeft:  NewKey(ebiten.KeyA, -1),
		KeyRunRight: NewKey(ebiten.KeyD, -1),
		KeyJump:     NewKey(ebiten.KeySpace, -1),
		KeyRoll:     NewKey(ebiten.KeyControlLeft, -1),
	}
}

func NewEmulationKeyboard() Keyboard {
	return Keyboard{
		KeyAttack:   NewKey(-1, -1),
		KeyRunLeft:  NewKey(-1, -1),
		KeyRunRight: NewKey(-1, -1),
		KeyJump:     NewKey(-1, -1),
		KeyRoll:     NewKey(-1, -1),
	}
}

func (k Keyboard) Update() {
	for str := range k {
		k[str].update()
	}
}
