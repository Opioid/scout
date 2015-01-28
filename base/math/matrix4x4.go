package math

import (
	_ "fmt"
)

/*
[ 0]:m00	[ 1]:m01	[ 2]:m02	[ 3]:m03
[ 4]:m10	[ 5]:m11	[ 6]:m12	[ 7]:m13
[ 8]:m20	[ 9]:m21	[10]:m22	[11]:m23
[12]:m30	[13]:m31	[14]:m32	[15]:m33
*/

type Matrix4x4 struct {
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float32
}

func MakeMatrix4x4FromBasisScaleOrigin(b *Matrix3x3, s, o Vector3) Matrix4x4 {
	return Matrix4x4{
		b.m00 * s.X, b.m01 * s.X, b.m02 * s.X, 0.0,
		b.m10 * s.Y, b.m11 * s.Y, b.m12 * s.Y, 0.0,
		b.m20 * s.Z, b.m21 * s.Z, b.m22 * s.Z, 0.0,
		o.X,         o.Y,         o.Z,         1.0,
	}
}

func (m *Matrix4x4) Row(i int) Vector4 {
//	return MakeVector4(m.m[i * 4], m.m[i * 4 + 1], m.m[i * 4 + 2], m.m[i * 4 + 3])
	switch i {
	case 0:
		return MakeVector4(m.m00, m.m01, m.m02, m.m03)
	case 1:
		return MakeVector4(m.m10, m.m11, m.m12, m.m13)
	case 2:
		return MakeVector4(m.m20, m.m21, m.m22, m.m23)
	default:
		return MakeVector4(m.m30, m.m31, m.m32, m.m33)
	}		
}

func (m *Matrix4x4) Right() Vector3 {
	return MakeVector3(m.m00, m.m01, m.m02)
}

func (m *Matrix4x4) Up() Vector3 {
	return MakeVector3(m.m10, m.m11, m.m12)
}

func (m *Matrix4x4) Direction() Vector3 {
	return MakeVector3(m.m20, m.m21, m.m22)
}

func (m *Matrix4x4) Translation() Vector3 {
	return MakeVector3(m.m30, m.m31, m.m32)
}

func (m *Matrix4x4) Div(s float32) Matrix4x4 {
	is := 1.0 / s
	return Matrix4x4{
		m.m00 * is, m.m01 * is, m.m02 * is, m.m03 * is,
		m.m10 * is, m.m11 * is, m.m12 * is, m.m13 * is,
		m.m20 * is, m.m21 * is, m.m22 * is, m.m23 * is,
		m.m30 * is, m.m31 * is, m.m32 * is, m.m33 * is,
	}
}

func (m *Matrix4x4) SetIdentity() {
	m.m00 = 1.0; m.m01 = 0.0; m.m02 = 0.0; m.m03 = 0.0
	m.m10 = 0.0; m.m11 = 1.0; m.m12 = 0.0; m.m13 = 0.0
	m.m20 = 0.0; m.m21 = 0.0; m.m22 = 1.0; m.m23 = 0.0
	m.m30 = 0.0; m.m31 = 0.0; m.m32 = 0.0; m.m33 = 1.0
}

func (m *Matrix4x4) SetFromBasisScaleOrigin(b *Matrix3x3, s, o Vector3) {
	m.m00 = b.m00 * s.X; m.m01 = b.m01 * s.X; m.m02 = b.m02 * s.X; m.m03 = 0.0;
	m.m10 = b.m10 * s.Y; m.m11 = b.m11 * s.Y; m.m12 = b.m12 * s.Y; m.m13 = 0.0;
	m.m20 = b.m20 * s.Z; m.m21 = b.m21 * s.Z; m.m22 = b.m22 * s.Z; m.m23 = 0.0;
	m.m30 = o.X;         m.m31 = o.Y;         m.m32 = o.Z;         m.m33 = 1.0;
}

func (m *Matrix4x4) SetBasis(b *Matrix3x3) {
	m.m00 = b.m00; m.m01 = b.m01; m.m02 = b.m02
	m.m10 = b.m10; m.m11 = b.m11; m.m12 = b.m12
	m.m20 = b.m20; m.m21 = b.m21; m.m22 = b.m22
}

func (m *Matrix4x4) SetOrigin(v Vector3) {
	m.m30 = v.X
	m.m31 = v.Y
	m.m32 = v.Z
}

func (m *Matrix4x4) Scale(v Vector3) {
	m.m00 *= v.X; m.m01 *= v.X; m.m02 *= v.X
	m.m10 *= v.Y; m.m11 *= v.Y; m.m12 *= v.Y
	m.m20 *= v.Z; m.m21 *= v.Z; m.m22 *= v.Z
}

func (m *Matrix4x4) TransformPoint(v Vector3) Vector3 {
	return MakeVector3(
		v.X * m.m00 + v.Y * m.m10 + v.Z * m.m20 + m.m30,
		v.X * m.m01 + v.Y * m.m11 + v.Z * m.m21 + m.m31,
		v.X * m.m02 + v.Y * m.m12 + v.Z * m.m22 + m.m32,
	)
}

func (m *Matrix4x4) TransformVector3(v Vector3) Vector3 {
	return MakeVector3(
		v.X * m.m00 + v.Y * m.m10 + v.Z * m.m20,
		v.X * m.m01 + v.Y * m.m11 + v.Z * m.m21,
		v.X * m.m02 + v.Y * m.m12 + v.Z * m.m22,
	)
}

func (m *Matrix4x4) TransposedTransformVector3(v Vector3) Vector3 {
	return MakeVector3(
		v.X * m.m00 + v.Y * m.m01 + v.Z * m.m02,
		v.X * m.m10 + v.Y * m.m11 + v.Z * m.m12,
		v.X * m.m20 + v.Y * m.m21 + v.Z * m.m22,
	)
}

func (m *Matrix4x4) Det() float32 {
	return m.m00 * m.m11 * m.m22 * m.m33 + m.m00 * m.m12 * m.m23 * m.m31 + m.m00 * m.m13 * m.m21 * m.m32 +
		   m.m01 * m.m10 * m.m23 * m.m32 + m.m01 * m.m12 * m.m20 * m.m33 + m.m01 * m.m13 * m.m22 * m.m30 +
		   m.m02 * m.m10 * m.m21 * m.m33 + m.m02 * m.m11 * m.m23 * m.m30 + m.m02 * m.m13 * m.m20 * m.m31 +
		   m.m03 * m.m10 * m.m22 * m.m31 + m.m03 * m.m11 * m.m21 * m.m32 + m.m03 * m.m12 * m.m21 * m.m30 -

		   m.m00 * m.m11 * m.m23 * m.m32 - m.m00 * m.m12 * m.m21 * m.m33 - m.m00 * m.m13 * m.m22 * m.m31 -
		   m.m01 * m.m10 * m.m22 * m.m33 - m.m01 * m.m12 * m.m23 * m.m30 - m.m01 * m.m13 * m.m20 * m.m32 -
		   m.m02 * m.m10 * m.m23 * m.m31 - m.m02 * m.m11 * m.m20 * m.m33 - m.m02 * m.m13 * m.m21 * m.m30 -
		   m.m03 * m.m10 * m.m21 * m.m32 - m.m03 * m.m11 * m.m22 * m.m30 - m.m03 * m.m12 * m.m20 * m.m31
}

// Unnecessarily slow
/*
func (m *Matrix4x4) Inverted() Matrix4x4 {
	idet := 1.0 / m.Det()
	return Matrix4x4{
		(m.m11*m.m22*m.m33 + m.m12*m.m23*m.m31 + m.m13*m.m21*m.m32 - m.m11*m.m23*m.m32 - m.m12*m.m21*m.m33 - m.m13*m.m22*m.m31) * idet,
		(m.m01*m.m23*m.m32 + m.m02*m.m21*m.m33 + m.m03*m.m22*m.m31 - m.m01*m.m22*m.m33 - m.m02*m.m23*m.m31 - m.m03*m.m21*m.m32) * idet,
		(m.m01*m.m12*m.m33 + m.m02*m.m13*m.m31 + m.m03*m.m11*m.m32 - m.m01*m.m13*m.m32 - m.m02*m.m11*m.m33 - m.m03*m.m12*m.m31) * idet,
		(m.m01*m.m13*m.m22 + m.m02*m.m11*m.m23 + m.m03*m.m12*m.m21 - m.m01*m.m12*m.m23 - m.m02*m.m13*m.m21 - m.m03*m.m11*m.m22) * idet,

		(m.m10*m.m23*m.m32 + m.m12*m.m20*m.m33 + m.m13*m.m22*m.m30 - m.m10*m.m22*m.m33 - m.m12*m.m23*m.m30 - m.m13*m.m20*m.m32) * idet,
		(m.m00*m.m22*m.m33 + m.m02*m.m23*m.m30 + m.m03*m.m20*m.m32 - m.m00*m.m23*m.m32 - m.m02*m.m20*m.m33 - m.m03*m.m22*m.m30) * idet,
		(m.m00*m.m13*m.m32 + m.m02*m.m10*m.m33 + m.m03*m.m12*m.m30 - m.m00*m.m12*m.m33 - m.m02*m.m13*m.m30 - m.m03*m.m10*m.m32) * idet,
		(m.m00*m.m12*m.m23 + m.m02*m.m13*m.m20 + m.m03*m.m10*m.m22 - m.m00*m.m13*m.m22 - m.m02*m.m10*m.m23 - m.m03*m.m12*m.m20) * idet,

		(m.m10*m.m21*m.m33 + m.m11*m.m23*m.m30 + m.m13*m.m20*m.m31 - m.m10*m.m23*m.m31 - m.m11*m.m20*m.m33 - m.m13*m.m21*m.m30) * idet,
		(m.m00*m.m23*m.m31 + m.m01*m.m20*m.m33 + m.m03*m.m21*m.m30 - m.m00*m.m21*m.m33 - m.m01*m.m23*m.m30 - m.m03*m.m20*m.m31) * idet,
		(m.m00*m.m11*m.m33 + m.m01*m.m13*m.m30 + m.m03*m.m10*m.m31 - m.m00*m.m13*m.m31 - m.m01*m.m10*m.m33 - m.m03*m.m11*m.m30) * idet,
		(m.m00*m.m13*m.m21 + m.m01*m.m10*m.m23 + m.m03*m.m11*m.m20 - m.m00*m.m11*m.m23 - m.m01*m.m13*m.m20 - m.m03*m.m10*m.m21) * idet,

		(m.m10*m.m22*m.m31 + m.m11*m.m20*m.m32 + m.m12*m.m21*m.m30 - m.m10*m.m21*m.m32 - m.m11*m.m22*m.m30 - m.m12*m.m20*m.m31) * idet,
		(m.m00*m.m21*m.m32 + m.m01*m.m22*m.m30 + m.m02*m.m20*m.m31 - m.m00*m.m22*m.m31 - m.m01*m.m20*m.m32 - m.m02*m.m21*m.m30) * idet,
		(m.m00*m.m12*m.m31 + m.m01*m.m10*m.m32 + m.m02*m.m11*m.m30 - m.m00*m.m11*m.m32 - m.m01*m.m12*m.m30 - m.m02*m.m10*m.m31) * idet,
		(m.m00*m.m11*m.m22 + m.m01*m.m12*m.m20 + m.m02*m.m10*m.m21 - m.m00*m.m12*m.m21 - m.m01*m.m10*m.m22 - m.m02*m.m11*m.m20) * idet,
	}
}
*/

// This gives the same result as the above and is a bit faster
// http://rodolphe-vaillant.fr/?e=7
func (m *Matrix4x4) Inverted() Matrix4x4 {
	inv00 :=  m.m11 * m.m22 * m.m33 - m.m11 * m.m23 * m.m32 - m.m21 * m.m12 * m.m33 + m.m21 * m.m13 * m.m32 + m.m31 * m.m12 * m.m23 - m.m31 * m.m13 * m.m22
	inv04 := -m.m10 * m.m22 * m.m33 + m.m10 * m.m23 * m.m32 + m.m20 * m.m12 * m.m33 - m.m20 * m.m13 * m.m32 - m.m30 * m.m12 * m.m23 + m.m30 * m.m13 * m.m22
	inv08 :=  m.m10 * m.m21 * m.m33 - m.m10 * m.m23 * m.m31 - m.m20 * m.m11 * m.m33 + m.m20 * m.m13 * m.m31 + m.m30 * m.m11 * m.m23 - m.m30 * m.m13 * m.m21
	inv12 := -m.m10 * m.m21 * m.m32 + m.m10 * m.m22 * m.m31 + m.m20 * m.m11 * m.m32 - m.m20 * m.m12 * m.m31 - m.m30 * m.m11 * m.m22 + m.m30 * m.m12 * m.m21
	inv01 := -m.m01 * m.m22 * m.m33 + m.m01 * m.m23 * m.m32 + m.m21 * m.m02 * m.m33 - m.m21 * m.m03 * m.m32 - m.m31 * m.m02 * m.m23 + m.m31 * m.m03 * m.m22
	inv05 :=  m.m00 * m.m22 * m.m33 - m.m00 * m.m23 * m.m32 - m.m20 * m.m02 * m.m33 + m.m20 * m.m03 * m.m32 + m.m30 * m.m02 * m.m23 - m.m30 * m.m03 * m.m22
	inv09 := -m.m00 * m.m21 * m.m33 + m.m00 * m.m23 * m.m31 + m.m20 * m.m01 * m.m33 - m.m20 * m.m03 * m.m31 - m.m30 * m.m01 * m.m23 + m.m30 * m.m03 * m.m21
	inv13 :=  m.m00 * m.m21 * m.m32 - m.m00 * m.m22 * m.m31 - m.m20 * m.m01 * m.m32 + m.m20 * m.m02 * m.m31 + m.m30 * m.m01 * m.m22 - m.m30 * m.m02 * m.m21
	inv02 :=  m.m01 * m.m12 * m.m33 - m.m01 * m.m13 * m.m32 - m.m11 * m.m02 * m.m33 + m.m11 * m.m03 * m.m32 + m.m31 * m.m02 * m.m13 - m.m31 * m.m03 * m.m12
	inv06 := -m.m00 * m.m12 * m.m33 + m.m00 * m.m13 * m.m32 + m.m10 * m.m02 * m.m33 - m.m10 * m.m03 * m.m32 - m.m30 * m.m02 * m.m13 + m.m30 * m.m03 * m.m12
	inv10 :=  m.m00 * m.m11 * m.m33 - m.m00 * m.m13 * m.m31 - m.m10 * m.m01 * m.m33 + m.m10 * m.m03 * m.m31 + m.m30 * m.m01 * m.m13 - m.m30 * m.m03 * m.m11
	inv14 := -m.m00 * m.m11 * m.m32 + m.m00 * m.m12 * m.m31 + m.m10 * m.m01 * m.m32 - m.m10 * m.m02 * m.m31 - m.m30 * m.m01 * m.m12 + m.m30 * m.m02 * m.m11
	inv03 := -m.m01 * m.m12 * m.m23 + m.m01 * m.m13 * m.m22 + m.m11 * m.m02 * m.m23 - m.m11 * m.m03 * m.m22 - m.m21 * m.m02 * m.m13 + m.m21 * m.m03 * m.m12
	inv07 :=  m.m00 * m.m12 * m.m23 - m.m00 * m.m13 * m.m22 - m.m10 * m.m02 * m.m23 + m.m10 * m.m03 * m.m22 + m.m20 * m.m02 * m.m13 - m.m20 * m.m03 * m.m12
	inv11 := -m.m00 * m.m11 * m.m23 + m.m00 * m.m13 * m.m21 + m.m10 * m.m01 * m.m23 - m.m10 * m.m03 * m.m21 - m.m20 * m.m01 * m.m13 + m.m20 * m.m03 * m.m11
	inv15 :=  m.m00 * m.m11 * m.m22 - m.m00 * m.m12 * m.m21 - m.m10 * m.m01 * m.m22 + m.m10 * m.m02 * m.m21 + m.m20 * m.m01 * m.m12 - m.m20 * m.m02 * m.m11

	det := m.m00 * inv00 + m.m01 * inv04 + m.m02 * inv08 + m.m03 * inv12

	idet := 1.0 / det

	return Matrix4x4{
		inv00 * idet, inv01 * idet, inv02 * idet, inv03 * idet,
		inv04 * idet, inv05 * idet, inv06 * idet, inv07 * idet,
		inv08 * idet, inv09 * idet, inv10 * idet, inv11 * idet,
		inv12 * idet, inv13 * idet, inv14 * idet, inv15 * idet,
	}
}