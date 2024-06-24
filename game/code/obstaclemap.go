package code

import (
	"image/color"
	"time"

	"gorl/fw/core/logging"
	"gorl/fw/core/math"
	"gorl/fw/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// PheromoneType is a custom type for representing different pheromone channels.
type PheromoneType uint8

const (
	PheromoneEdgeObstacle PheromoneType = iota // Encodes the alpha channel as zero (represents obstacle or edge)
	PheromoneLeaving                           // Encodes the blue channel (seeking food)
	PheromoneReturning                         // Encodes the red channel (returning home)
	PheromoneNoFood                            // Encodes the green channel (indicates no food)
)

// ObstacleMap represents a map where each cell is a color.RGBA value.
type ObstacleMap struct {
	data       []color.RGBA    // One-dimensional array storing the map data
	width      int             // Width of the map
	height     int             // Height of the map
	resolution math.Vector2Int // Logical resolution of the map
	maxAge     time.Duration   // Maximum age for pheromones
}

// NewObstacleMap creates a new map with the specified size and resolution.
func NewObstacleMap(size, resolution math.Vector2Int) *ObstacleMap {
	totalCells := size.X * size.Y
	data := make([]color.RGBA, totalCells)

	// Initialize obstacles as an example
	areaWidth, areaHeight := 100, 100
	startX, startY := 10, 10
	for x := startX; x < startX+areaWidth; x++ {
		for y := startY; y < startY+areaHeight; y++ {
			if x < size.X && y < size.Y {
				data[y*size.X+x] = color.RGBA{0, 0, 0, 255} // Obstacle with zero alpha
			}
		}
	}

	return &ObstacleMap{
		data:       data,
		width:      size.X,
		height:     size.Y,
		resolution: resolution,
		maxAge:     20 * time.Second,
	}
}

// GetSize returns the width and height of the map.
func (om *ObstacleMap) GetSize() math.Vector2Int {
	return math.Vector2Int{X: om.width, Y: om.height}
}

// scalePosition converts from render resolution to the map's logical resolution.
func (om *ObstacleMap) scalePosition(position math.Vector2Int) math.Vector2Int {
	return math.Vector2Int{
		X: position.X * om.width / om.resolution.X,
		Y: position.Y * om.height / om.resolution.Y,
	}
}

// SetPheromone sets the pheromone level at a given position based on the specified pheromone type.
func (om *ObstacleMap) SetPheromone(position math.Vector2Int, pheromoneType PheromoneType) {
	logicalPosition := om.scalePosition(position)
	index := logicalPosition.Y*om.width + logicalPosition.X

	if index < 0 || index >= len(om.data) {
		return // Out of bounds check
	}

	currData := om.data[index]

	// Avoid setting pheromones on obstacles
	if currData.A == 255 {
		return
	}

	// Set the appropriate pheromone channel based on the type
	const addAmount = 50
	switch pheromoneType {
	case PheromoneLeaving:
		om.data[index] = color.RGBA{R: currData.R, G: currData.G, B: currData.B + addAmount, A: 200}
	case PheromoneReturning:
		om.data[index] = color.RGBA{R: currData.R + addAmount, G: currData.G, B: currData.B, A: 200}
	case PheromoneNoFood:
		om.data[index] = color.RGBA{R: currData.R, G: currData.G + addAmount, B: currData.B, A: 200}
	default:
		// Ignore setting if the type is PheromoneEdgeObstacle or any invalid type
	}
}

// DecayPheromones decreases the pheromone levels over time by decrementing the color channels.
func (om *ObstacleMap) DecayPheromones(decayRate uint8) {
	for i := range om.data {
		cell := &om.data[i]
		if cell.A != 255 { // Skip obstacles
			if cell.R > decayRate {
				cell.R = util.Clamp(cell.R-decayRate, 1, 255)
				cell.A = min(254, max(cell.R, cell.B, cell.G))
			}
			if cell.G > decayRate {
				cell.G = util.Clamp(cell.G-decayRate, 1, 255)
				cell.A = min(254, max(cell.R, cell.B, cell.G))
			}
			if cell.B > decayRate {
				cell.B = util.Clamp(cell.B-decayRate, 1, 255)
				cell.A = min(254, max(cell.R, cell.B, cell.G))
			}
		}
	}
}

// HasInCircle returns the number of cells with pheromone within a circle in the render resolution,
// and calculates the aged count considering the intensity of the specified pheromone channel.
func (om *ObstacleMap) HasInCircle(center math.Vector2Int, radius float32, pheromoneType PheromoneType) (count int32, agedCount float32) {
	logicalCenter := om.scalePosition(center)
	logicalRadius := radius * float32(om.width) / float32(om.resolution.X) // Assuming uniform scaling

	startX := logicalCenter.X - int(logicalRadius)
	endX := logicalCenter.X + int(logicalRadius)
	startY := logicalCenter.Y - int(logicalRadius)
	endY := logicalCenter.Y + int(logicalRadius)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {
			dx := float32(x - logicalCenter.X)
			dy := float32(y - logicalCenter.Y)
			if dx*dx+dy*dy <= logicalRadius*logicalRadius {
				if x >= 0 && x < om.width && y >= 0 && y < om.height {
					index := y*om.width + x
					cell := om.data[index]

					var pheromoneLevel uint8
					switch pheromoneType {
					case PheromoneEdgeObstacle:
						if cell.A == 255 {
							count++
						}
					case PheromoneLeaving:
						pheromoneLevel = cell.B
					case PheromoneReturning:
						pheromoneLevel = cell.R
					case PheromoneNoFood:
						pheromoneLevel = cell.G
					default:
						continue // Skip invalid or obstacle types
					}

					if pheromoneLevel > 0 {
						count++
						agedLevel := float32(pheromoneLevel) / 255.0
						agedCount += agedLevel * agedLevel // squaring has the effect of emphasizing fresher pheromones
					}
				} else if pheromoneType == PheromoneEdgeObstacle {
					count++ // Count out-of-bounds cells as pheromones
				}
			}
		}
	}

	return count, agedCount
}

// ToRlTexture updates a texture with the map's data.
func (om *ObstacleMap) ToRlTexture(tex rl.Texture2D) {
	if tex.Width != int32(om.width) || tex.Height != int32(om.height) {
		logging.Fatal("Texture size does not match map size, should be %dx%d", om.width, om.height)
	}
	rl.UpdateTexture(tex, om.data)
}
