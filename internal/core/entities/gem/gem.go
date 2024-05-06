package gem

import (
	"gorl/internal/core/entities/proto"
	"gorl/internal/util"
	"sort"
	"strings"
)

type GlobalEntityManager struct {
	root               proto.Entity
	fixedUpdateTimer   util.Timer
	fixedUpdateEnabled bool
}

var Gem *GlobalEntityManager

// InitGem initializes the global entity manager with the given physics delta
// time. If the physics delta time is 0, fixed update is disabled.
func InitGem(physicsDeltatimeSeconds float32) {
	rootEntity := &proto.BaseEntity{Name: "GemRoot"}
	rootEntity.Init()
	Gem = &GlobalEntityManager{
		root:               rootEntity,
		fixedUpdateTimer:   *util.NewTimer(physicsDeltatimeSeconds),
		fixedUpdateEnabled: physicsDeltatimeSeconds > 0.0,
	}
}

// AddEntity adds a child to the root entity and initializes it. The child is
// returned.
func AddEntity(parent proto.Entity, child proto.Entity) proto.Entity {
	parent.AddChild(child)
	child.SetParent(parent)
	print_tree(Gem.root, 0, false)
	child.Init()
	return child
}

func print_tree(entity proto.Entity, depth int32, last bool) {
	prefix := strings.Repeat("   ", int(depth))
	if depth > 0 {
		if last {
			prefix += "└─ "
		} else {
			prefix += "├─ "
		}
	}

	children := entity.GetChildren()
	for i, child := range children {
		isLast := i == len(children)-1
		print_tree(child, depth+1, isLast)
	}
}

func RemoveEntity(entity proto.Entity) {
	// this function is weird. try to not touch it.
	// Slices seem to behave weirdly in Go. When a slice is returned, it is not
	// a copy, but rather a reference, since the returned slice shares the same
	// underlying array.
	children := make([]proto.Entity, len(entity.GetChildren()))
	copy(children, entity.GetChildren()) // gets entity.children (slice of pointers)
	for _, child := range children {
		RemoveEntity(child)
	}
	if entity.GetParent() != nil {
		entity.GetParent().RemoveChild(entity) // removes itself from parent.children
	}
	entity.Deinit()
}

func UpdateEntities() {
	if Gem.fixedUpdateEnabled && Gem.fixedUpdateTimer.Check() {
		updateEntity(Gem.root, true)
	} else {
		updateEntity(Gem.root, false)
	}
}

func DrawEntities() {
	draw_entities_flattened(Gem.root, false)
}

func DrawEntitiesGUI() {
	draw_entities_flattened(Gem.root, true)
}

func updateEntity(entity proto.Entity, with_fixed_update bool) {
	if with_fixed_update {
		entity.FixedUpdate()
	}
	entity.Update()
	for _, child := range entity.GetChildren() {
		updateEntity(child, with_fixed_update)
	}
}

func drawEntity(entity proto.Entity) {
	entity.Draw()
	for _, child := range entity.GetChildren() {
		drawEntity(child)
	}
}

// flattened_entity represents an entity with its depth in the tree.
// Depth will help us make sure that the tree structure is respected.
type flattened_entity struct {
	Entity proto.Entity
	Depth  int
}

// flatten_entities creates a flat list of all entities with their depth.
func flatten_entities(entity proto.Entity, depth int) []flattened_entity {
	var flattened_list []flattened_entity
	flattened_list = append(flattened_list, flattened_entity{Entity: entity, Depth: depth})

	for _, child := range entity.GetChildren() {
		flattened_list = append(flattened_list, flatten_entities(child, depth+1)...)
	}

	return flattened_list
}

func draw_entities_flattened(root proto.Entity, draw_gui bool) {
	// Flatten the entities
	flattened_entities := flatten_entities(root, 0)

	// Sort the flattened entities first by their GetDrawIndex() and then by their depth
	sort.SliceStable(flattened_entities, func(i, j int) bool {
		// If indices are equal, sort by depth (to respect tree structure)
		if flattened_entities[i].Entity.GetDrawIndex() == flattened_entities[j].Entity.GetDrawIndex() {
			return flattened_entities[i].Depth < flattened_entities[j].Depth
		}
		// Otherwise sort by draw index
		return flattened_entities[i].Entity.GetDrawIndex() < flattened_entities[j].Entity.GetDrawIndex()
	})

	// Draw in sorted order
	for _, fe := range flattened_entities {
		if draw_gui {
			fe.Entity.DrawGUI()
		} else {
			fe.Entity.Draw()
		}
	}
}

func Root() proto.Entity {
	return Gem.root
}
