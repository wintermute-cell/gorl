package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that AnimMapToFfEntity implements IEntity.
var _ entities.IEntity = &AnimMapToFfEntity{}

// AnimMapToFf Entity
type AnimMapToFfEntity struct {
	*entities.Entity // Required!
	MapPng           rl.Texture2D
	sclSec           float32
}

// NewAnimMapToFfEntity creates a new instance of the AnimMapToFfEntity.
func NewAnimMapToFfEntity() *AnimMapToFfEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &AnimMapToFfEntity{
		Entity: entities.NewEntity("AnimMapToFfEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		MapPng: rl.LoadTexture("./map_thresh.png"),
		// Initialize custom fields here
		// ...
	}
	return new_ent
}

func (ent *AnimMapToFfEntity) Init() {
	// Initialization logic for the entity
	// ...

}

func (ent *AnimMapToFfEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *AnimMapToFfEntity) Update() {
	// Update logic for the entity per frame
	// ...
	// TODO: Logic here to update draw positions of the grid
	ent.sclSec += rl.GetFrameTime() * 1000
}

func (ent *AnimMapToFfEntity) Draw() {
	tileSize := 40 // tilesize from grid graph entity
	// Draw the map png
	rl.DrawTextureV(ent.MapPng, ent.GetPosition(), rl.White)
	// Draw the fade in of the grid, but only for a certain time period, until
	// we draw the graph visualization
	// TODO: fade out when finished
	if ent.sclSec < 8000 {
		for i := range 50 {
			// top to bottom
			delay := -500
			rl.DrawLine(
				int32(i*tileSize),
				int32(-i*tileSize+delay),
				int32(i*tileSize),
				int32(-i*tileSize+int(ent.sclSec)+delay),
				rl.Red,
			)
			// left to right
			rl.DrawLine(
				int32(-i*tileSize),
				int32(i*tileSize),
				int32(-i*tileSize+int(ent.sclSec)),
				int32(i*tileSize),
				rl.Red,
			)
		}
	}
	// TODO: convert the raw map png to a grid graph with the correct obstacles set and draw it

	// draw graph edges, we use -20 to start drawing the edges out of the actual frame,
	// because with an unlimited map the graph would be endless to all sides (and it is a lot simpler)
	// we draw the edges before the nodes so that they are drawn first
	if ent.sclSec > 10000 {
		for i := range 50 {
			for j := range 50 {
				rl.DrawLine(
					int32(i*tileSize-20),
					int32(j*tileSize-20),
					int32(i*tileSize+2000),
					int32(j*tileSize-20),
					rl.Blue,
				)
				rl.DrawLine(
					int32(i*tileSize-20),
					int32(j*tileSize-20),
					int32(i*tileSize-20),
					int32(j*tileSize+2000),
					rl.Blue,
				)
			}
		}
	}
	// draw underlying graph
	if ent.sclSec > 6000 {
		for i := range 50 {
			for j := range 50 {
				rl.DrawCircle(
					int32(i*tileSize+20),
					int32(j*tileSize+20),
					10,
					rl.Red,
				)
			}
		}
	}
}

func (ent *AnimMapToFfEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
