package bvh

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Tree struct {
	root buildNode
	infiniteProps []*prop.StaticProp

//	nodes []miniNode
}

func (t *Tree) Intersect(ray *math.OptimizedRay, props []*prop.StaticProp, intersection *prop.Intersection) bool {
	hit := false

	for _, p := range t.infiniteProps {
		if p.Intersect(ray, intersection) {
			intersection.Prop = &p.Prop
			hit = true
		}
	}

	if t.root.intersect(ray, props, intersection) {
		hit = true
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
	for _, p := range t.infiniteProps {
		if p.IntersectP(ray) {
			return true
		}
	}

	return t.root.intersectP(ray, props)
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