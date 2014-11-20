package surrounding

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/ibl"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
	"os"
	"fmt"
)

type sphere struct {
	sphereMap texture.SamplerSphere

	ambientCube *light.AmbientCube

	diffuseSampler texture.SamplerSphere

	maxRoughnessMip float32
}

func NewSphere(sphericalTexture *texture.Texture2D) *sphere {
	s := new(sphere)

	s.sphereMap = texture.NewSamplerSpherical_linear(sphericalTexture)

//	s.ambientCube = NewAmbientCubeFromSurrounding(s)

	diffuse := texture.NewTexture2D(math.MakeVector2i(32, 16), 1)
//	diffuse := texture.NewTexture2D(math.MakeVector2i(256, 128), 1)

	ibl.CalculateSphereMapSolidAngleWeights(&sphericalTexture.Image.Buffers[0])

	ibl.IntegrateHemisphereSphereMap(s, 512, &diffuse.Image.Buffers[0])

	s.diffuseSampler = texture.NewSamplerSpherical_linear(diffuse) 

	numMipLevels := sphericalTexture.Image.NumMipLevels()

	if numMipLevels > 1 {
		s.maxRoughnessMip = float32(numMipLevels - 1)
		fmt.Println("We loaded the cache sometime previously.")
		return s
	}

	sphericalTexture.AllocateMipLevels(8)

	numMipLevels = sphericalTexture.Image.NumMipLevels()

	s.maxRoughnessMip = float32(numMipLevels - 1)

	roughnessIncrement := 1 / s.maxRoughnessMip

	for i := uint32(1); i < numMipLevels; i++ {
		ibl.IntegrateConeSphereMap(s, float32(i) * roughnessIncrement, 128, &sphericalTexture.Image.Buffers[i])
	}

	fo, err := os.Create("../cache/surrounding.sui")
	defer fo.Close()

	if err == nil {
		if err := texture.Save(fo, sphericalTexture); err != nil {
			fmt.Println(err)
		}
	}

	return s
}

func (s *sphere) Sample(ray *math.OptimizedRay) (math.Vector3, float32) {
	sample := s.sphereMap.Sample(ray.Direction)
	return sample.Vector3(), sample.W
}

func (s *sphere) SampleDiffuse(v math.Vector3) math.Vector3 {
//	return s.ambientCube.Evaluate(v)

	return s.diffuseSampler.Sample(v).Vector3()
}

func (s *sphere) SampleSpecular(v math.Vector3, roughness float32) math.Vector3 {
	return s.sphereMap.SampleLod(v, s.maxRoughnessMip * roughness).Vector3()
}

/*
func NewAmbientCubeFromSurrounding(surrounding Surrounding) *light.AmbientCube {
	ray := math.OptimizedRay{}
	ray.MaxT = 1000.0

	numSamples := uint32(256)
	numSamplesReciprocal := 1.0 / float32(numSamples)

	basis := math.Matrix3x3{}

	integrateHemisphere := func (d math.Vector3) math.Vector3 {
		basis.SetBasis(d)

		result := math.MakeVector3(0.0, 0.0, 0.0)

		for i := uint32(0); i < numSamples; i++ {
			sample := math.Hammersley(i, numSamples)

			s := math.HemisphereSample_cos(sample.X, sample.Y)

			v := basis.TransformVector3(s)
			ray.SetDirection(v)
			
			c, _ := surrounding.Sample(&ray)

			result.AddAssign(c.Scale(numSamplesReciprocal))
		}

		return result
	}

	ac := new(light.AmbientCube)

	ac.Colors[0] = integrateHemisphere(math.MakeVector3( 1.0,  0.0,  0.0))
	ac.Colors[1] = integrateHemisphere(math.MakeVector3(-1.0,  0.0,  0.0))
	ac.Colors[2] = integrateHemisphere(math.MakeVector3( 0.0,  1.0,  0.0))
	ac.Colors[3] = integrateHemisphere(math.MakeVector3( 0.0, -1.0,  0.0))
	ac.Colors[4] = integrateHemisphere(math.MakeVector3( 0.0,  0.0,  1.0))
	ac.Colors[5] = integrateHemisphere(math.MakeVector3( 0.0,  0.0, -1.0))

	return ac
}
*/