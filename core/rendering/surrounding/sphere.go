package surrounding

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/rendering/ibl"
	"github.com/Opioid/scout/core/scene/light"
	"github.com/Opioid/scout/base/math"
)

const (
	numSamples = 128
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

	ibl.CalculateSphereMapSolidAngleWeights(&sphericalTexture.Image.Buffers[0])

	ibl.IntegrateHemisphereSphereMap(s, numSamples, &diffuse.Image.Buffers[0])

	s.diffuseSampler = texture.NewSamplerSpherical_linear(diffuse) 

	sphericalTexture.AllocateMipLevelsDownTo(math.MakeVector2i(20, 10))

	// UGLY: have to reset the texture, so that the sampler registers the additional mip maps
	s.sphereMap.SetTexture(sphericalTexture)

	numMipLevels := sphericalTexture.Image.NumMipLevels()

	s.maxRoughnessMip = float32(numMipLevels - 1)

	roughnessIncrement := 1 / s.maxRoughnessMip

	for i := uint32(1); i < numMipLevels; i++ {
		ibl.IntegrateConeSphereMap(s, float32(i) * roughnessIncrement, numSamples, &sphericalTexture.Image.Buffers[i])
	}

	return s
}

func NewSphereFromCache(diffuseTexture *texture.Texture2D, specularTexture *texture.Texture2D) *sphere {
	s := new(sphere)

	s.diffuseSampler = texture.NewSamplerSpherical_linear(diffuseTexture) 
	s.sphereMap = texture.NewSamplerSpherical_linear(specularTexture)

	numMipLevels := specularTexture.Image.NumMipLevels()
	s.maxRoughnessMip = float32(numMipLevels - 1)

	return s
}

func (s *sphere) DiffuseTexture() *texture.Texture2D {
	return s.diffuseSampler.Texture()
}

func (s *sphere) SpecularTexture() *texture.Texture2D {
	return s.sphereMap.Texture()
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