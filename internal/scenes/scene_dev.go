package scenes

import (
	"gorl/internal/entities/gem"
	"gorl/internal/entities/proto"
	"gorl/internal/gui"
	"gorl/internal/input"

	//"gorl/internal/lighting"
	"gorl/internal/logging"
)

// This checks at compile time if the interface is implemented
var _ Scene = (*DevScene)(nil)

// Dev Scene
type DevScene struct {
	scn_root_ent *proto.BaseEntity
	g            *gui.Gui
}

func (scn *DevScene) Init() {
	scn.g = gui.NewGui()
	scn.scn_root_ent = &proto.BaseEntity{Name: "DevSceneRoot"}
	logging.Info("DevScene initialized.")
}

func (scn *DevScene) Deinit() {
	gem.RemoveEntity(scn.scn_root_ent)
	//lighting.Disable()
	logging.Info("DevScene de-initialized.")
}

func (scn *DevScene) DrawGUI() {
	scn.g.Draw()
}

func (scn *DevScene) Draw() {
	if input.Triggered(input.ActionEscape) {
		Sm.DisableAllScenes()
		Sm.EnableScene("main_menu")
	}
}
