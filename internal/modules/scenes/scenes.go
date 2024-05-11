package scenes

import "gorl/internal/core/entities/proto"

// Scene is an interface that every scene in the game should implement
type Scene interface {
	// Init and Deinit are implemented by the user.
	Init()
	Deinit()

	// Update, and FixedUpdate are implemented by the framework but may
	// be overridden by the user.
	Update()
	FixedUpdate()

	GetRoot() *proto.Entity
}
