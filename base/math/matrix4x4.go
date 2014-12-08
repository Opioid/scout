package math

type Matrix4x4 struct {
	
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float32
	/*
	 m[ 0],  m[ 1], m[ 2], m[ 3],
	 m10,  m11, m12, m13,
	 m20,  m21, m22, m23,
	 m30,  m31, m32, m33 */

//	m [16]float32
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
		v.X * m.m00 + v.Y * m.m10 + v.Z * m.m20  + m.m30,
		v.X * m.m01 + v.Y * m.m11 + v.Z * m.m21  + m.m31,
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

/*
template<typename T>
inline T det(const Matrix4x4<T>& m)
{
	return m.m00 * m.m11 * m.m22 * m.m33 + m.m00 * m.m12 * m.m23 * m.m31 + m.m00 * m.m13 * m.m21 * m.m32
		 + m.m01 * m.m10 * m.m23 * m.m32 + m.m01 * m.m12 * m.m20 * m.m33 + m.m01 * m.m13 * m.m22 * m.m30
		 + m.m02 * m.m10 * m.m21 * m.m33 + m.m02 * m.m11 * m.m23 * m.m30 + m.m02 * m.m13 * m.m20 * m.m31
		 + m.m03 * m.m10 * m.m22 * m.m31 + m.m03 * m.m11 * m.m21 * m.m32 + m.m03 * m.m12 * m.m21 * m.m30

		 - m.m00 * m.m11 * m.m23 * m.m32 - m.m00 * m.m12 * m.m21 * m.m33 - m.m00 * m.m13 * m.m22 * m.m31
		 - m.m01 * m.m10 * m.m22 * m.m33 - m.m01 * m.m12 * m.m23 * m.m30 - m.m01 * m.m13 * m.m20 * m.m32
		 - m.m02 * m.m10 * m.m23 * m.m31 - m.m02 * m.m11 * m.m20 * m.m33 - m.m02 * m.m13 * m.m21 * m.m30
		 - m.m03 * m.m10 * m.m21 * m.m32 - m.m03 * m.m11 * m.m22 * m.m30 - m.m03 * m.m12 * m.m20 * m.m31;
}
*/

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

/*
return Matrix4x4<T>(m.m11*m.m22*m.m33 + m.m12*m.m23*m.m31 + m.m13*m.m21*m.m32 - m.m11*m.m23*m.m32 - m.m12*m.m21*m.m33 - m.m13*m.m22*m.m31,
						m.m01*m.m23*m.m32 + m.m02*m.m21*m.m33 + m.m03*m.m22*m.m31 - m.m01*m.m22*m.m33 - m.m02*m.m23*m.m31 - m.m03*m.m21*m.m32,
						m.m01*m.m12*m.m33 + m.m02*m.m13*m.m31 + m.m03*m.m11*m.m32 - m.m01*m.m13*m.m32 - m.m02*m.m11*m.m33 - m.m03*m.m12*m.m31,
						m.m01*m.m13*m.m22 + m.m02*m.m11*m.m23 + m.m03*m.m12*m.m21 - m.m01*m.m12*m.m23 - m.m02*m.m13*m.m21 - m.m03*m.m11*m.m22,

						m.m10*m.m23*m.m32 + m.m12*m.m20*m.m33 + m.m13*m.m22*m.m30 - m.m10*m.m22*m.m33 - m.m12*m.m23*m.m30 - m.m13*m.m20*m.m32,
						m.m00*m.m22*m.m33 + m.m02*m.m23*m.m30 + m.m03*m.m20*m.m32 - m.m00*m.m23*m.m32 - m.m02*m.m20*m.m33 - m.m03*m.m22*m.m30,
						m.m00*m.m13*m.m32 + m.m02*m.m10*m.m33 + m.m03*m.m12*m.m30 - m.m00*m.m12*m.m33 - m.m02*m.m13*m.m30 - m.m03*m.m10*m.m32,
						m.m00*m.m12*m.m23 + m.m02*m.m13*m.m20 + m.m03*m.m10*m.m22 - m.m00*m.m13*m.m22 - m.m02*m.m10*m.m23 - m.m03*m.m12*m.m20,

						m.m10*m.m21*m.m33 + m.m11*m.m23*m.m30 + m.m13*m.m20*m.m31 - m.m10*m.m23*m.m31 - m.m11*m.m20*m.m33 - m.m13*m.m21*m.m30,
						m.m00*m.m23*m.m31 + m.m01*m.m20*m.m33 + m.m03*m.m21*m.m30 - m.m00*m.m21*m.m33 - m.m01*m.m23*m.m30 - m.m03*m.m20*m.m31,
						m.m00*m.m11*m.m33 + m.m01*m.m13*m.m30 + m.m03*m.m10*m.m31 - m.m00*m.m13*m.m31 - m.m01*m.m10*m.m33 - m.m03*m.m11*m.m30,
						m.m00*m.m13*m.m21 + m.m01*m.m10*m.m23 + m.m03*m.m11*m.m20 - m.m00*m.m11*m.m23 - m.m01*m.m13*m.m20 - m.m03*m.m10*m.m21,

						m.m10*m.m22*m.m31 + m.m11*m.m20*m.m32 + m.m12*m.m21*m.m30 - m.m10*m.m21*m.m32 - m.m11*m.m22*m.m30 - m.m12*m.m20*m.m31,
						m.m00*m.m21*m.m32 + m.m01*m.m22*m.m30 + m.m02*m.m20*m.m31 - m.m00*m.m22*m.m31 - m.m01*m.m20*m.m32 - m.m02*m.m21*m.m30,
						m.m00*m.m12*m.m31 + m.m01*m.m10*m.m32 + m.m02*m.m11*m.m30 - m.m00*m.m11*m.m32 - m.m01*m.m12*m.m30 - m.m02*m.m10*m.m31,
						m.m00*m.m11*m.m22 + m.m01*m.m12*m.m20 + m.m02*m.m10*m.m21 - m.m00*m.m12*m.m21 - m.m01*m.m10*m.m22 - m.m02*m.m11*m.m20) / det(m);
*/

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
