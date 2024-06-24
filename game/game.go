package game

import (
	"gorl/fw/modules/scenes"
	gscenes "gorl/game/scenes"
)

func Init() {
	scenes.RegisterScene("phero", &gscenes.PheroSceneScene{})
	scenes.EnableScene("phero")
}
