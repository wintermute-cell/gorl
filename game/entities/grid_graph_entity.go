package entities

import (
	"gorl/fw/core/entities"
	"gorl/fw/core/gem"
	input "gorl/fw/core/input/input_event"
	"image/color"
	"math"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Ensure that GridGraphEntity implements IEntity.
var _ entities.IEntity = &GridGraphEntity{}

// GridGraph Entity
type GridGraphEntity struct {
	*entities.Entity // Required!j
	Gg               *GridGraph
	pixelTracks      map[rl.Vector2]rl.Color
}

// NewGridGraphEntity creates a new instance of the GridGraphEntity.
func NewGridGraphEntity() *GridGraphEntity {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &GridGraphEntity{
		Entity:      entities.NewEntity("GridGraphEntity", rl.Vector2Zero(), 0, rl.Vector2One()),
		Gg:          NewGridGraph(100, 100),
		pixelTracks: make(map[rl.Vector2]rl.Color),
	}
	return new_ent
}

func (ent *GridGraphEntity) Init() {
	mapImage := rl.LoadImage("./map_thresh.png")
	ent.Gg = NewGridGraph(48, 27)
	ent.Gg.CalculateGridGraphFromImage(mapImage, 40)
	// ent.gg.RemoveUnreachableTiles(rl.NewVector2(10, 0))
	// ent.gg.Dijkstra(rl.NewVector2(10, 0))
}

func (ent *GridGraphEntity) Deinit() {
}

func (ent *GridGraphEntity) Update() {
	// ROBOT MOVEMENT
	for _, robot := range gem.GetChildren(ent) {
		robotEntity, ok := robot.(*RobotEntity)
		if !ok {
			// ERROR
		}
		// set the target, in case it is not set
		if ent.Gg.Target != nil {
			robotEntity.FinalTarget = rl.Vector2Scale(ent.Gg.Target.Coordinate, 40)
			// add 20, 20 to find the center of a tile
			robotEntity.FinalTarget = rl.Vector2Add(robotEntity.FinalTarget, rl.NewVector2(20, 20))
		}

		robotEntity.WallAvoidanceVelocity = robotEntity.FindClosestWall(ent.Gg.ObstaclesVRenderSpace)

		//==========================================
		// flow field movement
		flowVector := ent.Gg.GetFlowVector(robotEntity.GetPosition())
		// if the flowVector is a Vector2Zero, we are either at the target, or drove into a wall,
		// so we want to stop the vehicle immediately
		if rl.Vector2Equals(rl.Vector2Zero(), flowVector) {
			robotEntity.CurrentTarget = robotEntity.GetPosition()
			robotEntity.Acceleration = rl.Vector2Zero()
			robotEntity.Velocity = rl.Vector2Zero()
		}

		// check for target != nil, because in that case dijkstra has not been called yet
		if ent.Gg.Target != nil && rl.Vector2Equals(robotEntity.GetTilePosition(), ent.Gg.Target.Coordinate) {
			// no acceleration
			robotEntity.CurrentTarget = robotEntity.GetPosition()
		} else {
			robotEntity.CurrentTarget = rl.Vector2Add(robotEntity.GetPosition(), flowVector)
			ent.pixelTracks[robotEntity.GetPosition()] = robotEntity.Color
		}
	}
}

func (ent *GridGraphEntity) Draw() {
	// Draw vertices
	rl.DrawRectangleV(
		ent.GetPosition(),
		rl.NewVector2(
			float32(ent.Gg.Width)*float32(ent.Gg.TileSize),
			float32(ent.Gg.Height)*float32(ent.Gg.TileSize),
		),
		rl.Black,
	)
	// TODO: put this in a vertex.GetColor() function
	for _, vertex := range ent.Gg.VertexMap {
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
			int32(vertex.Coordinate.X)*ent.Gg.TileSize+int32(ent.GetPosition().X),
			int32(vertex.Coordinate.Y)*ent.Gg.TileSize+int32(ent.GetPosition().Y),
			int32(ent.Gg.TileSize),
			int32(ent.Gg.TileSize),
			vertexColor,
		)
		// display the distance if it is not at max value
		// if vertex.Distance != math.MaxInt {
		// 	rl.DrawText(
		// 		strconv.Itoa(vertex.Distance),
		// 		int32(vertex.Coordinate.X)*ent.Gg.TileSize+int32(ent.GetPosition().X),
		// 		int32(vertex.Coordinate.Y)*ent.Gg.TileSize+int32(ent.GetPosition().Y),
		// 		int32(ent.Gg.TextSize),
		// 		rl.Black,
		// 	)
		// }

	}
	// draw the arrows
	// for _, vert := range ent.Gg.VertexMap {
	// 	if vert.ClosestNeighbour != nil {
	// 		// NOTE: HARDCODED!
	// 		scale := 40
	// 		arrowTipDirection := rl.Vector2Subtract(vert.ClosestNeighbour.Coordinate, vert.Coordinate)
	// 		arrowTipDirection = rl.Vector2Scale(arrowTipDirection, float32(scale))
	// 		// arrowTipDirection,
	// 		// rl.NewVector2(float32(scale/2), float32(scale/2)),
	// 		// )
	// 		// arrowTipPosition = rl.Vector2Scale(arrowTipDirection, 10)
	// 		// arrowTipPosition = rl.Vector2Scale(vert.Coordinate, float32(scale))
	//
	// 		// convert to world space
	// 		arrowOrigin := rl.Vector2Scale(vert.Coordinate, float32(scale))
	// 		// mid shift
	// 		arrowOrigin = rl.Vector2Add(arrowOrigin, rl.NewVector2(float32(scale)/2, float32(scale)/2))
	//
	// 		arrowTipPosition := rl.Vector2Add(
	// 			arrowOrigin,
	// 			rl.Vector2Scale(arrowTipDirection, 0.3),
	// 		)
	// 		// rl.Vector2Scale(arrowTipDirection, float32(scale)/3),
	// 		// )
	//
	// 		rl.DrawLineEx(
	// 			arrowOrigin,
	// 			arrowTipPosition,
	// 			3,
	// 			rl.Black,
	// 		)
	// 		rl.DrawCircleV(
	// 			arrowTipPosition,
	// 			4,
	// 			rl.Black,
	// 		)
	// 	}
	// }
	// draw grid
	for i := range ent.Gg.Width + 1 {
		rl.DrawLine(
			int32(i)*ent.Gg.TileSize+int32(ent.GetPosition().X),
			0+int32(ent.GetPosition().Y),
			int32(i)*ent.Gg.TileSize+int32(ent.GetPosition().X),
			int32(ent.Gg.Height)*ent.Gg.TileSize+int32(ent.GetPosition().Y),
			rl.Black,
		)
	}
	for i := range ent.Gg.Height + 1 {
		rl.DrawLine(
			0+int32(ent.GetPosition().X),
			int32(i)*ent.Gg.TileSize+int32(ent.GetPosition().Y),
			int32(ent.Gg.Width)*ent.Gg.TileSize+int32(ent.GetPosition().X),
			int32(i)*ent.Gg.TileSize+int32(ent.GetPosition().Y),
			rl.Black,
		)
	}

	// Draw tracks of robots
	for p, c := range ent.pixelTracks {
		rl.DrawPixelV(rl.Vector2Add(ent.GetPosition(), p), c)
	}

	// DEBUG draw obstacles
	for _, obst := range ent.Gg.ObstaclesVRenderSpace {
		rl.DrawCircleV(rl.Vector2Add(obst, ent.GetPosition()), 3, rl.White)
	}
}

func (ent *GridGraphEntity) OnInputEvent(event *input.InputEvent) bool {
	if event.Action == input.ActionClickRightHeld {
		ent.SetPosition(rl.Vector2Add(ent.GetPosition(), rl.GetMouseDelta()))
	}
	sclMousePos := rl.Vector2Scale(
		rl.Vector2Subtract(
			rl.GetMousePosition(),
			ent.GetPosition(),
		),
		1/float32(ent.Gg.TileSize),
	)
	sclMousePos = rl.NewVector2(float32(int(sclMousePos.X)), float32(int(sclMousePos.Y)))

	if event.Action == input.ActionClickDown {
		ent.Gg.RemoveUnreachableTiles(sclMousePos)
		ent.Gg.Dijkstra(sclMousePos)
	}
	if event.Action == input.ActionPlaceObstacle {
		ent.Gg.SetObstacle(sclMousePos)
	}
	return true
}

// =======================================//
// ===== GRID GRAPH IMPLEMENTATION ===== //
// basically used as an enum
const (
	DIJKSTRA_BLACK = "BLACK"
	DIJKSTRA_GREY  = "GREY"
	DIJKSTRA_WHITE = "WHITE"
)

// A GridGraph is a structure that consists of a set of vertices. NOTE: other graphs need edges, only a grid graph has them implicitly!
// TODO: move robots in own entity, the grid graph should only tell you what direction you should go
type GridGraph struct {
	// Vertices  []*Vertex NOTE: because I only need a grid graph that has coordinates, the vertices[] are not needed,
	// NOTE: they are included in the VertexMap
	Width                 int
	Height                int
	VertexMap             map[rl.Vector2]*Vertex
	ObstaclesVRenderSpace []rl.Vector2
	TileSize              int32
	DrawOffset            rl.Vector2
	TextSize              int
	Target                *Vertex
}

// A Vertex is a node that belongs to a graph and can have an arbitrary number
// of neighbouring vertices. (this implies an edge between two vertices)
type Vertex struct {
	// NOTE: the coordinates are just inside the Vertex struct for convenience, because for now I only need a grid graph.
	// in a real arbitrary graph there is no "Coordinate", thus it has to be removed in a proper implementation
	Coordinate       rl.Vector2
	Neighbours       []*Vertex
	Distance         int
	DijkstraColor    string
	Predecessor      *Vertex
	ClosestNeighbour *Vertex
}

// Builds a new Grid Graph in the given dimensions
func NewGridGraph(width int, height int) *GridGraph {
	gridGraph := GridGraph{}
	gridGraph.Width = width
	gridGraph.Height = height
	gridGraph.VertexMap = make(map[rl.Vector2]*Vertex)
	// gridGraph.Obstacles = make([]rl.Vector2, 50)
	gridGraph.DrawOffset = rl.Vector2Zero()
	gridGraph.TileSize = 40
	gridGraph.TextSize = 20
	// Loop width and height for initializing the array.
	for x := range width {
		for y := range height {
			newCoord := rl.NewVector2(float32(x), float32(y))
			newVertex := &Vertex{}
			newVertex.Coordinate = newCoord
			newVertex.Distance = math.MaxInt
			newVertex.DijkstraColor = DIJKSTRA_WHITE
			newVertex.Predecessor = nil
			newVertex.ClosestNeighbour = nil
			gridGraph.VertexMap[newCoord] = newVertex
			// gridGraph.Vertices = append(gridGraph.Vertices, newVertex)
		}
	}
	// Loop another time for connecting the vertices with each other.
	// Watch out for the borders of the grid.
	for x := range width {
		for y := range height {
			// for simplicity and readybility we check all 4 corners seperately
			// after that, we check the border vertecies, if no condition
			// applies, a vertex has 4 neighbours
			vertex := gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y))]
			if x == 0 && y == 0 {
				// edge case top left
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y+1))])
			} else if x == 0 && y == height-1 {
				// edge case bottom left
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y-1))])
			} else if x == width-1 && y == 0 {
				// edge case top right
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y+1))])
			} else if x == width-1 && y == height-1 {
				// edge case bottom right
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y-1))])
			} else if x == 0 {
				// edge case left
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y-1))])
			} else if x == width-1 {
				// edge case right
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y-1))])
			} else if y == 0 {
				// edge case top
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y+1))])
			} else if y == height-1 {
				// edge case bottom
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y-1))])
			} else {
				// all vertices in the "middle" of the graph
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y-1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x-1), float32(y+1))])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[rl.NewVector2(float32(x+1), float32(y-1))])
			}
		}
	}
	return &gridGraph
}

// Runs the dijkstra algorithm for the tile at target. It calculates the distance
// of all other tiles to the target. The algorithm will fail, if there are
// encircled tiles present. Call RemoveUnreachableTiles to avoid that problem
// and declare these encircled paths as obstacles.
func (gg *GridGraph) Dijkstra(target rl.Vector2) {
	// target is a raw mouse position, we need to cast it to an int
	targetScaled := rl.NewVector2(float32(int(target.X)), float32(int(target.Y)))
	// Only run the algorithm if therp is a target
	if targetVertex, ok := gg.VertexMap[targetScaled]; ok {
		// set the target if it is valid
		gg.Target = targetVertex
		// reset dijkstra color, distance, predecessor and closest neighbour
		for _, vertex := range gg.VertexMap {
			vertex.Distance = math.MaxInt
			vertex.DijkstraColor = DIJKSTRA_WHITE
			vertex.Predecessor = nil
			vertex.ClosestNeighbour = nil
		}

		targetVertex.DijkstraColor = DIJKSTRA_GREY
		targetVertex.Distance = 0

		// copy VertexMap
		remainingGraphMap := make(map[rl.Vector2]*Vertex)
		for k, v := range gg.VertexMap {
			remainingGraphMap[k] = v
		}
		// loop the graph until we run out of unvisited elements
		for len(remainingGraphMap) > 0 {
			// extract minimum distance vertex
			currMinDistance := math.MaxInt
			var currVertex *Vertex
			for _, vert := range remainingGraphMap {
				if vert.Distance < currMinDistance && vert.DijkstraColor == DIJKSTRA_GREY {
					currMinDistance = vert.Distance
					currVertex = vert
				}
			}
			// check neighbours of the vertex with the smallest distance
			for _, nVert := range currVertex.Neighbours {
				// Relax - NOTE: hardcoded 1 as weight, because in this use case they are always 1
				if nVert.DijkstraColor != DIJKSTRA_BLACK && nVert.Distance > currVertex.Distance+1 {
					nVert.DijkstraColor = DIJKSTRA_GREY
					nVert.Distance = currVertex.Distance + 1
					nVert.Predecessor = currVertex
				}
			}
			currVertex.DijkstraColor = DIJKSTRA_BLACK
			delete(remainingGraphMap, currVertex.Coordinate)
		}
		// set the closest neighbour for each vertex. Prefer horizontal
		// and vertical paths (one component of the neighbouring
		// position is 0
		for _, v := range gg.VertexMap {
			closestVertex := v
			var preferredNeighbour *Vertex
			for _, nVert := range v.Neighbours {
				if nVert != nil && nVert.Distance < closestVertex.Distance {
					nDirX := nVert.Coordinate.X - v.Coordinate.X
					nDirY := nVert.Coordinate.Y - v.Coordinate.Y
					if nDirX == 0 || nDirY == 0 {
						preferredNeighbour = nVert
					}
					closestVertex = nVert
				}
			}
			// check for distance != 0 to exclude the target itsel
			if v.Distance != 0 {
				// use preferred neighbour if there is one
				if preferredNeighbour != nil {
					v.ClosestNeighbour = preferredNeighbour
				} else {
					v.ClosestNeighbour = closestVertex
				}

			} else {
				// reset for target NOTE: kp was ich mir dabei gedacht hab, wird schon stimmen
				v.ClosestNeighbour = nil
			}
		}

	}
}

// The Dijkstra algorithm can only calculate the distance to the target if it is reachable, if a tile can
// not reach the target, the algorithm will fail. Call this function before Dijkstra to ensure that unreachable
// tiles will be removed (set as an obstacle), to avoid this fail.
// position is expected to be already scaled to a grid position
func (gg *GridGraph) RemoveUnreachableTiles(position rl.Vector2) {
	// if the tile is not in the VertexMap (already removed), do nothing
	if _, ok := gg.VertexMap[position]; !ok {
		return
	}
	// nodes that are connectet to position
	var reachableNodes []*Vertex
	// nodes that are neighbours of the current node, but not yet completely checked
	var nodesToCheck []*Vertex
	// completely checked nodes
	var checkedNodes []*Vertex

	// initially add the starting position to the nodesToVisit
	nodesToCheck = append(nodesToCheck, gg.VertexMap[position])

	for len(nodesToCheck) > 0 {
		currentNode := nodesToCheck[0]
		// adds the neighbours of the current position to the nodesToCheck...
		for _, nVert := range currentNode.Neighbours {
			//... but only if they are not already checked and not queued to be checked
			if !slices.Contains(checkedNodes, nVert) && !slices.Contains(nodesToCheck, nVert) {
				nodesToCheck = append(nodesToCheck, nVert)
			}

		}
		// if we are done with adding nodes to be checked, we add this node to the reachableNodes
		// and the checkedNodes, and remove it from nodesToCheck
		reachableNodes = append(reachableNodes, currentNode)
		checkedNodes = append(checkedNodes, currentNode)
		// deletes currentNode from nodesToCheck
		index := slices.Index(nodesToCheck, currentNode)
		if index != -1 {
			nodesToCheck = slices.Delete(nodesToCheck, index, index+1)
		}
	}
	// replace the old vertex map with the new one, remove it from the old one
	// so that we can call SetObstacle on them, to add them properly to
	// ObstaclesVRenderSpace
	newVertexMap := make(map[rl.Vector2]*Vertex)
	for _, vert := range reachableNodes {
		newVertexMap[vert.Coordinate] = vert
		delete(gg.VertexMap, vert.Coordinate)
	}
	// call set obstacles on the unreachable tiles
	for k, _ := range gg.VertexMap {
		gg.SetObstacle(k)
	}
	gg.VertexMap = newVertexMap
}

// Sets an "obstacle" in the graph. Basically removes a vertex from the grid graph
func (gg *GridGraph) SetObstacle(position rl.Vector2) {
	// add the obstacle to its slice if it is not already in there
	obstVRenderSpace := rl.Vector2Add(rl.Vector2Scale(position, 40), rl.NewVector2(20, 20))
	i := slices.Index(gg.ObstaclesVRenderSpace, obstVRenderSpace)
	if i == -1 {
		gg.ObstaclesVRenderSpace = append(gg.ObstaclesVRenderSpace, obstVRenderSpace)
	}
	// check if position is (still) in the grid graph
	if obstacleVert, ok := gg.VertexMap[position]; ok {
		// remove the vert from its neighbours
		for _, nVert := range obstacleVert.Neighbours {
			// remove the obstacle it self from its neighbour
			index := slices.Index(nVert.Neighbours, obstacleVert)
			if index != -1 {
				nVert.Neighbours = slices.Delete(nVert.Neighbours, index, index+1)
			}
			// positions for the horizontal and vertical neighbours
			directionLeft := rl.Vector2Add(obstacleVert.Coordinate, rl.NewVector2(-1, 0))
			directionRight := rl.Vector2Add(obstacleVert.Coordinate, rl.NewVector2(1, 0))
			directionUp := rl.Vector2Add(obstacleVert.Coordinate, rl.NewVector2(0, -1))
			directionDown := rl.Vector2Add(obstacleVert.Coordinate, rl.NewVector2(0, 1))
			// remove the vertices that are 1 tile apart in horizontal or vertical direction as neighbours as well
			// for that, we need the direction of the obstacle seen from the neighbouring vertex, we only need either
			// the horizontal two vertices, or the vertical two vertices
			obstacleDirection := rl.Vector2Subtract(obstacleVert.Coordinate, nVert.Coordinate)
			if obstacleDirection.X == 0 {

				if obstacleDirection.Y == 1 || obstacleDirection.Y == -1 {
					if otherNVert, ok := gg.VertexMap[directionLeft]; ok {
						// remove the upper and lower obstaclet from its neighbour
						index := slices.Index(nVert.Neighbours, otherNVert)
						if index != -1 {
							nVert.Neighbours = slices.Delete(nVert.Neighbours, index, index+1)
						}
					}
					if otherNVert, ok := gg.VertexMap[directionRight]; ok {
						// remove the upper and lower obstaclet from its neighbour
						index := slices.Index(nVert.Neighbours, otherNVert)
						if index != -1 {
							nVert.Neighbours = slices.Delete(nVert.Neighbours, index, index+1)
						}
					}
				}
			}
			if obstacleDirection.Y == 0 {
				if obstacleDirection.X == 1 || obstacleDirection.X == -1 {
					if otherNVert, ok := gg.VertexMap[directionUp]; ok {
						// remove the upper and lower obstaclet from its neighbour
						index := slices.Index(nVert.Neighbours, otherNVert)
						if index != -1 {
							nVert.Neighbours = slices.Delete(nVert.Neighbours, index, index+1)
						}
					}
					if otherNVert, ok := gg.VertexMap[directionDown]; ok {
						// remove the upper and lower obstaclet from its neighbour
						index := slices.Index(nVert.Neighbours, otherNVert)
						if index != -1 {
							nVert.Neighbours = slices.Delete(nVert.Neighbours, index, index+1)
						}
					}
				}
			}
		}
		// remove the vert from the VertexMap itself
		delete(gg.VertexMap, obstacleVert.Coordinate)
	}
}

// Returns the direction a robot should take (if it is
// inside a valid vertex of the grid graph), otherwise return Vector2Zero
// pos is a screen space vector
func (gg *GridGraph) GetFlowVector(pos rl.Vector2) rl.Vector2 {
	// determine which tile the robot is on
	tilePosition := rl.Vector2Scale(pos, 1/40.0)
	tilePosition = rl.NewVector2(
		float32(int(tilePosition.X)),
		float32(int(tilePosition.Y)),
	)
	direction := rl.Vector2Zero()
	// return the vector to the closest neighbour, if there is one
	if vert, ok := gg.VertexMap[tilePosition]; ok {
		if vert.ClosestNeighbour != nil {
			direction = rl.Vector2Subtract(vert.ClosestNeighbour.Coordinate, vert.Coordinate)
		}
	}
	return direction
}

// Takes an black/white image and converts it into a silce of vectors which
// repres a grid graph. The slice contains only the obstacles, it is implied
// that every tile that is not declared an obstacle is a movable tile
func (gg *GridGraph) CalculateGridGraphFromImage(mapImage *rl.Image, tileSize int) {
	// this is a one dimensional array
	imgColors := rl.LoadImageColors(mapImage)
	var gridGraphObstacles []rl.Vector2

	// converting to 2D slice
	imgColors2D := make([][]color.RGBA, mapImage.Height)
	for i := range imgColors2D {
		imgColors2D[i] = make([]color.RGBA, mapImage.Width)
	}
	// add values to the 2D slice
	for x := range imgColors2D {
		for y := range imgColors2D[x] {
			imgColors2D[x][y] = imgColors[x*int(mapImage.Width)+y]
		}
	}
	sclWidth := len(imgColors2D) / tileSize
	sclHeight := len(imgColors2D[0]) / tileSize
	// loop the tile size
	for x := range sclWidth {
		for y := range sclHeight {
			// now loop the inner pixels of each tile and determine the predominant color
			whiteCount := 0
			blackCount := 0
			for p := range tileSize {
				for q := range tileSize {
					// since R, G and B are all 255 in a white pixel, its enough to just check one value
					if imgColors2D[x*tileSize+p][y*tileSize+q].R == 255 {
						whiteCount++
					} else {
						blackCount++
					}
				}
			}
			if whiteCount <= blackCount {
				gridGraphObstacles = append(gridGraphObstacles, rl.NewVector2(float32(y), float32(x)))
			}
		}
	}
	for _, ob := range gridGraphObstacles {
		gg.SetObstacle(rl.NewVector2(ob.X, ob.Y))
	}
}
