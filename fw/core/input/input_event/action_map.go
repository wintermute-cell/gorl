package input

import rl "github.com/gen2brain/raylib-go/raylib"

// Actions
type Action int32

const (
	ActionMoveUp Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionClickDown
	ActionRightClickDown
	ActionClickHeld
	ActionClickUp
	ActionMouseHover
	ActionEscape
	// Add other actions as needed
	ActionZoomIn
	ActionZoomOut
)

var ActionMap = map[Action][]Trigger{
	ActionMoveUp: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyW},
	},
	ActionMoveDown: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyS},
	},
	ActionMoveLeft: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyA},
	},
	ActionMoveRight: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyD},
	},
	ActionClickDown: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypePressed, MouseButton: rl.MouseLeftButton},
	},
	ActionRightClickDown: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypePressed, MouseButton: rl.MouseRightButton},
	},
	ActionClickHeld: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypeDown, MouseButton: rl.MouseLeftButton},
	},
	ActionClickUp: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypeReleased, MouseButton: rl.MouseLeftButton},
	},
	ActionMouseHover: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypePassive},
	},
	ActionEscape: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyEscape},
	},

	ActionZoomIn: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyQ},
	},
	ActionZoomOut: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyE},
	},
	// Add other action-trigger mappings
}
