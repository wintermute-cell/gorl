package datastructures

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Based on https://stackoverflow.com/a/48330314

type QTBounds struct {
	X      int32
	Y      int32
	Width  int32
	Height int32
}

type QTNode struct {
	// Index of the first child qtnode if this node is a branch node.
	// Index of the first element node if this node is a leaf node.
	firstChild int32
	// NOTE: children are stored contiguously in blocks of 4.
	// NW, SW, NE, SE

	// -1 if this node is a branch node.
	// Number of elements if this node is a leaf node.
	count int32
}

type QTElement struct {
	// ID of the element. Can be used to refer to external data.
	ID int32

	// Bounds of the element.
	Bounds QTBounds
}

func (qt *QuadTree) NewQTElement(id int32, bounds QTBounds) QTElement {
	return QTElement{ID: id, Bounds: bounds}
}

type QTElemNode struct {
	// Index of the next element in the leaf node. -1 if this is the last element.
	nextElemNode int32

	// Index of the element.
	elementIdx int32
}

type QuadTree struct {
	// Stores all the elements of the quadtree.
	elements *FreeList[QTElement]

	// Stores all the element nodes of the quadtree. An element may consist of
	// multiple nodes if it spans multiple cells.
	elementNodes *FreeList[QTElemNode]

	// Stores all the nodes of the quadtree. The first node is always the root.
	nodes []*QTNode

	// Index of the first free node. This and the next three indices are free,
	// since subdivisions are managed contiguously.
	// A value of -1 indicates that the free list is empty, so we just append 4
	// nodes after the current last element.
	firstFreeNode int32

	// Stores the bounds of the entire quadtree.
	bounds QTBounds

	// The maximum number of elements that can be stored in a leaf node.
	leafCapacity int32

	// The maximum depth of the quadtree.
	maxDepth int32
}

func NewQuadTree(bounds QTBounds, leafCapacity int32, maxDepth int32) *QuadTree {
	root := &QTNode{firstChild: -1, count: 0}
	return &QuadTree{
		elements:      NewFreeList[QTElement](0),
		elementNodes:  NewFreeList[QTElemNode](0),
		nodes:         []*QTNode{root},
		firstFreeNode: -1,
		bounds:        bounds,
		leafCapacity:  leafCapacity,
		maxDepth:      maxDepth,
	}
}

// leavesInRect returns all the leaf nodes that intersect with the given rectangle, and a slice of their depths in the tree.
func (qt *QuadTree) leavesInRect(target QTBounds) ([]*QTNode, []int32) {

	// NodeData is a helper struct that stores a node and its dynamically
	// calculated bounds.
	type NodeData struct {
		node  *QTNode
		bound QTBounds
		depth int32
	}

	ret := make([]*QTNode, 0, 32)
	depths := make([]int32, 0, 32)
	toProcess := NewStack[NodeData](0)
	toProcess.Push(NodeData{node: qt.nodes[0], bound: qt.bounds, depth: 0}) // Push the root node.

	for toProcess.Size() > 0 {
		ndat, _ := toProcess.Pop()

		// If this node is a leaf node, add it to the result.
		if ndat.node.count != -1 {
			ret = append(ret, ndat.node)
			depths = append(depths, ndat.depth)
		} else {
			// If this node is a branch node, process its children.
			// First, we calculate the AABBs of the children by subdividing the parent bounds.
			// This can be done quickly with bit manipulation.

			hw := ndat.bound.Width >> 1  // Half width.
			hh := ndat.bound.Height >> 1 // Half height.
			mx := ndat.bound.X + hw      // Horizontal mid.
			hy := ndat.bound.Y + hh      // Vertical mid.

			// Then we check which children intersect with the target
			// rectangle, and push them to the stack.

			if target.X <= mx {
				// Target is in the left half...
				if target.Y <= hy {
					// ...and in the top half.
					toProcess.Push(NodeData{
						node:  qt.nodes[ndat.node.firstChild],
						bound: QTBounds{X: ndat.bound.X, Y: ndat.bound.Y, Width: hw, Height: hh},
						depth: ndat.depth + 1,
					})
				}
				if target.Y+target.Height > hy {
					// ...and in the bottom half.
					toProcess.Push(NodeData{
						node:  qt.nodes[ndat.node.firstChild+1],
						bound: QTBounds{X: ndat.bound.X, Y: hy, Width: hw, Height: hh},
						depth: ndat.depth + 1,
					})
				}
			}
			if target.X+target.Width > mx {
				// Target is in the right half...
				if target.Y <= hy {
					// ...and in the top half.
					toProcess.Push(NodeData{
						node:  qt.nodes[ndat.node.firstChild+2],
						bound: QTBounds{X: mx, Y: ndat.bound.Y, Width: hw, Height: hh},
						depth: ndat.depth + 1,
					})
				}
				if target.Y+target.Height > hy {
					// ...and in the bottom half.
					toProcess.Push(NodeData{
						node:  qt.nodes[ndat.node.firstChild+3],
						bound: QTBounds{X: mx, Y: hy, Width: hw, Height: hh},
						depth: ndat.depth + 1,
					})
				}
			}
		}
	}
	return ret, depths
}

func (qt *QuadTree) ElementsInRect(target QTBounds) []QTElement {
	leaves, _ := qt.leavesInRect(target)
	ret := make([]QTElement, 0, 32)
	for _, leaf := range leaves {
		if leaf.count == 0 {
			continue
		}
		elemNode := qt.elementNodes.Get(int(leaf.firstChild))
		for {
			ret = append(ret, qt.elements.Get(int(elemNode.elementIdx)))
			if elemNode.nextElemNode == -1 {
				break
			}
			elemNode = qt.elementNodes.Get(int(elemNode.nextElemNode))
		}
	}
	return ret
}

// Cleanup removes empty leaf nodes from the quadtree.
func (qt *QuadTree) Cleanup() {
	// NOTE: This function only does a single pass through the tree.
	// This will not completely remove all empty branches, but should suffice
	// if called repeatedly over time.

	toProcess := NewStack[int32](0)

	// Only process if the root is a branch node.
	if qt.nodes[0].count == -1 {
		toProcess.Push(0)
	}

	for toProcess.Size() > 0 {
		idx, _ := toProcess.Pop()
		firstChild := qt.nodes[idx].firstChild

		// Check if all children are empty.
		emptyLeaves := 0
		for i := 0; i < 4; i++ {
			child := qt.nodes[firstChild+int32(i)]
			if child.count == 0 {
				// If the child is an empty leaf node, increment the counter.
				emptyLeaves++
			} else if child.count == -1 {
				// If the child is a branch node, push it to the stack.
				toProcess.Push(firstChild + int32(i))
			}
		}

		if emptyLeaves == 4 {
			// Remove the children by adding them to the free list.
			qt.nodes[firstChild].firstChild = qt.firstFreeNode
			qt.firstFreeNode = firstChild

			// Make this node the new empty leaf node instead of its children.
			qt.nodes[idx].firstChild = -1
			qt.nodes[idx].count = 0
		}
	}
}

func (qt *QuadTree) Remove(element QTElement) {
	// Find all leaf nodes that intersect with the elements bounds.
	intersectingLeaves, _ := qt.leavesInRect(element.Bounds)

	for _, leaf := range intersectingLeaves {
		// Find the element node that contains the element.
		// Multiple leaves may contain one with a matching ID.
		if leaf.count == 0 {
			continue
		}

		lastIdx := int(leaf.firstChild)
		elemNode := qt.elementNodes.Get(lastIdx)
		prevNodeIdx := -1
		found := false
		for {
			if qt.elements.Get(int(elemNode.elementIdx)).ID == element.ID {
				found = true
				break
			}
			if elemNode.nextElemNode == -1 {
				break
			}
			prevNodeIdx = lastIdx
			lastIdx = int(elemNode.nextElemNode)
			elemNode = qt.elementNodes.Get(int(elemNode.nextElemNode))
		}

		// If the element was not found, continue with the next leaf.
		if !found {
			continue
		}

		// Remove the element node.
		qt.elements.Remove(int(elemNode.elementIdx))
		qt.elementNodes.Remove(lastIdx)
		if prevNodeIdx == -1 {
			leaf.firstChild = elemNode.nextElemNode
		} else {
			prevNode := qt.elementNodes.Get(prevNodeIdx)
			qt.elementNodes.Set(prevNodeIdx,
				QTElemNode{
					nextElemNode: elemNode.nextElemNode,
					elementIdx:   prevNode.elementIdx,
				})
		}
		leaf.count -= 1
	}
}

// Insert adds an element to the quadtree.
func (qt *QuadTree) Insert(element QTElement) {
	// Find all leaf nodes that intersect with the elements bounds.
	intersectingLeaves, depths := qt.leavesInRect(element.Bounds)

	for leafIdx, leaf := range intersectingLeaves {
		elemNode := QTElemNode{
			nextElemNode: leaf.firstChild,
			elementIdx:   int32(qt.elements.Insert(element)),
		}
		leaf.firstChild = int32(qt.elementNodes.Insert(elemNode))
		leaf.count += 1

		// If the leaf is above capacity, and below max depth, subdivide it.
		if depths[leafIdx] < qt.maxDepth && leaf.count > qt.leafCapacity {
			qt.subdivideLeaf(leaf)
		}
	}
}

// subdivideLeaf subdivides a leaf node into 4 children and distributes the
// elements to the children.
func (qt *QuadTree) subdivideLeaf(leaf *QTNode) {
	firstChild := leaf.firstChild

	// First, we remove all elements from the leaf.
	elements := make([]QTElement, 0)
	next := qt.elementNodes.Get(int(firstChild))
	for {
		element := qt.elements.Get(int(next.elementIdx))
		elements = append(elements, element)
		qt.Remove(element)
		// If there are no more elements, break.
		if next.nextElemNode == -1 {
			break
		}
		next = qt.elementNodes.Get(int(next.nextElemNode))
	}

	leaf.count = -1 // Mark this node as a branch node.

	newNodes := make([]*QTNode, 4)

	// Make space for the children.
	if qt.firstFreeNode == -1 {
		// -1 means there are no holes in the array
		leaf.firstChild = int32(len(qt.nodes))
		for i := 0; i < 4; i++ {
			newNode := &QTNode{firstChild: -1, count: 0} // Create a new empty node.
			qt.nodes = append(qt.nodes, newNode)
			newNodes[i] = newNode
		}
	} else {
		// reuse a hole in this case
		idx := qt.firstFreeNode
		qt.firstFreeNode = qt.nodes[idx].firstChild
		leaf.firstChild = idx
		for i := 0; i < 4; i++ {
			qt.nodes[idx] = &QTNode{firstChild: -1, count: 0} // Create a new empty node.
			newNodes[i] = qt.nodes[idx]
		}
	}

	// Insert the elements back into the children.
	for _, elem := range elements {
		qt.Insert(elem)
	}
}

// Draw draws the quadtree for debugging purposes, using raylib.
func (qt *QuadTree) Draw() {
	// NodeData is a helper struct that stores a node and its dynamically
	// calculated bounds.
	type NodeData struct {
		node  *QTNode
		bound QTBounds
	}

	toProcess := NewStack[NodeData](0)
	toProcess.Push(NodeData{node: qt.nodes[0], bound: qt.bounds}) // Push the root node.

	for toProcess.Size() > 0 {
		ndat, _ := toProcess.Pop()

		// Draw the bounds of the node.
		rl.DrawRectangleLines(ndat.bound.X, ndat.bound.Y, ndat.bound.Width, ndat.bound.Height, rl.Gray)

		// Draw the count of the node.
		rl.DrawText(
			fmt.Sprintf("%d", ndat.node.count),
			ndat.bound.X+ndat.bound.Width/2-10,
			ndat.bound.Y+ndat.bound.Height/2-10,
			10,
			rl.Black,
		)

		// If this node is a branch node, process its children.
		if ndat.node.count == -1 {
			// First, we calculate the AABBs of the children by subdividing the parent bounds.
			// This can be done quickly with bit manipulation.

			hw := ndat.bound.Width >> 1  // Half width.
			hh := ndat.bound.Height >> 1 // Half height.
			mx := ndat.bound.X + hw      // Horizontal mid.
			hy := ndat.bound.Y + hh      // Vertical mid.

			// Then we check which children intersect with the target
			// rectangle, and push them to the stack.

			toProcess.Push(NodeData{
				node:  qt.nodes[ndat.node.firstChild],
				bound: QTBounds{X: ndat.bound.X, Y: ndat.bound.Y, Width: hw, Height: hh},
			})
			toProcess.Push(NodeData{
				node:  qt.nodes[ndat.node.firstChild+1],
				bound: QTBounds{X: ndat.bound.X, Y: hy, Width: hw, Height: hh},
			})
			toProcess.Push(NodeData{
				node:  qt.nodes[ndat.node.firstChild+2],
				bound: QTBounds{X: mx, Y: ndat.bound.Y, Width: hw, Height: hh},
			})
			toProcess.Push(NodeData{
				node:  qt.nodes[ndat.node.firstChild+3],
				bound: QTBounds{X: mx, Y: hy, Width: hw, Height: hh},
			})
		}
	}
}
