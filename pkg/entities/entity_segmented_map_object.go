package entities

import (
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// MapObject Entity
type SegmentedMapObjectEntity2D struct {
	// Required fields
	proto.BaseEntity2D

	// Custom Fields
    sprite *rl.Image
    segment_region rl.Rectangle
}

func NewSegmentedMapObjectEntity2D(
    position rl.Vector2,
    sprite *rl.Image,
    //rotation float32, // should not be needed i think
) *SegmentedMapObjectEntity2D {
	new_ent := &SegmentedMapObjectEntity2D{
		BaseEntity2D: proto.BaseEntity2D{Transform: proto.Transform2D{Position: position, Rotation: 0, Scale: rl.Vector2One()}},

        sprite: sprite,
	}
	return new_ent
}

func (ent *SegmentedMapObjectEntity2D) Init() {
    ent.BaseEntity2D.Init()

    heightThreshold := int32(5.0)  // Define your threshold here
    totalWidth := int32(ent.sprite.Width)
    texture := rl.LoadTextureFromImage(ent.sprite)

    // Get lowest pixel height for every column
    lowestPixels := make([]int32, ent.sprite.Width)
    for x := int32(0); x < int32(ent.sprite.Width); x++ {
        lowestPixels[x] = int32(findLowestNonTransparentPixel(ent.sprite, x))
    }

    createSegment := func(startX, endX int) {
        srcRect := rl.NewRectangle(
            float32(startX),
            0,
            float32(endX-startX),
            float32(ent.sprite.Height),
        )

        lowestY := lowestPixels[startX]

        if lowestY >= 0 {
            drawIndexHeight := int32(ent.GetPosition().Y) + lowestY
            segment := NewMapObjectSegmentEntity2D(
                rl.NewVector2(float32(startX), 0),
                texture,
                srcRect,
                drawIndexHeight,
            )
            segment.Init()
            gem.AddEntity(ent, segment)
        }
    }

    // Dynamic segmentation based on lowest pixel height changes
    startX := 0
    for x := 1; x < int(totalWidth); x++ {
        if util.Abs(lowestPixels[x]-lowestPixels[startX]) > heightThreshold {
            createSegment(startX, x)
            startX = x
        }
    }
    // Create the final segment
    createSegment(startX, int(totalWidth))
}

func findLowestNonTransparentPixel(image *rl.Image, x int32) int32 {
    lowestY := int32(-1)

    for y := int32(0); y < int32(image.Height); y++ {
        color := rl.GetImageColor(*image, x, y)
        if color.A > 0 && (lowestY == -1 || y > int32(lowestY)) {
            lowestY = int32(y)
        }
    }

    return lowestY
}


func (ent *SegmentedMapObjectEntity2D) Deinit() {
	// Required de-initialization
	ent.BaseEntity2D.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *SegmentedMapObjectEntity2D) Update() {
	// Required update
	ent.BaseEntity2D.Update()

	// Update logic for the entity
	// ...
}

func (ent *SegmentedMapObjectEntity2D) Draw() {
    // drawing is done in segments
}

func (ent *SegmentedMapObjectEntity2D) GetName() string {
    return "SegmentedMapObject"
}

// GetDrawIndex returns the draw index of this entity. Entities with a higher
// index are drawn in front of entities with a lower index.
func (ent *SegmentedMapObjectEntity2D) GetDrawIndex() int32 {
	return int32(ent.GetPosition().Y)+13
}
