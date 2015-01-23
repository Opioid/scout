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

func (t *Tree) Intersect(ray *math.OptimizedRay, props []*prop.Prop, transformation *math.ComposedTransformation, intersection *prop.Intersection) bool {
	hit := false

	if t.root.intersect(ray, props, transformation, intersection) {
		hit = true
	}

	for i := t.infinitePropsBegin; i < t.infinitePropsEnd; i++ {
		p := props[i]
		if p.Intersect(ray, transformation, &intersection.Geo) {
			intersection.Prop = p
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


func (t *Tree) IntersectP(ray *math.OptimizedRay, props []*prop.Prop, transformation *math.ComposedTransformation) bool {
	if t.root.intersectP(ray, props, transformation) {
		return true
	}

	for i := t.infinitePropsBegin; i < t.infinitePropsEnd; i++ {
		if  props[i].CastsShadow && props[i].IntersectP(ray, transformation) {
			return true
		}
	}	

	return false

}