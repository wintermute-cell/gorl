package gem

import (
	"gorl/fw/core/entities/proto"
	"sort"
)

// Calls the Draw method of all entities in the list, sorted by their draw index.
// Returns the sorted list of entities.
func Draw(entities []proto.IEntity) []proto.IEntity {
	if len(entities) > 2 {
		// NOTE: i think we can get away with sorting the entities slice, since no
		// one should rely on the order of entities in the slice
		sort.Slice(entities, func(i, j int) bool {
			return entities[i].GetDrawIndex() < entities[j].GetDrawIndex()
		})
	}

	for _, entity := range entities {
		entity.Draw()
	}

	return entities
}

// DrawLayer wraps GetByLayer and Draw to draw all entities on the specified layer.
// Returns the sorted list of entities.
func DrawLayer(layerIndex int64) []proto.IEntity {
	entities := GetByLayer(layerIndex)
	Draw(entities)
	return entities
}

// Update calls the Update method of all entities in the list.
func Update(entities []proto.IEntity) {
	for _, entity := range entities {
		entity.Update()
	}
}
