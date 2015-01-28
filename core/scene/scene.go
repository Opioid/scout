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

func (scene *Scene) Intersect(ray *math.OptimizedRay, visibility uint8, scratch *prop.ScratchBuffer, intersection *prop.Intersection) bool {
	hit := false

	if scene.bvh.Intersect(ray, visibility, scene.StaticProps, scratch, intersection) {
		hit = true
	}

	for _, p := range scene.DynamicProps {
		if p.Intersect(ray, scratch, &intersection.Geo) {
			intersection.Prop = p
			hit = true
		}
	}

	return hit
}

func (scene *Scene) IntersectP(ray *math.OptimizedRay, scratch *prop.ScratchBuffer) bool {
	if scene.bvh.IntersectP(ray, scene.StaticProps, scratch) {
		return true
	}

	for _, p := range scene.DynamicProps {
		if p.IntersectP(ray, scratch) {
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

	if p := l.Prop(); p.Shape != nil {
		p.CastsShadow = false
		scene.AddProp(p)
	}
}

func (scene *Scene) CreateComplex(typename string) Complex {
	c := scene.ComplexProvider.NewComplex(typename)

	if c != nil {
		scene.Complexes = append(scene.Complexes, c)
	}

	return c
}

func (scene *Scene) MonteCarloLight(r float32) (light.Light, float32) {
	numLights := len(scene.Lights)
	num := float32(numLights)
	l := int(num * r - 0.001)

	probability := 1.0 / num

	if l >= numLights {
		// the intention is that this symbolizes the surrounding light
		return nil, probability
	} else {
		return scene.Lights[l], probability
	}
}