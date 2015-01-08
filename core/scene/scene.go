package scene

import (
	"github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/core/scene/bvh"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/light"
	_ "github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Scene struct {
	Surrounding surrounding.Surrounding

	bvh bvh.Tree

	StaticProps []*prop.StaticProp

	Actors []*prop.Actor
	
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

	/*
	buffer := texture.Buffer{}
	buffer.Resize(math.MakeVector2i(512, 256))
	surrounding.BakeSphereMap(scene.Surrounding, &buffer)
	*/
}

func (scene *Scene) Intersect(ray *math.OptimizedRay, intersection *prop.Intersection) bool {
	hit := scene.bvh.Intersect(ray, scene.StaticProps, intersection)

	for _, a := range scene.Actors {
		if a.Intersect(ray, intersection) {
			intersection.Prop = &a.Prop
			hit = true
		}
	}

	return hit
}

func (scene *Scene) IntersectP(ray *math.OptimizedRay) bool {
	if scene.bvh.IntersectP(ray, scene.StaticProps) {
		return true
	}

	for _, a := range scene.Actors {
		if a.IntersectP(ray) {
			return true
		}
	}

	return false
}

func (scene *Scene) CreateStaticProp() *prop.StaticProp {
	p := prop.NewStaticProp()
	scene.StaticProps = append(scene.StaticProps, p)
	return p
}

func (scene *Scene) AddLight(l light.Light) {
	scene.Lights = append(scene.Lights, l)
}

func (scene *Scene) CreateActor() *prop.Actor {
	a := prop.NewActor()
	scene.Actors = append(scene.Actors, a)
	return a
}

func (scene *Scene) CreateComplex(typename string) Complex {
	c := scene.ComplexProvider.NewComplex(typename)
	scene.Complexes = append(scene.Complexes, c)
	return c
}