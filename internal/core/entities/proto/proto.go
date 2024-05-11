package proto

import (
	"gorl/internal/ai"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// IEntity is an interface that every entity in the game should implement
type IEntity interface {
	Init()
	Deinit()
	Update()
	FixedUpdate() // only called when phyics enabled
	Draw()

	// new

	// OnChildAdded is called every time an entity is added to the Gem with this Entity as the parent.
	OnChildAdded(child IEntity)
	// OnChildRemoved is called every time a child of this entity is removed from the Gem.
	OnChildRemoved(child IEntity)

	// old
	DrawGUI()
	AddChild(child IEntity)
	RemoveChild(child IEntity)
	GetChildren() []IEntity
	GetParent() IEntity
	SetParent(parent IEntity)
	GetDrawIndex() int32
	SetDrawIndex(index int32)
	GetName() string
}

// -------------------
//
//	ENTITY 2D
//
// -------------------
type Transform2D struct {
	Position rl.Vector2
	Rotation float32
	Scale    rl.Vector2
}

func BaseTransform2D() Transform2D {
	return Transform2D{
		Position: rl.Vector2{X: 0, Y: 0},
		Rotation: 0,
		Scale:    rl.Vector2{X: 1, Y: 1},
	}
}

type IEntity2D interface {
	IEntity
	GetPosition() rl.Vector2
	SetPosition(new_position rl.Vector2)
	GetScale() rl.Vector2
	SetScale(new_size rl.Vector2)
	GetRotation() float32
}

type Entity2DPlayer interface {
	IEntity2D
	SendMessage(message string, sender IEntity)
}

// -------------------
//
//	ENTITY 2D AI
//
// -------------------
type Entity2DAI interface {
	IEntity2D
	ai.AiControllable
}
