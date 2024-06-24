package code

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type FoodPiles struct {
	FoodPoints      [][]rl.Vector2
	FoodPointRadius float32

	FoodPilePoints []rl.Vector2
	FoodPileRadius float32
}

func NewFoodPiles() *FoodPiles {
	// NOTE: you can modify the constructor to take any parameters you need to
	// initialize the entity.
	new_ent := &FoodPiles{
		FoodPoints:      [][]rl.Vector2{},
		FoodPointRadius: 3,
		FoodPileRadius:  20,
		FoodPilePoints: []rl.Vector2{
			rl.NewVector2(730, 840),
			rl.NewVector2(1560, 440),
		},
	}

	// Add food points
	amountFoodPerPile := 1000
	for idx, point := range new_ent.FoodPilePoints {
		new_ent.FoodPoints = append(new_ent.FoodPoints, []rl.Vector2{})
		for i := 0; i < amountFoodPerPile; i++ {
			randPoint := rl.NewVector2(
				point.X+((rand.Float32()-0.5)*2)*new_ent.FoodPileRadius,
				point.Y+((rand.Float32()-0.5)*2)*new_ent.FoodPileRadius,
			)
			new_ent.FoodPoints[idx] = append(new_ent.FoodPoints[idx], randPoint)
		}
	}

	return new_ent
}

// CheckForFoodInCircle checks if there is food within a circle.
// If food is found, it returns true and removes the food point.
func (ent *FoodPiles) CheckForFoodInCircle(center rl.Vector2, radius float32) bool {
	for pileIdx, pile := range ent.FoodPoints {
		// we first check if the circle is colliding with the food pile
		// as this is more performant than checking each food point.
		if rl.CheckCollisionCircles(center, radius, ent.FoodPilePoints[pileIdx], ent.FoodPileRadius) {
			for pointIdx, point := range pile {
				if point.X < 0 {
					continue
				}
				if rl.CheckCollisionCircles(center, radius, point, ent.FoodPointRadius) {
					ent.FoodPoints[pileIdx][pointIdx] = rl.NewVector2(-1, -1)
					return true
				}
			}
		}
	}
	return false
}

// Draw the food piles.
func (ent *FoodPiles) Draw() {
	for _, pile := range ent.FoodPoints {
		for _, point := range pile {
			if point.X < 0 {
				continue
			}
			rl.DrawCircleV(point, ent.FoodPointRadius, rl.Green)
		}
	}
}
