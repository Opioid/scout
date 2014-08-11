package math

type Matrix4x4 struct {
	/*
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float32 */
	/*
	 m[ 0],  m[ 1], m[ 2], m[ 3],
	 m[ 4],  m[ 5], m[ 6], m[ 7],
	 m[ 8],  m[ 9], m[10], m[11],
	 m[12],  m[13], m[14], m[15] */

	m [16]float32
}

func (m *Matrix4x4) Row(i int) Vector4 {
	return Vector4{m.m[i * 4], m.m[i * 4 + 1], m.m[i * 4 + 2], m.m[i * 4 + 3]}
}

func (m *Matrix4x4) Div(s float32) Matrix4x4 {
	is := 1.0 / s
	return Matrix4x4{[16]float32{
		m.m[ 0] * is, m.m[ 1] * is, m.m[ 2] * is, m.m[ 3] * is,
		m.m[ 4] * is, m.m[ 5] * is, m.m[ 6] * is, m.m[ 7] * is,
		m.m[ 8] * is, m.m[ 9] * is, m.m[10] * is, m.m[11] * is,
		m.m[12] * is, m.m[13] * is, m.m[14] * is, m.m[15] * is,
	}}
}

func (m *Matrix4x4) SetIdentity() {
	m.m[0]  = 1.0; m.m[1]  = 0.0; m.m[2]  = 0.0; m.m[3]  = 0.0
	m.m[4]  = 0.0; m.m[5]  = 1.0; m.m[6]  = 0.0; m.m[7]  = 0.0
	m.m[8]  = 0.0; m.m[9]  = 0.0; m.m[10] = 1.0; m.m[11] = 0.0
	m.m[12] = 0.0; m.m[13] = 0.0; m.m[14] = 0.0; m.m[15] = 1.0
}

func (m *Matrix4x4) SetBasis(b *Matrix3x3) {
	m.m[0] = b.m[0]; m.m[1] = b.m[1]; m.m[2]  = b.m[2]
	m.m[4] = b.m[3]; m.m[5] = b.m[4]; m.m[6]  = b.m[5]
	m.m[8] = b.m[6]; m.m[9] = b.m[7]; m.m[10] = b.m[8]
}

func (m *Matrix4x4) SetOrigin(v Vector3) {
	m.m[12] = v.X
	m.m[13] = v.Y
	m.m[14] = v.Z
}

func (m *Matrix4x4) Scale(v Vector3) {
	m.m[0] *= v.X; m.m[1] *= v.X; m.m[ 2] *= v.X
	m.m[4] *= v.Y; m.m[5] *= v.Y; m.m[ 6] *= v.Y
	m.m[8] *= v.Z; m.m[9] *= v.Z; m.m[10] *= v.Z
}

func (m *Matrix4x4) TransformPoint(v Vector3) Vector3 {
	return Vector3{
		v.X * m.m[0] + v.Y * m.m[4] + v.Z * m.m[8]  + m.m[12],
		v.X * m.m[1] + v.Y * m.m[5] + v.Z * m.m[9]  + m.m[13],
		v.X * m.m[2] + v.Y * m.m[6] + v.Z * m.m[10] + m.m[14],
	}
}

func (m *Matrix4x4) TransformVector(v Vector3) Vector3 {
	return Vector3{
		v.X * m.m[0] + v.Y * m.m[4] + v.Z * m.m[8],
		v.X * m.m[1] + v.Y * m.m[5] + v.Z * m.m[9],
		v.X * m.m[2] + v.Y * m.m[6] + v.Z * m.m[10],
	}
}

/*
template<typename T>
inline T det(const Matrix4x4<T>& m)
{
	return m.m[0] * m.m[5] * m.m[10] * m.m[15] + m.m[0] * m.m[6] * m.m[11] * m.m[13] + m.m[0] * m.m[7] * m.m[9] * m.m[14]
		 + m.m[1] * m.m[4] * m.m[11] * m.m[14] + m.m[1] * m.m[6] * m.m[8] * m.m[15] + m.m[1] * m.m[7] * m.m[10] * m.m[12]
		 + m.m[2] * m.m[4] * m.m[9] * m.m[15] + m.m[2] * m.m[5] * m.m[11] * m.m[12] + m.m[2] * m.m[7] * m.m[8] * m.m[13]
		 + m.m[3] * m.m[4] * m.m[10] * m.m[13] + m.m[3] * m.m[5] * m.m[9] * m.m[14] + m.m[3] * m.m[6] * m.m[9] * m.m[12]

		 - m.m[0] * m.m[5] * m.m[11] * m.m[14] - m.m[0] * m.m[6] * m.m[9] * m.m[15] - m.m[0] * m.m[7] * m.m[10] * m.m[13]
		 - m.m[1] * m.m[4] * m.m[10] * m.m[15] - m.m[1] * m.m[6] * m.m[11] * m.m[12] - m.m[1] * m.m[7] * m.m[8] * m.m[14]
		 - m.m[2] * m.m[4] * m.m[11] * m.m[13] - m.m[2] * m.m[5] * m.m[8] * m.m[15] - m.m[2] * m.m[7] * m.m[9] * m.m[12]
		 - m.m[3] * m.m[4] * m.m[9] * m.m[14] - m.m[3] * m.m[5] * m.m[10] * m.m[12] - m.m[3] * m.m[6] * m.m[8] * m.m[13];
}
*/

func (m *Matrix4x4) Det() float32 {
	return m.m[0] * m.m[5] * m.m[10] * m.m[15] + m.m[0] * m.m[6] * m.m[11] * m.m[13] + m.m[0] * m.m[7] * m.m[ 9] * m.m[14] +
		   m.m[1] * m.m[4] * m.m[11] * m.m[14] + m.m[1] * m.m[6] * m.m[ 8] * m.m[15] + m.m[1] * m.m[7] * m.m[10] * m.m[12] +
		   m.m[2] * m.m[4] * m.m[ 9] * m.m[15] + m.m[2] * m.m[5] * m.m[11] * m.m[12] + m.m[2] * m.m[7] * m.m[ 8] * m.m[13] +
		   m.m[3] * m.m[4] * m.m[10] * m.m[13] + m.m[3] * m.m[5] * m.m[ 9] * m.m[14] + m.m[3] * m.m[6] * m.m[ 9] * m.m[12] -

		   m.m[0] * m.m[5] * m.m[11] * m.m[14] - m.m[0] * m.m[6] * m.m[ 9] * m.m[15] - m.m[0] * m.m[7] * m.m[10] * m.m[13] -
		   m.m[1] * m.m[4] * m.m[10] * m.m[15] - m.m[1] * m.m[6] * m.m[11] * m.m[12] - m.m[1] * m.m[7] * m.m[ 8] * m.m[14] -
		   m.m[2] * m.m[4] * m.m[11] * m.m[13] - m.m[2] * m.m[5] * m.m[ 8] * m.m[15] - m.m[2] * m.m[7] * m.m[ 9] * m.m[12] -
		   m.m[3] * m.m[4] * m.m[ 9] * m.m[14] - m.m[3] * m.m[5] * m.m[10] * m.m[12] - m.m[3] * m.m[6] * m.m[ 8] * m.m[13]
}

/*
return Matrix4x4<T>(m.m[5]*m.m[10]*m.m[15] + m.m[6]*m.m[11]*m.m[13] + m.m[7]*m.m[9]*m.m[14] - m.m[5]*m.m[11]*m.m[14] - m.m[6]*m.m[9]*m.m[15] - m.m[7]*m.m[10]*m.m[13],
						m.m[1]*m.m[11]*m.m[14] + m.m[2]*m.m[9]*m.m[15] + m.m[3]*m.m[10]*m.m[13] - m.m[1]*m.m[10]*m.m[15] - m.m[2]*m.m[11]*m.m[13] - m.m[3]*m.m[9]*m.m[14],
						m.m[1]*m.m[6]*m.m[15] + m.m[2]*m.m[7]*m.m[13] + m.m[3]*m.m[5]*m.m[14] - m.m[1]*m.m[7]*m.m[14] - m.m[2]*m.m[5]*m.m[15] - m.m[3]*m.m[6]*m.m[13],
						m.m[1]*m.m[7]*m.m[10] + m.m[2]*m.m[5]*m.m[11] + m.m[3]*m.m[6]*m.m[9] - m.m[1]*m.m[6]*m.m[11] - m.m[2]*m.m[7]*m.m[9] - m.m[3]*m.m[5]*m.m[10],

						m.m[4]*m.m[11]*m.m[14] + m.m[6]*m.m[8]*m.m[15] + m.m[7]*m.m[10]*m.m[12] - m.m[4]*m.m[10]*m.m[15] - m.m[6]*m.m[11]*m.m[12] - m.m[7]*m.m[8]*m.m[14],
						m.m[0]*m.m[10]*m.m[15] + m.m[2]*m.m[11]*m.m[12] + m.m[3]*m.m[8]*m.m[14] - m.m[0]*m.m[11]*m.m[14] - m.m[2]*m.m[8]*m.m[15] - m.m[3]*m.m[10]*m.m[12],
						m.m[0]*m.m[7]*m.m[14] + m.m[2]*m.m[4]*m.m[15] + m.m[3]*m.m[6]*m.m[12] - m.m[0]*m.m[6]*m.m[15] - m.m[2]*m.m[7]*m.m[12] - m.m[3]*m.m[4]*m.m[14],
						m.m[0]*m.m[6]*m.m[11] + m.m[2]*m.m[7]*m.m[8] + m.m[3]*m.m[4]*m.m[10] - m.m[0]*m.m[7]*m.m[10] - m.m[2]*m.m[4]*m.m[11] - m.m[3]*m.m[6]*m.m[8],

						m.m[4]*m.m[9]*m.m[15] + m.m[5]*m.m[11]*m.m[12] + m.m[7]*m.m[8]*m.m[13] - m.m[4]*m.m[11]*m.m[13] - m.m[5]*m.m[8]*m.m[15] - m.m[7]*m.m[9]*m.m[12],
						m.m[0]*m.m[11]*m.m[13] + m.m[1]*m.m[8]*m.m[15] + m.m[3]*m.m[9]*m.m[12] - m.m[0]*m.m[9]*m.m[15] - m.m[1]*m.m[11]*m.m[12] - m.m[3]*m.m[8]*m.m[13],
						m.m[0]*m.m[5]*m.m[15] + m.m[1]*m.m[7]*m.m[12] + m.m[3]*m.m[4]*m.m[13] - m.m[0]*m.m[7]*m.m[13] - m.m[1]*m.m[4]*m.m[15] - m.m[3]*m.m[5]*m.m[12],
						m.m[0]*m.m[7]*m.m[9] + m.m[1]*m.m[4]*m.m[11] + m.m[3]*m.m[5]*m.m[8] - m.m[0]*m.m[5]*m.m[11] - m.m[1]*m.m[7]*m.m[8] - m.m[3]*m.m[4]*m.m[9],

						m.m[4]*m.m[10]*m.m[13] + m.m[5]*m.m[8]*m.m[14] + m.m[6]*m.m[9]*m.m[12] - m.m[4]*m.m[9]*m.m[14] - m.m[5]*m.m[10]*m.m[12] - m.m[6]*m.m[8]*m.m[13],
						m.m[0]*m.m[9]*m.m[14] + m.m[1]*m.m[10]*m.m[12] + m.m[2]*m.m[8]*m.m[13] - m.m[0]*m.m[10]*m.m[13] - m.m[1]*m.m[8]*m.m[14] - m.m[2]*m.m[9]*m.m[12],
						m.m[0]*m.m[6]*m.m[13] + m.m[1]*m.m[4]*m.m[14] + m.m[2]*m.m[5]*m.m[12] - m.m[0]*m.m[5]*m.m[14] - m.m[1]*m.m[6]*m.m[12] - m.m[2]*m.m[4]*m.m[13],
						m.m[0]*m.m[5]*m.m[10] + m.m[1]*m.m[6]*m.m[8] + m.m[2]*m.m[4]*m.m[9] - m.m[0]*m.m[6]*m.m[9] - m.m[1]*m.m[4]*m.m[10] - m.m[2]*m.m[5]*m.m[8]) / det(m);
*/

func (m *Matrix4x4) Inverted() Matrix4x4 {
	idet := 1.0 / m.Det()
	return Matrix4x4{[16]float32{
		(m.m[5]*m.m[10]*m.m[15] + m.m[6]*m.m[11]*m.m[13] + m.m[7]*m.m[ 9]*m.m[14] - m.m[5]*m.m[11]*m.m[14] - m.m[6]*m.m[ 9]*m.m[15] - m.m[7]*m.m[10]*m.m[13]) * idet,
		(m.m[1]*m.m[11]*m.m[14] + m.m[2]*m.m[ 9]*m.m[15] + m.m[3]*m.m[10]*m.m[13] - m.m[1]*m.m[10]*m.m[15] - m.m[2]*m.m[11]*m.m[13] - m.m[3]*m.m[ 9]*m.m[14]) * idet,
		(m.m[1]*m.m[ 6]*m.m[15] + m.m[2]*m.m[ 7]*m.m[13] + m.m[3]*m.m[ 5]*m.m[14] - m.m[1]*m.m[ 7]*m.m[14] - m.m[2]*m.m[ 5]*m.m[15] - m.m[3]*m.m[ 6]*m.m[13]) * idet,
		(m.m[1]*m.m[ 7]*m.m[10] + m.m[2]*m.m[ 5]*m.m[11] + m.m[3]*m.m[ 6]*m.m[ 9] - m.m[1]*m.m[ 6]*m.m[11] - m.m[2]*m.m[ 7]*m.m[ 9] - m.m[3]*m.m[ 5]*m.m[10]) * idet,

		(m.m[4]*m.m[11]*m.m[14] + m.m[6]*m.m[ 8]*m.m[15] + m.m[7]*m.m[10]*m.m[12] - m.m[4]*m.m[10]*m.m[15] - m.m[6]*m.m[11]*m.m[12] - m.m[7]*m.m[ 8]*m.m[14]) * idet,
		(m.m[0]*m.m[10]*m.m[15] + m.m[2]*m.m[11]*m.m[12] + m.m[3]*m.m[ 8]*m.m[14] - m.m[0]*m.m[11]*m.m[14] - m.m[2]*m.m[ 8]*m.m[15] - m.m[3]*m.m[10]*m.m[12]) * idet,
		(m.m[0]*m.m[ 7]*m.m[14] + m.m[2]*m.m[ 4]*m.m[15] + m.m[3]*m.m[ 6]*m.m[12] - m.m[0]*m.m[ 6]*m.m[15] - m.m[2]*m.m[ 7]*m.m[12] - m.m[3]*m.m[ 4]*m.m[14]) * idet,
		(m.m[0]*m.m[ 6]*m.m[11] + m.m[2]*m.m[ 7]*m.m[ 8] + m.m[3]*m.m[ 4]*m.m[10] - m.m[0]*m.m[ 7]*m.m[10] - m.m[2]*m.m[ 4]*m.m[11] - m.m[3]*m.m[ 6]*m.m[ 8]) * idet,

		(m.m[4]*m.m[ 9]*m.m[15] + m.m[5]*m.m[11]*m.m[12] + m.m[7]*m.m[ 8]*m.m[13] - m.m[4]*m.m[11]*m.m[13] - m.m[5]*m.m[ 8]*m.m[15] - m.m[7]*m.m[ 9]*m.m[12]) * idet,
		(m.m[0]*m.m[11]*m.m[13] + m.m[1]*m.m[ 8]*m.m[15] + m.m[3]*m.m[ 9]*m.m[12] - m.m[0]*m.m[ 9]*m.m[15] - m.m[1]*m.m[11]*m.m[12] - m.m[3]*m.m[ 8]*m.m[13]) * idet,
		(m.m[0]*m.m[ 5]*m.m[15] + m.m[1]*m.m[ 7]*m.m[12] + m.m[3]*m.m[ 4]*m.m[13] - m.m[0]*m.m[ 7]*m.m[13] - m.m[1]*m.m[ 4]*m.m[15] - m.m[3]*m.m[ 5]*m.m[12]) * idet,
		(m.m[0]*m.m[ 7]*m.m[ 9] + m.m[1]*m.m[ 4]*m.m[11] + m.m[3]*m.m[ 5]*m.m[ 8] - m.m[0]*m.m[ 5]*m.m[11] - m.m[1]*m.m[ 7]*m.m[ 8] - m.m[3]*m.m[ 4]*m.m[ 9]) * idet,

		(m.m[4]*m.m[10]*m.m[13] + m.m[5]*m.m[ 8]*m.m[14] + m.m[6]*m.m[ 9]*m.m[12] - m.m[4]*m.m[ 9]*m.m[14] - m.m[5]*m.m[10]*m.m[12] - m.m[6]*m.m[ 8]*m.m[13]) * idet,
		(m.m[0]*m.m[ 9]*m.m[14] + m.m[1]*m.m[10]*m.m[12] + m.m[2]*m.m[ 8]*m.m[13] - m.m[0]*m.m[10]*m.m[13] - m.m[1]*m.m[ 8]*m.m[14] - m.m[2]*m.m[ 9]*m.m[12]) * idet,
		(m.m[0]*m.m[ 6]*m.m[13] + m.m[1]*m.m[ 4]*m.m[14] + m.m[2]*m.m[ 5]*m.m[12] - m.m[0]*m.m[ 5]*m.m[14] - m.m[1]*m.m[ 6]*m.m[12] - m.m[2]*m.m[ 4]*m.m[13]) * idet,
		(m.m[0]*m.m[ 5]*m.m[10] + m.m[1]*m.m[ 6]*m.m[ 8] + m.m[2]*m.m[ 4]*m.m[ 9] - m.m[0]*m.m[ 6]*m.m[ 9] - m.m[1]*m.m[ 4]*m.m[10] - m.m[2]*m.m[ 5]*m.m[ 8]) * idet,
	}}
}
