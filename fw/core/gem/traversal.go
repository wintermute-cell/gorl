package gem

import (
	"gorl/fw/core/collections"
	"gorl/fw/core/entities"
	input "gorl/fw/core/input/input_handling"
	"gorl/fw/core/math"
	"gorl/fw/core/render"
)

var _ render.Drawable = &WrappedEntity{}

type WrappedEntity struct {
	entities.IEntity
	absTransform math.Transform2D
}

// ShouldDraw checks if the entity should be drawn based on its layer flags,
// enabled and visible properties.
func (d WrappedEntity) ShouldDraw(layerFlags math.BitFlag) bool {
	e := d.IEntity.IsEnabled()
	v := d.IEntity.IsVisible()
	f := d.IEntity.GetLayerFlags().IsAny(layerFlags)
	return e && v && f
}

// Draw draws the entity.
func (d WrappedEntity) Draw() {
	oldTransform := *d.IEntity.GetTransform() // save the entity's old *local* transform
	d.IEntity.SetTransform(d.absTransform)    // set the entity's transform to the new transform matrix
	d.IEntity.Draw()                          // draw the entity
	d.IEntity.SetTransform(oldTransform)      // restore the entity's old *local* transform
}

// GetEntity retrieves the wrapped entity.
func (d WrappedEntity) GetEntity() entities.IEntity {
	return d.IEntity
}

// AsDrawable returns the entity as a drawable.
func (d WrappedEntity) AsDrawable() render.Drawable {
	return d
}

// AsInputReceiver returns the entity as an input receiver.
func (d WrappedEntity) AsInputReceiver() input.InputReceiver {
	return d
}

// Traverse traverses through the entity graph, updating the entities.
// In the process, it produces a list of DrawableEntity objects.
func Traverse(withFixedUpdate bool) ([]render.Drawable, []input.InputReceiver) {

	root := gemInstance.root

	nodeStack := collections.NewStack[*gemNode](0)
	nodeStack.Push(root)

	transformStack := collections.NewStack[math.Matrix3](0)
	transformStack.Push(math.Matrix3Identity())

	drawables := make([]render.Drawable, 0)
	inputReceivers := make([]input.InputReceiver, 0)

	for !nodeStack.IsEmpty() {

		node, _ := nodeStack.Pop()
		tMat3, _ := transformStack.Pop()

		// of the entity is not enabled, skip it and its children
		if !node.entity.IsEnabled() {
			continue
		}

		// add the enabler to the input receivers
		inputReceivers = append(inputReceivers, node.entity)

		// Update the entity
		node.entity.Update()
		if withFixedUpdate {
			node.entity.FixedUpdate()
		}

		drawables = append(drawables, WrappedEntity{
			IEntity:      node.entity,
			absTransform: math.NewTransform2DFromMatrix3(tMat3),
		})

		for _, child := range node.children {
			nodeStack.Push(child)
			transformStack.Push( // we push M_child * M_stack
				child.entity.
					GetTransform().
					GenerateMatrix().
					Multiply(tMat3))
		}
	}

	return drawables, inputReceivers
}
