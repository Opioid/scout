package scene

import (
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
)

type Scene struct {
	StaticProps []*StaticProp
	Lights []light.Light
	Complexes []Complex

	ComplexProvider ComplexProvider
}

func (scene *Scene) Init() {
	scene.ComplexProvider.Init()
}

func (scene *Scene) Compile() {
	/*
	for _, light := range scene.Lights {
		light.Entity.Transformation.Update()
	}*/
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

func (scene *Scene) AddLight(l light.Light) {
	scene.Lights = append(scene.Lights, l)
}

func (scene *Scene) CreateComplex(typename string) Complex {
	c := scene.ComplexProvider.NewComplex(typename)

	scene.Complexes = append(scene.Complexes, c)

	return c
}