package robot

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Robot struct {
	position  rl.Vector2
	direction rl.Vector2
	Color     rl.Color // NOTE: REMOVE WHEN ICON IS USED
	Speed     float32
}

// TODO: REMOVE WHOLE FILE

func NewRobot(position rl.Vector2) *Robot {
	return &Robot{
		position:  position,
		direction: rl.Vector2Zero(),
		Color:     rl.Green,
		Speed:     40,
	}

}

// ===== GETTER/SETTER =====//
func (rb *Robot) GetPosition() rl.Vector2 {
	return rb.position
}
func (rb *Robot) SetPosition(newPos rl.Vector2) {
	rb.position = newPos
}
func (rb *Robot) GetDirection() rl.Vector2 {
	return rb.direction
}
func (rb *Robot) SetDirection(newDir rl.Vector2) {
	rb.direction = newDir
}
