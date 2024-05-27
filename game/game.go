package game

type game struct {
	// This struct is used to store any game state shared between Init() and
	// Update() that can't be part of any existing systems like 'entities',
	// 'scenes', the global 'store', etc...
}

var state game

func Init() {
	// This code is run before the game loop starts.
	// NOTE: hier szenen initialisieren
}

func Update() {
	// This code is run every frame, before drawing.
}
