package light

import (
	pkgsurrounding "github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/base/math"
)

type AmbientCube struct {
	Colors [6]math.Vector3
}

func (a *AmbientCube) Evaluate(n math.Vector3) math.Vector3 {
/*
	vec3 normal_squared = normal_ws * normal_ws; 
	
	int is_negative_x = normal_ws.x < 0.f ? 1 : 0;
	int is_negative_y = normal_ws.y < 0.f ? 1 : 0;
	int is_negative_z = normal_ws.z < 0.f ? 1 : 0;
	
	vec3 color = normal_squared.x * ambient_cube[is_negative_x]     
		       + normal_squared.y * ambient_cube[is_negative_y + 2] 
		       + normal_squared.z * ambient_cube[is_negative_z + 4]; 
			  
    return color; 
    */

    n_squared := n.Mul(n)

    is_negative_x, is_negative_y, is_negative_z := 0, 0, 0

    if n.X < 0.0 {
    	is_negative_x = 1
    }

    if n.Y < 0.0 {
    	is_negative_y = 1
    }

    if n.Z < 0.0 {
    	is_negative_z = 1
    }

    color := a.Colors[is_negative_x    ].Scale(n_squared.X).Add(
    		 a.Colors[is_negative_y + 2].Scale(n_squared.Y).Add(
    		 a.Colors[is_negative_z + 4].Scale(n_squared.Z)))

    return color
}


func NewAmbientCubeFromSurrounding(surrounding pkgsurrounding.Surrounding) *AmbientCube {
	ac := new(AmbientCube)

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
			
			c := surrounding.Sample(&ray)

			result.AddAssign(c.Scale(numSamplesReciprocal))
		}

		return result
	}

	ac.Colors[0] = integrateHemisphere(math.MakeVector3( 1.0,  0.0,  0.0))
	ac.Colors[1] = integrateHemisphere(math.MakeVector3(-1.0,  0.0,  0.0))
	ac.Colors[2] = integrateHemisphere(math.MakeVector3( 0.0,  1.0,  0.0))
	ac.Colors[3] = integrateHemisphere(math.MakeVector3( 0.0, -1.0,  0.0))
	ac.Colors[4] = integrateHemisphere(math.MakeVector3( 0.0,  0.0,  1.0))
	ac.Colors[5] = integrateHemisphere(math.MakeVector3( 0.0,  0.0, -1.0))

	return ac
}