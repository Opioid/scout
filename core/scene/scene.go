package scene

import (
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
)

type Scene struct {
	StaticProps []*StaticProp
	Lights []*light.Light

	Epsilon float32
}

func (scene *Scene) Compile() {
	for _, light := range scene.Lights {
		light.Entity.Transformation.Update()
	}
}

func (scene *Scene) Intersect(ray *math.Ray, intersection *Intersection) bool {
	hit := false
	var thit float32

	for _, prop := range scene.StaticProps {
		if prop.Intersect(ray, &thit) {
			ray.MaxT = thit
			intersection.Prop = &prop.Prop
			hit = true
		}
	}

	if hit {
		intersection.Dg.Point = ray.Origin.Add(ray.Direction.Scale(ray.MaxT - scene.Epsilon))
	}

	return hit
}

func (scene *Scene) IntersectP(ray *math.Ray) bool {
	hit := false
	var thit float32

	for _, prop := range scene.StaticProps {
		if prop.IntersectP(ray, &thit) {
			ray.MaxT = thit
			hit = true
		}
	}

	return hit
}

func (scene *Scene) CreateLight(lightType light.Type) *light.Light {
	l := &light.Light{Type: lightType}

	scene.Lights = append(scene.Lights, l)

	return l
}