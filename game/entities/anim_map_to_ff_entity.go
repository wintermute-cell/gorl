package entities

import (
	"fmt"
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that AnimMapToFfEntity implements IEntity.
var _ entities.IEntity = &AnimMapToFfEntity{}

// AnimMapToFf Entity
type AnimMapToFfEntity struct {
	*entities.Entity // Required!
	MapPng           rl.Texture2D
	MapImage         rl.Image // used to get pixeldata from the texture
	TileSize         int
	sclSec           float32
	faderGrid        int
	faderNodes       int
	faderEdges       int
}

// NewAnimMapToFfEntity creates a new instance of the AnimMapToFfEntity.
func NewAnimMapToFfEntity() *AnimMapToFfEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &AnimMapToFfEntity{
		Entity:   entities.NewEntity("AnimMapToFfEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		MapPng:   rl.LoadTexture("./map_thresh.png"),
		TileSize: 40, // NOTE: not very clean to use this here
	}
	return new_ent
}

func (ent *AnimMapToFfEntity) Init() {
	// Initialization logic for the entity
	// ...
	ent.MapImage = *rl.LoadImageFromTexture(ent.MapPng)
	ent.faderGrid = 255
	ent.faderNodes = 0
	ent.faderEdges = 0
	ent.CalculateGridGraphImage()
	// TODO: set this to the grid graph
}

func (ent *AnimMapToFfEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *AnimMapToFfEntity) Update() {
	// update draw position of the grid (sclSec) and faders
	ent.sclSec += rl.GetFrameTime() * 1000
	if ent.sclSec > 6000 && ent.faderNodes < 255 {
		if ent.faderNodes < 255 {

			ent.faderNodes++
		}
		if ent.faderGrid > 0 {
			ent.faderGrid--
		}
	}
	if ent.sclSec > 10000 && ent.faderEdges < 255 {
		ent.faderEdges += 5
	}
	// TODO: IDEE: timer und flag abfragen, ob genug zeit fuer die animation vergangen ist
	// und ob schon die karte erstell wurde (CalculateGridGraphImage gecallt wurde),
	// dann erstellen, und displayen, machste morgen brudi :)
}

func (ent *AnimMapToFfEntity) Draw() {
	// Draw the map png
	rl.DrawTextureV(ent.MapPng, ent.GetPosition(), rl.White)

	// Draw the fade in of the grid, but only for a certain time period, until
	// we draw the graph visualization
	if ent.sclSec < 8000 {
		for i := range 50 {
			// top to bottom
			delay := -500
			rl.DrawLine(
				int32(i*ent.TileSize),
				int32(-i*ent.TileSize+delay),
				int32(i*ent.TileSize),
				int32(-i*ent.TileSize+int(ent.sclSec)+delay),
				rl.NewColor(255, 0, 0, uint8(ent.faderGrid)),
			)
			// left to right
			rl.DrawLine(
				int32(-i*ent.TileSize),
				int32(i*ent.TileSize),
				int32(-i*ent.TileSize+int(ent.sclSec)),
				int32(i*ent.TileSize),
				rl.NewColor(255, 0, 0, uint8(ent.faderGrid)),
			)
		}
	}

	// draw graph edges, we use -20 to start drawing the edges out of the actual frame,
	// because with an unlimited map the graph would be endless to all sides (and it is a lot simpler)
	// we draw the edges before the nodes so that they are drawn first
	// TODO: draw diagonals
	if ent.sclSec > 10000 {
		for i := range 50 {
			for j := range 50 {
				rl.DrawLine(
					int32(i*ent.TileSize-20),
					int32(j*ent.TileSize-20),
					int32(i*ent.TileSize+2000),
					int32(j*ent.TileSize-20),
					rl.NewColor(0, 0, 255, uint8(ent.faderEdges)),
				)
				rl.DrawLine(
					int32(i*ent.TileSize-20),
					int32(j*ent.TileSize-20),
					int32(i*ent.TileSize-20),
					int32(j*ent.TileSize+2000),
					rl.NewColor(0, 0, 255, uint8(ent.faderEdges)),
				)
			}
		}
	}
	// draw graph nodes
	if ent.sclSec > 6000 {
		for i := range 50 {
			for j := range 50 {
				rl.DrawCircle(
					int32(i*ent.TileSize+20),
					int32(j*ent.TileSize+20),
					10,
					rl.NewColor(255, 0, 0, uint8(ent.faderNodes)),
				)
			}
		}
	}
}

func (ent *AnimMapToFfEntity) CalculateGridGraphImage() {
	// this is a one dimensional array
	imgColors := rl.LoadImageColors(&ent.MapImage)
	gridGraphMap := make(map[rl.Vector2]rl.Color)
	// converting to 2D slice
	imgColors2D := make([][]color.RGBA, ent.MapImage.Height)
	for i := range imgColors2D {
		imgColors2D[i] = make([]color.RGBA, ent.MapImage.Width)
	}
	// add values to the 2D slice
	for x, _ := range imgColors2D {
		for y, _ := range imgColors2D[x] {
			imgColors2D[x][y] = imgColors[x*int(ent.MapImage.Width)+y]
		}
	}
	sclWidth := len(imgColors2D) / ent.TileSize
	sclHeight := len(imgColors2D[0]) / ent.TileSize
	// loop the tile size
	for x := range sclWidth {
		for y := range sclHeight {
			// now loop the inner pixels of each tile and determine the predominant color
			whiteCount := 0
			blackCount := 0
			for p := range ent.TileSize {
				for q := range ent.TileSize {
					// since R, G and B are all 255 in a white pixel, its enough to just check one value
					if imgColors2D[x*ent.TileSize+p][y*ent.TileSize+q].R == 255 {
						whiteCount++
					} else {
						blackCount++
					}
				}
			}
			if whiteCount > blackCount {
				gridGraphMap[rl.NewVector2(float32(x), float32(y))] = rl.White
			} else {
				gridGraphMap[rl.NewVector2(float32(x), float32(y))] = rl.Black
			}
		}
	}
	fmt.Println(gridGraphMap)
	// NOTE: OMG ICH GLAUB ES FUNKTIONIERT ENDLICH
}

func (ent *AnimMapToFfEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
