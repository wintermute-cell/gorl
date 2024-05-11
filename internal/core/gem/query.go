package gem

import (
	"gorl/internal/core/entities/proto"
	"gorl/internal/logging"
	"gorl/internal/util/langutils"
)

var preallocSize = 100

// GetByLayer returns all entities in the Gem that are on the specified layer.
//
// The default layer is gem.DefaultLayer (alias for 0).
func GetByLayer(layerIndex int64) []proto.IEntity {
	if layerIndex < 0 || layerIndex >= int64(len(gemInstance.entities)) {
		logging.Error("Layer index out of bounds: %v", layerIndex)
		return []proto.IEntity{}
	}
	return gemInstance.entities[layerIndex]
}

// GetAll returns all entities in the Gem.
func GetAll() []proto.IEntity {
	res := make([]proto.IEntity, 0, preallocSize)
	for _, layer := range gemInstance.entities {
		res = append(res, layer...)
	}
	return res
}

// FilterByType filters the entities by the type T, returning only the entities
// that implement the interface T.
//
// If inverse is true, it returns the entities that do not implement the
// interface T.
func FilterByType[T any](entities []proto.IEntity, inverse bool) []T {
	res := make([]T, 0, preallocSize)
	for _, entity := range entities {
		if langutils.ImplementsInterface[T](entity) {
			res = append(res, entity.(T))
		}
	}
	return res
}
