package input

import rl "github.com/gen2brain/raylib-go/raylib"

// Actions
type Action int32

const (
	ActionMoveUp Action = iota
	ActionMoveDown
	ActionMoveLeft
	ActionMoveRight
	ActionPlaceObstacle
	ActionMoveRobotsToTarget
	ActionPlaceRobot
	ActionClickDown
	ActionClickHeld
	ActionClickUp
	ActionClickRightHeld
	ActionMouseHover
	ActionEscape
	ActionEnterUiMode
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
	ActionPlaceObstacle: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyO},
	},
	ActionMoveRobotsToTarget: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeReleased, Key: rl.KeyM},
	},
	ActionPlaceRobot: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeReleased, Key: rl.KeyP},
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
	ActionClickRightHeld: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypeDown, MouseButton: rl.MouseRightButton},
	},
	ActionMouseHover: {
		{InputType: InputTypeMouse, TriggerType: TriggerTypePassive},
	},
	ActionEscape: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeDown, Key: rl.KeyEscape},
	},
	ActionEnterUiMode: {
		{InputType: InputTypeKey, TriggerType: TriggerTypeReleased, Key: rl.KeyU},
	},
	// Add other action-trigger mappings
}
