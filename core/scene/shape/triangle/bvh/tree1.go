package bvh

import (
	_ "github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	_ "math"
	_ "fmt"
)

const hasChildrenFlag uint32 = 0xFFFFFFFF

type Tree1 struct {
	nodes []node

	Triangles []primitive.Triangle
}

func (t *Tree1) AABB() bounding.AABB {
	return t.nodes[0].aabb
}

func (t *Tree1) Intersect(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, intersection *primitive.Intersection) bool {
	return t.intersectNode(0, ray, intersection)
}

func (t *Tree1) IntersectP(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32) bool {
	return t.intersectNodeP(0, ray)
}

func (t *Tree1) intersectNode(n uint32, ray *math.OptimizedRay, intersection *primitive.Intersection) bool {
	node := &t.nodes[n]

	if !node.aabb.IntersectP(ray) {
		return false
	}

	hit := false

	if node.hasChildren() {
		a, b := node.children(ray.Sign[node.axis], n)

		if t.intersectNode(a, ray, intersection) {
			hit = true
		} 

		if t.intersectNode(b, ray, intersection) {
			hit = true
		}
	} else {
		ti := primitive.Intersection{}
		ti.T = ray.MaxT

		for i := node.startIndex; i < node.endIndex; i++ {
			if h, c := t.Triangles[i].Intersect(ray); h && c.T < ti.T {
				ti.Coordinates = c
				ti.Index = i
				hit = true
			}
		}

		if hit {
			// the idea was not to write to these pointers in the loop... Don't know whether it makes a difference
			*intersection = ti
			ray.MaxT = ti.T
		}
	}

	return hit
}

func (t *Tree1) intersectNodeP(n uint32, ray *math.OptimizedRay) bool {
	node := &t.nodes[n]

	if !node.aabb.IntersectP(ray) {
		return false
	}

	if node.hasChildren() {
		a, b := node.children(ray.Sign[node.axis], n)

		if t.intersectNodeP(a, ray) {
			return true
		}

		return t.intersectNodeP(b, ray)
	}

	for i := node.startIndex; i < node.endIndex; i++ {
		if t.Triangles[i].IntersectP(ray) {
			return true
		}
	}

	return false

}

func (t *Tree1) allocateNodes(numNodes uint32) []node {
	if uint32(len(t.nodes)) < numNodes {
		t.nodes = make([]node, numNodes)
	}

	return t.nodes
}

type node struct {
	aabb bounding.AABB

	startIndex, endIndex uint32

	axis int8
}

func (n *node) hasChildren() bool {
	return (n.startIndex & hasChildrenFlag) == hasChildrenFlag
}

func (n *node) setHasChildren(children bool) {
	if children {
		n.startIndex |= hasChildrenFlag
	} else {
		n.startIndex &= ^hasChildrenFlag
	}
}

func (n *node) children(sign uint32, id uint32) (uint32, uint32) {
	if sign == 0 {
		return id + 1, n.endIndex
	} else {
		return n.endIndex, id + 1
	}
}

func (n *node) setRightChild(offset uint32) {
	n.endIndex = offset
}