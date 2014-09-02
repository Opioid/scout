package bvh

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Tree struct {
	root buildNode

	infinitePropsBegin uint32
	infinitePropsEnd   uint32
}

func (t *Tree) Intersect(ray *math.OptimizedRay, props []*prop.StaticProp, intersection *prop.Intersection) bool {
	hit := false

	if t.root.intersect(ray, props, intersection) {
		hit = true
	}

	for i := t.infinitePropsBegin; i < t.infinitePropsEnd; i++ {
		p := props[i]
		if p.Intersect(ray, intersection) {
			intersection.Prop = &p.Prop
			hit = true
		}
	}	

/*
	currentNode := uint32(0)
	n := &t.nodes[currentNode]

	for n != nil {
		if !n.aabb.Intersect(ray) {
			currentNode = n.skipOffset()
			if currentNode == 0 {
				return hit
			} else {
				n = &t.nodes[currentNode]
			}

			continue
		}

		if !n.hasChildren() {
			for _, i := range n.indices {
				p := props[i]
				if p.Intersect(ray, intersection) {
					intersection.Prop = &p.Prop
					hit = true
				}
			}

			currentNode = n.skipOffset()
			if currentNode == 0 {
				return hit
			} else {
				n = &t.nodes[currentNode]
			}

		} else {
			currentNode++
			n = &t.nodes[currentNode]
		}
	}
*/
	return hit
}

func (t *Tree) IntersectP(ray *math.OptimizedRay, props []*prop.StaticProp) bool {
	if t.root.intersectP(ray, props) {
		return true
	}

	for i := t.infinitePropsBegin; i < t.infinitePropsEnd; i++ {
		if props[i].IntersectP(ray) {
			return true
		}
	}	

	return false
/*
	currentNode := uint32(0)
	n := &t.nodes[currentNode]

	for n != nil {
		if !n.aabb.Intersect(ray) {
			currentNode = n.skipOffset()
			if currentNode == 0 {
				return false
			} else {
				n = &t.nodes[currentNode]
			}

			continue
		}

		if !n.hasChildren() {
			for _, i := range n.indices {
				if props[i].IntersectP(ray) {
					return true
				}
			}

			currentNode = n.skipOffset()
			if currentNode == 0 {
				return false
			} else {
				n = &t.nodes[currentNode]
			}

		} else {
			currentNode++
			n = &t.nodes[currentNode]
		}
	}

	return false
	*/
}

/*
func (t *Tree) allocateNodes(numNodes uint32) []miniNode {
	if uint32(len(t.nodes)) < numNodes {
		t.nodes = make([]miniNode, numNodes)
	}

	return t.nodes
}
*/