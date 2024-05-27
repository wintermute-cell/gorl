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
	ActionClickHeld
	ActionClickUp
	ActionMouseHover
	ActionEscape
	// Add other actions as needed
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
	// Add other action-trigger mappings
}
