package scene

import (
	"github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/core/scene/bvh"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
)

type Scene struct {
	Surrounding surrounding.Surrounding

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
	builder := bvh.Builder{}
	var outProps []*prop.StaticProp
	builder.Build(scene.StaticProps, 4, &scene.bvh, &outProps)

	scene.StaticProps = outProps
}

func (scene *Scene) Intersect(ray *math.OptimizedRay, intersection *prop.Intersection) bool {
	return scene.bvh.Intersect(ray, scene.StaticProps, intersection)
}

func (scene *Scene) IntersectP(ray *math.OptimizedRay) bool {
	return scene.bvh.IntersectP(ray, scene.StaticProps)
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