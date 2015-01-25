package material 

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	_ "fmt"
)

type Values struct {
	Color   math.Vector3
	DiffuseColor math.Vector3

	N, V math.Vector3
	N_dot_v float32

	F0 math.Vector3
	Roughness float32
	A2 float32
}

func MakeValues(color math.Vector3, opacity, roughness, metallic float32, n, v math.Vector3) Values {
	values := Values{}

	values.Color = color
	values.DiffuseColor = color.Scale(1.0 - metallic)
	values.N = n
	values.V = v
	values.N_dot_v = math32.Max(n.Dot(v), 0.0)

	values.F0 = math.MakeVector3(0.03, 0.03, 0.03).Lerp(color, metallic)
	
	values.Roughness = roughness
	a := roughness * roughness
	values.A2 = a * a

	return values
}

func (values *Values) Set(color math.Vector3, opacity, roughness, metallic float32, n, v math.Vector3) {
	values.Color = color
	values.DiffuseColor = color.Scale(1.0 - metallic)
	values.N = n
	values.V = v
	values.N_dot_v = math32.Max(n.Dot(v), 0.0)

	values.F0 = math.MakeVector3(0.03, 0.03, 0.03).Lerp(color, metallic)
	
	values.Roughness = roughness
	a := roughness * roughness
	values.A2 = a * a
}