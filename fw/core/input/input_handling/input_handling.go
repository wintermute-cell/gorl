package input

import (
	"gorl/fw/core/entities/proto"
	input "gorl/fw/core/input/input_event"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// HandleInputEvents checks for input events and propagates them to the entities.
// Receives a sorted slice of layers, each containing a slice of entities.
// Both must be sorted from back to front (from far away to close to camera).
func HandleInputEvents(orderedEntities [][]proto.IEntity) {
	events := checkForInputs()
	for _, event := range events {
	entities: // stop iterating over all entitie if one blocks propagation.
		//for _, layer := range orderedEntities {
		//	for _, entity := range layer {
		// iterate in reverse
		for i := len(orderedEntities) - 1; i >= 0; i-- {
			layer := orderedEntities[i]
			for j := len(layer) - 1; j >= 0; j-- {
				entity := layer[j]
				shouldContinue := entity.OnInputEvent(event)
				if !shouldContinue {
					break entities
				}
			}
		}
	}
}

func checkForInputs() []*input.InputEvent {

	events := []*input.InputEvent{}
	mousePosition := rl.GetMousePosition()

	for action, triggers := range input.ActionMap {
		for _, trigger := range triggers {
			switch trigger.InputType {
			case input.InputTypeKey:
				switch trigger.TriggerType {
				case input.TriggerTypeDown:
					if rl.IsKeyDown(trigger.Key) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypePressed:
					if rl.IsKeyPressed(trigger.Key) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypeReleased:
					if rl.IsKeyReleased(trigger.Key) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				}
			case input.InputTypeMouse:
				switch trigger.TriggerType {
				case input.TriggerTypeDown:
					if rl.IsMouseButtonDown(trigger.MouseButton) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypePressed:
					if rl.IsMouseButtonPressed(trigger.MouseButton) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypeReleased:
					if rl.IsMouseButtonReleased(trigger.MouseButton) {
						events = append(events, input.NewInputEvent(action, mousePosition))
					}
				case input.TriggerTypePassive:
					events = append(events, input.NewInputEvent(action, mousePosition))
				}
			case input.InputTypeGamepad:
				// Implement the checks for gamepad buttons using a similar pattern
			}
		}
	}

	return events
}
