package material

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/math32"
	gomath "math"
	_ "fmt"
)

const (
	// magic roughness constant that doesn't cause INF in specular_d
	// instead there is a max() now
	// 0.01313900625 
//	minRoughness = 0.0// 0.01313900625
	// Ran into another issue in specular_g, which doesn't (??) require such a high minRoughness. 
	// Mabye keeping it above 0 makes sense anyway. Don't know about the specular_d issue now.
	minRoughness = 1.0 / 255.0
)

func specular_f(v_dot_h float32, f0 math.Vector3) math.Vector3 {
	return f0.Add(math.MakeVector3(1.0 - f0.X, 1.0 - f0.Y, 1.0 - f0.Z).Scale(math.Exp2((-5.55473 * v_dot_h - 6.98316) * v_dot_h)))
}

func specular_d(n_dot_h, a2 float32) float32 {
	d := n_dot_h * n_dot_h * (a2 - 1.0) + 1.0
//	return a2 / math32.Max((gomath.Pi * d * d), gomath.SmallestNonzeroFloat32)
	return a2 / (gomath.Pi * d * d)
}

func specular_g(n_dot_l, n_dot_v, a2 float32) float32 {
	g_v := n_dot_v + math32.Sqrt((n_dot_v - n_dot_v * a2) * n_dot_v + a2)
	g_l := n_dot_l + math32.Sqrt((n_dot_l - n_dot_l * a2) * n_dot_l + a2)
	return math32.Rsqrt(g_v * g_l)
}

type SubstituteBrdf struct {
	Color   math.Vector3
	DiffuseColor math.Vector3
 	Opacity float32

 	N, v math.Vector3
 	N_dot_v float32

 	F0 math.Vector3
 	Roughness float32
 	a2 float32
}

func MakeSubstituteBrdf(color math.Vector3, opacity, roughness, metallic float32, n, v math.Vector3) SubstituteBrdf {
	brdf := SubstituteBrdf{}

	brdf.Color = color
	brdf.DiffuseColor = color.Scale(1.0 - metallic).Scale(opacity)
	brdf.Opacity = opacity
	brdf.N = n
	brdf.v = v
	brdf.N_dot_v = math32.Max(n.Dot(v), 0.0)

	brdf.F0 = math.MakeVector3(0.03, 0.03, 0.03).Lerp(color, metallic).Scale(opacity)
	
	brdf.Roughness = roughness
	a := roughness * roughness
	brdf.a2 = a * a

	return brdf
}

func (brdf *SubstituteBrdf) Evaluate(l math.Vector3) math.Vector3 {
	n_dot_l := math32.Max(brdf.N.Dot(l), 0.00001)

	h := brdf.v.Add(l).Normalized()

	n_dot_h := brdf.N.Dot(h)
	v_dot_h := brdf.v.Dot(h)

	specular := specular_f(v_dot_h, brdf.F0).Scale(specular_d(n_dot_h, brdf.a2)).Scale(specular_g(n_dot_l, brdf.N_dot_v, brdf.a2))

	r := brdf.DiffuseColor.Add(specular).Scale(n_dot_l)

	return r
}