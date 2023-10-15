package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// MapObject Entity
type MapObjectEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    sprite rl.Texture2D
}

func NewMapObjectEntity2D(
    position rl.Vector2,
    sprite rl.Texture2D,
) *MapObjectEntity2D {
	new_ent := &MapObjectEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},

        sprite: sprite,
	}
	return new_ent
}

func (ent *MapObjectEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *MapObjectEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *MapObjectEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *MapObjectEntity2D) Draw() {
    rl.DrawTexturePro(
        ent.sprite,
        rl.NewRectangle(0, 0, float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, float32(ent.sprite.Width), float32(ent.sprite.Height)),
        rl.Vector2Zero(), 0, rl.White,
        )
}

func (ent *MapObjectEntity2D) GetName() string {
    return "MapObject"
}

func (ent *MapObjectEntity2D) GetDrawIndex() int32 {
	return int32(ent.GetPosition().Y)+ent.sprite.Height
}
