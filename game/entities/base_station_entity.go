package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/fw/core/math"
	"gorl/game/code/astar"
	"gorl/game/code/colorscheme"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that BaseStationEntity implements IEntity.
var _ entities.IEntity = &BaseStationEntity{}

// BaseStation Entity
type BaseStationEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	constructedMapData []color.RGBA
	constructedMap     rl.Texture2D
	gridMap            [][]bool
	gridTileSize       int
}

// NewBaseStationEntity creates a new instance of the BaseStationEntity.
func NewBaseStationEntity(position rl.Vector2) *BaseStationEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &BaseStationEntity{
		Entity: entities.NewEntity("BaseStationEntity", position, 0, rl.Vector2One()),

		gridTileSize: 10,
	}

	// initialize the quantized grid map with false values
	new_ent.gridMap = make([][]bool, 1920/new_ent.gridTileSize, 1920/new_ent.gridTileSize)
	for idx := range new_ent.gridMap {
		new_ent.gridMap[idx] = make([]bool, 1080/new_ent.gridTileSize, 1080/new_ent.gridTileSize)
	}

	// fill the image with transparent black pixels
	imgData := make([]byte, 1920*1080*4)
	emptyImg := rl.NewImage(imgData, 1920, 1080, 1, rl.UncompressedR8g8b8a8)
	new_ent.constructedMap = rl.LoadTextureFromImage(emptyImg)
	new_ent.constructedMapData = make([]color.RGBA, 1920*1080)

	return new_ent
}

func (ent *BaseStationEntity) DiscoverPixel(x, y int) {
	ent.constructedMapData[y*1920+x] = colorscheme.Colorscheme.Color04.ToRGBA()
	gridX := x / ent.gridTileSize
	gridY := y / ent.gridTileSize
	ent.gridMap[gridX][gridY] = true
}

// GetPath returns a path from one point to another through the base stations
// grid map.
func (ent *BaseStationEntity) GetPath(from, to rl.Vector2) []rl.Vector2 {

	fromGrid := rl.Vector2Divide(from, rl.NewVector2(float32(ent.gridTileSize), float32(ent.gridTileSize)))
	toGrid := rl.Vector2Divide(to, rl.NewVector2(float32(ent.gridTileSize), float32(ent.gridTileSize)))

	path := astar.AstarPath(
		math.Vector2IntFromRl(fromGrid),
		math.Vector2IntFromRl(toGrid),
		ent.gridMap,
	)

	floatWorldSpacePath := make([]rl.Vector2, 0, len(path))
	for _, node := range path {
		floatWorldSpacePath = append(
			floatWorldSpacePath,
			rl.NewVector2(float32(node.X*ent.gridTileSize), float32(node.Y*ent.gridTileSize)),
		)
	}

	return floatWorldSpacePath
}

func (ent *BaseStationEntity) Init() {
	// Initialization logic for the entity
	// ...
}

func (ent *BaseStationEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *BaseStationEntity) Update() {
	// Make the texture match the map data
	rl.UpdateTexture(ent.constructedMap, ent.constructedMapData)

}

func (ent *BaseStationEntity) Draw() {

	// Draw the grid map
	gridX := len(ent.gridMap)
	gridY := len(ent.gridMap[0])
	for x := 0; x < gridX; x++ {
		for y := 0; y < gridY; y++ {
			if ent.gridMap[x][y] {
				rl.DrawCircleV(
					rl.NewVector2(float32(x*ent.gridTileSize), float32(y*ent.gridTileSize)),
					float32(ent.gridTileSize/2),
					rl.Green,
				)
			}
		}
	}

	// Draw the constructed map
	rl.DrawTexture(ent.constructedMap, 0, 0, rl.White)

	rl.DrawCircleV(ent.GetPosition(), 20, colorscheme.Colorscheme.Color01.ToRGBA())

}

func (ent *BaseStationEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	return true
}
