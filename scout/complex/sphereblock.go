package complex

import (
	pkgscene "github.com/Opioid/scout/core/scene"
	"github.com/Opioid/scout/core/resource"
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/core/scene/material"
	"github.com/Opioid/scout/base/math"
	gomath "math"
)

type sphereBlock struct {

}

func (c *sphereBlock) Init(scene *pkgscene.Scene, resourceManager *resource.Manager) {
	shape := shape.NewSphere()

	materials := make([]material.Material, 1)
	materials[0] = resourceManager.LoadMaterial("../data/materials/green.material")

	scale := math.MakeVector3(0.5, 0.5, 0.5)
	spacing := math.MakeVector3(1.2, 1.2, 1.2)

	numX, numY, numZ := 16, 10, 16

	offset := math.MakeVector3(-0.5 * float32(numX), 0.006, 0.0).Add(scale.Mul(spacing))

	for z := 0; z < numZ; z++ {
		for y := 0; y < numY; y++ {
			for x := 0; x < numX; x++ {
				p := scene.CreateProp()
				p.Shape = shape
				p.Materials = materials

				position := math.MakeVector3(float32(x), float32(y), float32(z)).Mul(spacing)

				modifier := math.MakeVector3(
					-math.Sin(float32(y) / float32(numY - 1) * 1.0 * gomath.Pi), 
					 math.Sin(float32(z) / float32(numZ - 1) * gomath.Pi), 
					-math.Sin(float32(x) / float32(numX - 1) * 2.0 * gomath.Pi))

				p.SetTransformation(offset.Add(position.Add(modifier)), scale, math.MakeIdentityQuaternion())
			}
		}
	}


}

type SphereBlockFactory struct {

}

func (f *SphereBlockFactory) New() pkgscene.Complex {
	return new(sphereBlock)
}