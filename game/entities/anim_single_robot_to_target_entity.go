package entities

import (
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_event"
	"gorl/game/code/grid_graph"
	"math"
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that AnimSingleRobotToTargetEntity implements IEntity.
var _ entities.IEntity = &AnimSingleRobotToTargetEntity{}

// AnimSingleRobotToTarget Entity
type AnimSingleRobotToTargetEntity struct {
	*entities.Entity // Required!

	// Custom Fields
	gg             *grid_graph.GridGraph
	sclSec         float32
	faderRobot     int
	moveDelay      float32 // delay between each robot movement
	isRobotsMoving bool
	robotSec       float32 // helper for robot seconds counting
	inputCounter   int     // for fixing double input action
	TextSize       int
}

// NewAnimSingleRobotToTargetEntity creates a new instance of the AnimSingleRobotToTargetEntity.
func NewAnimSingleRobotToTargetEntity() *AnimSingleRobotToTargetEntity {
	new_ent := &AnimSingleRobotToTargetEntity{
		Entity: entities.NewEntity("AnimSingleRobotToTargetEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		gg:     grid_graph.NewGridGraph(48, 27),
	}
	return new_ent
}

func (ent *AnimSingleRobotToTargetEntity) Init() {
	// Initialization logic for the entity
	// ...
	mapImage := rl.LoadImage("./map_thresh.png")
	ent.gg = grid_graph.NewGridGraph(48, 27)
	ent.gg.CalculateGridGraphFromImage(mapImage, 40)
	ent.gg.RemoveUnreachableTiles(grid_graph.Coordinate{X: 10, Y: 0})
	ent.gg.Dijkstra(grid_graph.Coordinate{X: 10, Y: 0})
	ent.faderRobot = 0
	ent.moveDelay = 1000
	ent.isRobotsMoving = false
	ent.TextSize = 20
}

func (ent *AnimSingleRobotToTargetEntity) Deinit() {
	// De-initialization logic for the entity
	// ...
}

func (ent *AnimSingleRobotToTargetEntity) Update() {
	// Update logic for the entity per frame
	// ...
	ent.sclSec += rl.GetFrameTime() * 1000
	// trigger robot display
	if ent.sclSec > 3000 && ent.faderRobot < 255 {
		ent.faderRobot++

	}
	// move the robot
	if ent.isRobotsMoving {
		ent.robotSec += rl.GetFrameTime() * 1000
		if ent.robotSec > ent.moveDelay {
			ent.robotSec -= ent.moveDelay
			ent.gg.MoveRobotsToTarget()
		}

	}
}

func (ent *AnimSingleRobotToTargetEntity) Draw() {
	// NOTE: THIS IS COPIED FROM grid_graph_entity.go,
	// for fading effects i need this here because i cant be
	// bothered with cleaning it up

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
	// draw the robots
	for _, robot := range ent.gg.Robots {
		robotColor := rl.NewColor(robot.Color.R, robot.Color.G, robot.Color.B, uint8(ent.faderRobot))
		rl.DrawCircle(
			int32(int32(robot.Coords.X)*ent.gg.TileSize+ent.gg.TileSize/2)+int32(ent.GetPosition().X),
			int32(int32(robot.Coords.Y)*ent.gg.TileSize+ent.gg.TileSize/2)+int32(ent.GetPosition().Y),
			10,
			robotColor,
		)
	}
}

func (ent *AnimSingleRobotToTargetEntity) OnInputEvent(event *input.InputEvent) bool {
	// NOTE: RIPPED INPUT FROM grid_graph_entity.go
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
	sclCoord := grid_graph.Coordinate{X: int(sclMousePos.X), Y: int(sclMousePos.Y)}

	if event.Action == input.ActionClickDown {
		ent.gg.Dijkstra(sclCoord)
	}
	if event.Action == input.ActionPlaceObstacle {
		ent.gg.SetObstacle(sclCoord)
	}
	// start the robots
	if event.Action == input.ActionMoveRobotsToTarget {
		// fix for double call
		ent.inputCounter++
		if ent.inputCounter >= 2 {
			ent.isRobotsMoving = !ent.isRobotsMoving
			ent.inputCounter = 0
		}
	}
	if event.Action == input.ActionPlaceRobot {
		// NOTE: warum wird das immer zweimal gecallet?
		ent.inputCounter++
		if ent.inputCounter >= 2 {
			ent.gg.AddRobot(sclCoord)
			ent.inputCounter = 0
		}
	}
	return true
}
