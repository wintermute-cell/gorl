package gem

import "gorl/fw/core/entities/proto"

// UpdateEntities traverses the hierarchy, calling Update() and FixedUpdate() and
// returns a map of layers each being a slice of entities.
// The returned map contains the slices already sorted, ready to be drawn.
func UpdateEntities(withFixedUpdate bool) map[int64][]proto.IEntity {
	entitiesAsLayers := traverseAndFlatten(gemInstance.root, true, withFixedUpdate)
	return entitiesAsLayers
}

// DrawEntitySlice draws the entities in the given slice.
// Choosing the entities to draw is up to the caller, and generally done with
// the return value of UpdateEntities().
func DrawEntitySlice(entities []proto.IEntity) {
	for _, entity := range entities {
		entity.Draw()
	}
}
