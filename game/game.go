package game

import "gorl/fw/modules/scenes"
import gsc "gorl/game/scenes"

type game struct {
	// This struct is used to store any game state shared between Init() and
	// Update() that can't be part of any existing systems like 'entities',
	// 'scenes', the global 'store', etc...
}

var state game

func Init() {
	// This code is run before the game loop starts.
	// NOTE: hier szenen initialisieren
	scenes.RegisterScene("flow field", &gsc.FlowFieldScene{})
	scenes.EnableScene("flow field")

}
