package grid_graph

import (
	"gorl/fw/core/logging"
	"image/color"
	"math"
	"math/rand"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// basically used as an enum
const (
	DIJKSTRA_BLACK = "BLACK"
	DIJKSTRA_GREY  = "GREY"
	DIJKSTRA_WHITE = "WHITE"
)

// A GridGraph is a structure that consists of a set of vertices. NOTE: other graphs need edges, only a grid graph has them implicitly!
type GridGraph struct {
	// Vertices  []*Vertex NOTE: because I only need a grid graph that has coordinates, the vertices[] are not needed,
	// NOTE: they are included in the VertexMap
	Width      int
	Height     int
	VertexMap  map[Coordinate]*Vertex
	TileSize   int32
	DrawOffset rl.Vector2
	Robots     map[int]*Robot // index - robot
}

// A Vertex is a node that belongs to a graph and can have an arbitrary number
// of neighbouring vertices. (this implies an edge between two vertices)
type Vertex struct {
	// NOTE: the coordinates are just inside the Vertex struct for convenience, because for now I only need a grid graph.
	// in a real arbitrary graph there is no "Coordinate", thus it has to be removed in a proper implementation
	Coordinate       Coordinate
	Neighbours       []*Vertex
	Distance         int
	DijkstraColor    string
	Predecessor      *Vertex
	ClosestNeighbour *Vertex
}

// This is basically a Vector2 with int components
type Coordinate struct {
	X int
	Y int
}

// robot
type Robot struct {
	Coords Coordinate
	Color  rl.Color
	//====== for smooth movement
	Direction rl.Vector2
}

// Builds a new Grid Graph in the given dimensions
func NewGridGraph(width int, height int) *GridGraph {
	logging.Info(">> graph.go: constructing new grid graph")
	gridGraph := GridGraph{}
	gridGraph.Width = width
	gridGraph.Height = height
	gridGraph.VertexMap = make(map[Coordinate]*Vertex)
	gridGraph.DrawOffset = rl.Vector2Zero()
	gridGraph.TileSize = 40
	gridGraph.Robots = make(map[int]*Robot)
	// Loop width and height for initializing the array.
	for x := range width {
		for y := range height {
			newCoord := Coordinate{x, y}
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
			vertex := gridGraph.VertexMap[Coordinate{x, y}]
			if x == 0 && y == 0 {
				// edge case top left
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y + 1}])
			} else if x == 0 && y == height-1 {
				// edge case bottom left
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y - 1}])
			} else if x == width-1 && y == 0 {
				// edge case top right
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y + 1}])
			} else if x == width-1 && y == height-1 {
				// edge case bottom right
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y - 1}])
			} else if x == 0 {
				// edge case left
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y - 1}])
			} else if x == width-1 {
				// edge case right
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y - 1}])
			} else if y == 0 {
				// edge case top
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y + 1}])
			} else if y == height-1 {
				// edge case bottom
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y - 1}])
			} else {
				// all vertices in the "middle" of the graph
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y - 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x - 1, y + 1}])
				vertex.Neighbours = append(vertex.Neighbours, gridGraph.VertexMap[Coordinate{x + 1, y - 1}])
			}
		}
	}
	// NOTE: just a sample robot
	gridGraph.Robots[0] = &Robot{Coordinate{X: 20, Y: 20}, rl.Green, rl.Vector2Zero()}
	return &gridGraph
}

// Runs the dijkstra algorithm for the tile at target. It calculates the distance
// of all other tiles to the target. The algorithm will fail, if there are
// encircled tiles present. Call RemoveUnreachableTiles to avoid that problem
// and declare these encircled paths as obstacles.
func (gg *GridGraph) Dijkstra(target Coordinate) {
	// Only run the algorithm if therp is a target
	if targetVertex, ok := gg.VertexMap[target]; ok {
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
		remainingGraphMap := make(map[Coordinate]*Vertex)
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
					if _, ok := gg.VertexMap[Coordinate{X: nDirX, Y: nDirY}]; ok {
						if nDirX == 0 || nDirY == 0 {
							preferredNeighbour = gg.VertexMap[Coordinate{X: nDirX, Y: nDirY}]
						}
						closestVertex = gg.VertexMap[Coordinate{X: nDirX, Y: nDirY}]
					}
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

// Moves every robot one tile closer to the target by checking their neighbours
// distance and moving towards the lowest one, if there is none, the robot has
// reached the target
func (gg *GridGraph) MoveRobotsToTarget() {
	for _, robot := range gg.Robots {
		// check if we are at the target already, if yes, then skip
		if gg.VertexMap[robot.Coords].Distance == 0 {
			continue
		}
		neighbours := gg.VertexMap[robot.Coords].Neighbours
		var closestNeighbour *Vertex
		smallestDistance := math.MaxInt
		// find smallest distance and corresponding vert
		for _, nVert := range neighbours {
			// check for smaller distance
			if nVert.Distance < smallestDistance {
				// if another robot is on the tile already, skip this one (flag it)
				continueFlag := false
				for _, rob := range gg.Robots {
					if nVert.Coordinate == rob.Coords {
						continueFlag = true
					}
				}
				// only set the tile as closets neighbour if it wasn flagged
				if !continueFlag {
					smallestDistance = nVert.Distance
					closestNeighbour = nVert
				}
			}
		}

		// set robot to new position
		if closestNeighbour != nil {
			// TODO: set offset
			robot.Direction = rl.NewVector2(
				float32(closestNeighbour.Coordinate.X)-float32(robot.Coords.X),
				float32(closestNeighbour.Coordinate.Y)-float32(robot.Coords.Y),
			)
			robot.Coords = closestNeighbour.Coordinate
		}
	}
}

// Adds a robot to the graph
func (gg *GridGraph) AddRobot(position Coordinate) {
	// check if position is in the grid graph
	if _, ok := gg.VertexMap[position]; ok {
		robotColor := rl.NewColor(
			uint8(rand.Int()%256),
			uint8(rand.Int()%256),
			uint8(rand.Int()%256),
			255)
		newRobot := &Robot{position, robotColor, rl.Vector2Zero()}
		gg.Robots[len(gg.Robots)] = newRobot
	}
}

// The Dijkstra algorithm can only calculate the distance to the target if it is reachable, if a tile can
// not reach the target, the algorithm will fail. Call this function before Dijkstra to ensure that unreachable
// tiles will be removed (set as an obstacle), to avoid this fail.
func (gg *GridGraph) RemoveUnreachableTiles(position Coordinate) {
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
	newVertexMap := make(map[Coordinate]*Vertex)
	for _, vert := range reachableNodes {
		newVertexMap[vert.Coordinate] = vert
	}
	gg.VertexMap = newVertexMap
}

// Sets an "obstacle" in the graph. Basically removes a vertex from the grid graph
func (gg *GridGraph) SetObstacle(position Coordinate) {
	// check if position is (still) in the grid graph
	if vert, ok := gg.VertexMap[position]; ok {
		// remove the vert from its neighbours
		for _, nVert := range vert.Neighbours {
			index := slices.Index(nVert.Neighbours, vert)
			if index != -1 {
				nVert.Neighbours = slices.Delete(nVert.Neighbours, index, index+1)
			}
		}
		// remove the vert from the VertexMap itself
		delete(gg.VertexMap, vert.Coordinate)
	}
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
		gg.SetObstacle(Coordinate{X: int(ob.X), Y: int(ob.Y)})
	}
}
