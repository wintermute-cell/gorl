package graphutils

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Takes an black/white image and converts it into a silce of vectors which
// repres a grid graph. The slice contains only the obstacles, it is implied
// that every tile that is not declared an obstacle is a movable tile
func CalculateGridGraphImage(mapImage *rl.Image, tileSize int) []rl.Vector2 {
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
				gridGraphObstacles = append(gridGraphObstacles, rl.NewVector2(float32(x), float32(y)))
			}
		}
	}
	return gridGraphObstacles
}
