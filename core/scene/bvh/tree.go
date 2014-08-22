package bvh

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Tree struct {
	staticProps []*prop.StaticProp
}

func (t *Tree) Split(staticProps []*prop.StaticProp) {

	numFinitShapes := 0

	for _, p := range staticProps {
		if p.Shape.IsFinite() {
			numFinitShapes++
		}
	}

	t.staticProps = make([]*prop.StaticProp, len(staticProps) - numFinitShapes)

	i := 0
	for _, p := range staticProps {
		if !p.Shape.IsFinite() {
			t.staticProps[i] = p
			i++
		}
	}

//	t.staticProps = staticProps

}

func (t *Tree) Intersect(ray *math.Ray, intersection *prop.Intersection) bool {
	hit := false

	for _, prop := range t.staticProps {
		if prop.Intersect(ray, intersection) {
			intersection.Prop = &prop.Prop
			hit = true
		}
	}

	return hit
}

func (t *Tree) IntersectP(ray *math.Ray) bool {
	for _, prop := range t.staticProps {
		if prop.IntersectP(ray) {
			return true
		}
	}

	return false
}
