package scene

import (
	"github.com/Opioid/scout/core/scene/bvh"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
)

type Scene struct {
	bvh bvh.Tree

	StaticProps []*prop.StaticProp
	Lights []light.Light
	Complexes []Complex

	ComplexProvider ComplexProvider
}

func (scene *Scene) Init() {
	scene.ComplexProvider.Init()
}

func (scene *Scene) Compile() {
	scene.bvh.Split(scene.StaticProps)
}

func (scene *Scene) Intersect(ray *math.Ray, intersection *prop.Intersection) bool {
/*	hit := false

	for _, prop := range scene.StaticProps {
		if prop.Intersect(ray, intersection) {
			intersection.Prop = &prop.Prop
			hit = true
		}
	}

	return hit
	*/
	return scene.bvh.Intersect(ray, intersection)
}

func (scene *Scene) IntersectP(ray *math.Ray) bool {
/*	for _, prop := range scene.StaticProps {
		if prop.IntersectP(ray) {
			return true
		}
	}

	return false
	*/
	return scene.bvh.IntersectP(ray)
}

func (scene *Scene) CreateStaticProp() *prop.StaticProp {
	p := prop.NewStaticProp()

	scene.StaticProps = append(scene.StaticProps, p)

	return p
}

func (scene *Scene) AddLight(l light.Light) {
	scene.Lights = append(scene.Lights, l)
}

func (scene *Scene) CreateComplex(typename string) Complex {
	c := scene.ComplexProvider.NewComplex(typename)

	scene.Complexes = append(scene.Complexes, c)

	return c
}