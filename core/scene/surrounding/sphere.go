package surrounding

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
)

type sphere struct {
	sphereMap texture.SamplerSphere

	ambientCube *light.AmbientCube

	diffuseSampler texture.SamplerSphere
}

func NewSphere(sphericalTexture *texture.Texture2D) *sphere {
	s := new(sphere)

	s.sphereMap = texture.NewSamplerSpherical_linear(sphericalTexture)

	s.ambientCube = NewAmbientCubeFromSurrounding(s)

	diffuse := texture.NewTexture2D(math.MakeVector2i(32, 16), 1)
//	diffuse := texture.NewTexture2D(math.MakeVector2i(256, 128), 1)

	calculateSphereMapSolidAngleWeights(&sphericalTexture.Image.Buffers[0])

	integrateHemisphereSphereMap(s, 512, &diffuse.Image.Buffers[0])

	s.diffuseSampler = texture.NewSamplerSpherical_linear(diffuse) 

	sphericalTexture := s.sphereMap.Texture()

	sphericalTexture.AllocateMipLevels(4)

	mipLevels := sphericalTexture.Image.MipLevels()

	roughnessIncrement := 1.0 / float32(mipLevels - 1) 

	for i := 1; i < mipLevels; i++ {
		integrateConeSphereMap(s, float32(i) * roughnessIncrement, 32, &sphericalTexture.Image.Buffers[i])
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

func (s *sphere) SampleSpecular(v math.Vector3) math.Vector3 {
	return s.sphereMap.Sample(v).Vector3()
}

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