package bvh

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Tree struct {
	root node
	props []*prop.StaticProp
}

func (t *Tree) Assign(props []*prop.StaticProp) {

	numFinitShapes := 0

	for _, p := range props {
		if p.Shape.IsFinite() {
			numFinitShapes++
		}
	}

	t.props = make([]*prop.StaticProp, len(props) - numFinitShapes)

	i := 0
	for _, p := range props {
		if !p.Shape.IsFinite() {
			t.props[i] = p
			i++
		}
	}

	finiteProps := make([]*prop.StaticProp, numFinitShapes)

	i = 0
	for _, p := range props {
		if p.Shape.IsFinite() {
			finiteProps[i] = p
			i++
		}
	}

	t.root.split(finiteProps)
}

func (t *Tree) Intersect(ray *math.OptimizedRay, intersection *prop.Intersection) bool {
	hit := false

	for _, p := range t.props {
		if p.Intersect(ray, intersection) {
			intersection.Prop = &p.Prop
			hit = true
		}
	}

	if t.root.intersect(ray, intersection) {
		hit = true
	}

	return hit
}

func (t *Tree) IntersectP(ray *math.OptimizedRay) bool {
	for _, p := range t.props {
		if p.IntersectP(ray) {
			return true
		}
	}

	return t.root.intersectP(ray)
}
