package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/game/code/grid_graph"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that GridGraphEntity implements IEntity.
var _ entities.IEntity = &GridGraphEntity{}

// GridGraph Entity
type GridGraphEntity struct {
	*entities.Entity // Required!
	gg               *grid_graph.GridGraph
	TextSize         int32
}

// NewGridGraphEntity creates a new instance of the GridGraphEntity.
func NewGridGraphEntity() *GridGraphEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &GridGraphEntity{
		Entity:   entities.NewEntity("GridGraphEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		gg:       grid_graph.NewGridGraph(100, 100),
		TextSize: 20,
	}
	return new_ent
}

func (ent *GridGraphEntity) Init() {
	// Initialization logic for the entity
	// ...

	mapImage := rl.LoadImage("./map_thresh.png")
	ent.gg = grid_graph.NewGridGraph(48, 27)
	ent.gg.CalculateGridGraphFromImage(mapImage, 40)
	// ent.gg.RemoveUnreachableTiles(rl.NewVector2(10, 0))
	// ent.gg.Dijkstra(rl.NewVector2(10, 0))
}

func (ent *GridGraphEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *GridGraphEntity) Update() {
	// Update logic for the entity per frame
	// ...
}

func (ent *GridGraphEntity) Draw() {
	// Draw vertices
	rl.DrawRectangleV(
		ent.GetPosition(),
		rl.NewVector2(
			float32(ent.gg.Width)*float32(ent.gg.TileSize),
			float32(ent.gg.Height)*float32(ent.gg.TileSize),
		),
		rl.Black,
	)
	// TODO: put this in a vertex.GetColor() function
	for _, vertex := range ent.gg.VertexMap {
		sclColorVal := vertex.Distance * 20
		var vertexColor rl.Color
		if sclColorVal <= 255 {
			vertexColor = rl.NewColor(
				255-uint8(sclColorVal),
				255-uint8(sclColorVal),
				255,
				255,
			)
		} else if sclColorVal <= 511 {
			vertexColor = rl.NewColor(
				uint8(sclColorVal)-255,
				0,
				255,
				255,
			)
		} else if sclColorVal <= 767 {
			diff := sclColorVal - 511
			vertexColor = rl.NewColor(
				255,
				0,
				255-uint8(diff),
				255,
			)
		} else {
			vertexColor = rl.NewColor(255, 0, 0, 255)
		}
		// draw the rectangle
		rl.DrawRectangle(
			int32(vertex.Coordinate.X)*ent.gg.TileSize+int32(ent.GetPosition().X),
			int32(vertex.Coordinate.Y)*ent.gg.TileSize+int32(ent.GetPosition().Y),
			int32(ent.gg.TileSize),
			int32(ent.gg.TileSize),
			vertexColor,
		)
		// display the distance if it is not at max value
		if vertex.Distance != math.MaxInt {
			rl.DrawText(
				strconv.Itoa(vertex.Distance),
				int32(vertex.Coordinate.X)*ent.gg.TileSize+int32(ent.GetPosition().X),
				int32(vertex.Coordinate.Y)*ent.gg.TileSize+int32(ent.GetPosition().Y),
				int32(ent.TextSize),
				rl.Black,
			)
		}

	}
	// draw an arrow to the closest neighbour (if there is one) TODO:

	// draw grid
	for i := range ent.gg.Width + 1 {
		rl.DrawLine(
			int32(i)*ent.gg.TileSize+int32(ent.GetPosition().X),
			0+int32(ent.GetPosition().Y),
			int32(i)*ent.gg.TileSize+int32(ent.GetPosition().X),
			int32(ent.gg.Height)*ent.gg.TileSize+int32(ent.GetPosition().Y),
			rl.Black,
		)
	}
	for i := range ent.gg.Height + 1 {
		rl.DrawLine(
			0+int32(ent.GetPosition().X),
			int32(i)*ent.gg.TileSize+int32(ent.GetPosition().Y),
			int32(ent.gg.Width)*ent.gg.TileSize+int32(ent.GetPosition().X),
			int32(i)*ent.gg.TileSize+int32(ent.GetPosition().Y),
			rl.Black,
		)
	}
}

func (ent *GridGraphEntity) OnInputEvent(event *input.InputEvent) bool {
	// Logic to run when an input event is received.
	// Return false if the event was consumed and should not be propagated
	// further.
	// ...
	if event.Action == input.ActionClickRightHeld {
		ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.GetMouseDelta()))
	}
	sclMousePos := rl.Vector2Scale(
		rl.Vector2Subtract(
			rl.GetMousePosition(),
			ent.GetPosition(),
		),
		1/float32(ent.gg.TileSize),
	)
	sclMousePos = rl.NewVector2(float32(int(sclMousePos.X)), float32(int(sclMousePos.Y)))

	if event.Action == input.ActionClickDown {
		ent.gg.RemoveUnreachableTiles(sclMousePos)
		ent.gg.Dijkstra(sclMousePos)
	}
	if event.Action == input.ActionPlaceObstacle {
		ent.gg.SetObstacle(sclMousePos)
	}
	return true
}
