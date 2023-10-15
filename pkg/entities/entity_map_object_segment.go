package entities

import (
	"cowboy-gorl/pkg/entities/proto"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// MapObjectSegment Entity
type MapObjectSegmentEntity2D struct {
    proto.BaseEntity2D

    sprite rl.Texture2D
    srcRect rl.Rectangle
    drawIndexHeight int32
}

func NewMapObjectSegmentEntity2D(
    position rl.Vector2,
    sprite rl.Texture2D,
    srcRect rl.Rectangle,
    drawIndexHeight int32,
) *MapObjectSegmentEntity2D {
    return &MapObjectSegmentEntity2D{
        BaseEntity2D: proto.BaseEntity2D{
            Transform: proto.Transform2D{
                Position: position,
            },
        },
        drawIndexHeight: drawIndexHeight,
        sprite: sprite,
        srcRect: srcRect,
    }
}

func (ent *MapObjectSegmentEntity2D) Init() {
	// Required initialization
	ent.BaseEntity2D.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *MapObjectSegmentEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *MapObjectSegmentEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *MapObjectSegmentEntity2D) Draw() {
    rl.DrawTexturePro(
        ent.sprite,
        ent.srcRect,
        rl.NewRectangle(ent.GetPosition().X, ent.GetPosition().Y, ent.srcRect.Width, ent.srcRect.Height),
        rl.NewVector2(0, 0),  // Rotation origin in the middle of the height
        0,  // If you use rotations - otherwise, keep 0
        rl.White,
    )
}

func (ent *MapObjectSegmentEntity2D) GetName() string {
    return "MapObjectSegment"
}

func (ent *MapObjectSegmentEntity2D) GetDrawIndex() int32 {
    // Compute draw index, maybe just Y, maybe Y + offset based on sprite/segment
    return ent.drawIndexHeight
}
