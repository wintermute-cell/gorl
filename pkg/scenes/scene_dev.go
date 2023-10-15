package scenes

import (
	"cowboy-gorl/pkg/ai/navigation"
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/gui"

	//"cowboy-gorl/pkg/lighting"
	"cowboy-gorl/pkg/logging"
)

// This checks at compile time if the interface is implemented
var _ Scene = (*DevScene)(nil)

// Dev Scene
type DevScene struct {
	entity_manager *proto.EntityManager
	scn_root_ent   *proto.BaseEntity
	pathable_world navigation.PathableWorld
	start_tile     navigation.Pathable
	end_tile       navigation.Pathable
	g              *gui.Gui
	boolmap        [][]bool
	navmap         *navigation.PathableWorld
}

func (scn *DevScene) Init() {
	scn.entity_manager = proto.NewEntityManager()
	scn.g = gui.NewGui()
	scn.scn_root_ent = &proto.BaseEntity{Name: "DevSceneRoot"}
	gem.AddEntity(gem.Root(), scn.scn_root_ent)
	//lighting.Enable()
	logging.Info("DevScene initialized.")
}

func (scn *DevScene) Deinit() {
    gem.RemoveEntity(scn.scn_root_ent)
	scn.entity_manager.DisableAllEntities()
	//lighting.Disable()
	logging.Info("DevScene de-initialized.")
}

func (scn *DevScene) DrawGUI() {
	scn.g.Draw()
}

func (scn *DevScene) Draw() {
	scn.entity_manager.UpdateEntities()
}
