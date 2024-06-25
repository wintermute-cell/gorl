package code

import (
	"image/color"

	"gorl/fw/core/logging"
	"gorl/fw/core/math"
	"gorl/fw/util"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// PheromoneType is a custom type for representing different pheromone channels.
type PheromoneType uint8

const (
	PheromoneEdgeObstacle PheromoneType = iota // Encodes the alpha channel in obstacle data
	PheromoneLeaving                           // Encodes the blue channel in pheromone data (seeking food)
	PheromoneReturning                         // Encodes the red channel in pheromone data (returning home)
	PheromoneNoFood                            // Encodes the green channel in pheromone data (indicates no food left)
)

// PheromoneMap represents a map where each cell is a color.RGBA value.
type PheromoneMap struct {
	pheromoneData []color.RGBA    // Array storing the pheromone data
	obstacleData  []color.RGBA    // Array storing the obstacle data
	width         int             // Width of the map
	height        int             // Height of the map
	resolution    math.Vector2Int // Logical resolution of the map
}

// NewPheromoneMap creates a new map with the specified size and resolution.
func NewPheromoneMap(size, resolution math.Vector2Int) *PheromoneMap {
	totalCells := size.X * size.Y
	pheromoneData := make([]color.RGBA, totalCells)
	obstacleData := make([]color.RGBA, totalCells)

	// Initialize obstacles as an example
	//areaWidth, areaHeight := 100, 100
	//startX, startY := 10, 10
	//for x := startX; x < startX+areaWidth; x++ {
	//	for y := startY; y < startY+areaHeight; y++ {
	//		if x < size.X && y < size.Y {
	//			obstacleData[y*size.X+x] = color.RGBA{0, 0, 0, 255} // Obstacle with zero alpha
	//		}
	//	}
	//}

	return &PheromoneMap{
		pheromoneData: pheromoneData,
		obstacleData:  obstacleData,
		width:         size.X,
		height:        size.Y,
		resolution:    resolution,
	}
}

// GetSize returns the width and height of the map.
func (om *PheromoneMap) GetSize() math.Vector2Int {
	return math.Vector2Int{X: om.width, Y: om.height}
}

// scalePosition converts from render resolution to the map's logical resolution.
func (om *PheromoneMap) scalePosition(position math.Vector2Int) math.Vector2Int {
	return math.Vector2Int{
		X: position.X * om.width / om.resolution.X,
		Y: position.Y * om.height / om.resolution.Y,
	}
}

// SetPheromone sets the pheromone level at a given position based on the specified pheromone type.
func (om *PheromoneMap) SetPheromone(position math.Vector2Int, pheromoneType PheromoneType, amount uint8) {
	logicalPosition := om.scalePosition(position)
	index := logicalPosition.Y*om.width + logicalPosition.X

	if index < 0 || index >= len(om.pheromoneData) {
		return // Out of bounds check
	}

	currObstacleData := om.obstacleData[index]
	currPheromoneData := om.pheromoneData[index]

	// Avoid setting pheromones on obstacles
	if currObstacleData.A == 255 {
		return
	}

	// Set the appropriate pheromone channel based on the type
	var addAmount = amount
	switch pheromoneType {
	case PheromoneLeaving:
		om.pheromoneData[index] = color.RGBA{R: currPheromoneData.R, G: currPheromoneData.G, B: currPheromoneData.B + addAmount, A: 200}
		om.pheromoneData[index].A = max(om.pheromoneData[index].R, om.pheromoneData[index].B, om.pheromoneData[index].G)
	case PheromoneReturning:
		om.pheromoneData[index] = color.RGBA{R: currPheromoneData.R + addAmount, G: currPheromoneData.G, B: currPheromoneData.B, A: 200}
		om.pheromoneData[index].A = max(om.pheromoneData[index].R, om.pheromoneData[index].B, om.pheromoneData[index].G)
	case PheromoneNoFood:
		om.pheromoneData[index] = color.RGBA{R: currPheromoneData.R, G: currPheromoneData.G + addAmount, B: currPheromoneData.B, A: 200}
		om.pheromoneData[index].A = max(om.pheromoneData[index].R, om.pheromoneData[index].B, om.pheromoneData[index].G)
	default:
		// Ignore setting if the type is PheromoneEdgeObstacle or any invalid type
	}
}

// DecayPheromones decreases the pheromone levels over time by decrementing the color channels.
func (om *PheromoneMap) DecayPheromones(decayRate uint8) {
	for i := range om.pheromoneData {
		cell := &om.pheromoneData[i]
		if om.obstacleData[i].A != 255 { // Skip obstacles
			if cell.R > decayRate {
				cell.R = util.Clamp(cell.R-decayRate, 1, 255)
				cell.A = max(cell.R, cell.B, cell.G)
			}
			if cell.G > decayRate {
				cell.G = util.Clamp(cell.G-decayRate, 1, 255)
				cell.A = max(cell.R, cell.B, cell.G)
			}
			if cell.B > decayRate {
				cell.B = util.Clamp(cell.B-decayRate, 1, 255)
				cell.A = max(cell.R, cell.B, cell.G)
			}
		}
	}
}

// HasInCircle returns the number of cells with pheromone within a circle in the render resolution,
// and calculates the aged count considering the intensity of the specified pheromone channel.
func (om *PheromoneMap) HasInCircle(center math.Vector2Int, radius float32, pheromoneType PheromoneType) (count int32, agedCount float32) {
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
					pheroCell := om.pheromoneData[index]
					obstacleCell := om.obstacleData[index]

					var pheromoneLevel uint8
					switch pheromoneType {
					case PheromoneEdgeObstacle:
						if obstacleCell.A == 255 {
							count++
						}
					case PheromoneLeaving:
						pheromoneLevel = pheroCell.B
					case PheromoneReturning:
						pheromoneLevel = pheroCell.R
					case PheromoneNoFood:
						pheromoneLevel = pheroCell.G
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

// PheromoneToTexture updates a texture with the map's data.
func (om *PheromoneMap) PheromoneToTexture(tex rl.Texture2D) {
	if tex.Width != int32(om.width) || tex.Height != int32(om.height) {
		logging.Fatal("Texture size does not match map size, should be %dx%d", om.width, om.height)
	}
	rl.UpdateTexture(tex, om.pheromoneData)
}

func (om *PheromoneMap) ObstaclesToTexture(tex rl.Texture2D) {
	if tex.Width != int32(om.width) || tex.Height != int32(om.height) {
		logging.Fatal("Texture size does not match map size, should be %dx%d", om.width, om.height)
	}
	rl.UpdateTexture(tex, om.obstacleData)
}
