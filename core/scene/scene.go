package scene

import (
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
)

type Scene struct {
	StaticProps []*StaticProp
	Lights []*light.Light
}

func (scene *Scene) Compile() {
	for _, light := range scene.Lights {
		light.Entity.Transformation.Update()
	}
}

func (scene *Scene) Intersect(ray *math.Ray, intersection *Intersection) bool {
	hit := false

	for _, prop := range scene.StaticProps {
		if prop.Intersect(ray, intersection) {
			intersection.Prop = &prop.Prop
			hit = true
		}
	}

	return hit
}

func (scene *Scene) IntersectP(ray *math.Ray) bool {
	for _, prop := range scene.StaticProps {
		if prop.IntersectP(ray) {
			return true
		}
	}

	return false
}

func (scene *Scene) CreateLight(lightType light.Type) *light.Light {
	l := &light.Light{Type: lightType}

	scene.Lights = append(scene.Lights, l)

	return l
}