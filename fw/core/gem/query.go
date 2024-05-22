package gem

import (
	"gorl/fw/core/entities/proto"
	"gorl/fw/core/logging"
)

var preallocSize = 1000

// GetByLayer returns all entities in the Gem that are on the specified layer.
//
// The default layer is gem.DefaultLayer (alias for 0).
//
// ( Returns the layer index back, so we can save in what order the layers are
// rendered in the main loop. We use this info to pass input events in the
// correct order. )
func GetByLayer(layerIndex int64) []proto.IEntity {
	if layerIndex < 0 || layerIndex >= int64(len(gemInstance.entities)) {
		logging.Error("Layer index out of bounds: %v", layerIndex)
		return []proto.IEntity{}
	}
	return gemInstance.entities[layerIndex]
}

// GetAll returns all entities in the Gem.
func GetAll() [][]proto.IEntity {
	return gemInstance.entities
}
