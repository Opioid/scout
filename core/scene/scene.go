package scene

import (
	"github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/core/scene/bvh"
	_ "github.com/Opioid/scout/core/scene/entity"
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/core/scene/light"
	_ "github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/base/math"
	_ "fmt"
)

type Scene struct {
	Surrounding surrounding.Surrounding

	bvh bvh.Tree

	StaticProps []*prop.Prop
	DynamicProps []*prop.Prop

	Lights []light.Light

	Complexes []Complex

	ComplexProvider ComplexProvider
}

func (scene *Scene) Init() {
	scene.ComplexProvider.Init()
}

func (scene *Scene) Compile() {
	builder := bvh.Builder{}
	var outProps []*prop.Prop
	builder.Build(scene.StaticProps, 4, &scene.bvh, &outProps)

	scene.StaticProps = outProps

	/*
	buffer := texture.Buffer{}
	buffer.Resize(math.MakeVector2i(512, 256))
	surrounding.BakeSphereMap(scene.Surrounding, &buffer)
	*/
}

func (scene *Scene) Intersect(ray *math.OptimizedRay, intersection *prop.Intersection) bool {
	hit := false

	if scene.bvh.Intersect(ray, scene.StaticProps, intersection) {
		hit = true
	}

	for _, p := range scene.DynamicProps {
		if p.Intersect(ray, &intersection.Geo) {
			intersection.Prop = p
			hit = true
		}
	}

	return hit
}

func (scene *Scene) IntersectP(ray *math.OptimizedRay) bool {
	if scene.bvh.IntersectP(ray, scene.StaticProps) {
		return true
	}

	for _, p := range scene.DynamicProps {
		if p.IntersectP(ray) {
			return true
		}
	}

	return false
}

func (scene *Scene) AddProp(p *prop.Prop) {
	scene.StaticProps = append(scene.StaticProps, p)
}

func (scene *Scene) AddLight(l light.Light) {
	scene.Lights = append(scene.Lights, l)
/*
	if p := l.Prop(); p.Shape != nil {
		p.CastsShadow = false
		scene.AddProp(p)
	}*/
}

func (scene *Scene) CreateComplex(typename string) Complex {
	c := scene.ComplexProvider.NewComplex(typename)
	scene.Complexes = append(scene.Complexes, c)
	return c
}