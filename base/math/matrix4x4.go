package math

type Matrix4x4 struct {
	/*
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float32 */
	/*
	 m0,  m1,  m2,  m3,
	 m4,  m5,  m6,  m7,
	 m8,  m9, m10, m11,
	m12, m13, m14, m15 */

	m [16]float32
}

func (m *Matrix4x4) Row(i int) Vector4 {
	return Vector4{m.m[i * 4], m.m[i * 4 + 1], m.m[i * 4 + 2], m.m[i * 4 + 3]}
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

func (m *Matrix4x4) TransformPoint(v Vector3) Vector3 {
	/*
	return Vector3<T>(v.x * m.m00 + v.y * m.m10 + v.z * m.m20 + m.m30,
					  v.x * m.m01 + v.y * m.m11 + v.z * m.m21 + m.m31,
					  v.x * m.m02 + v.y * m.m12 + v.z * m.m22 + m.m32);
	*/

	return Vector3{
		v.X * m.m[0] + v.Y * m.m[4] + v.Z * m.m[8]  + m.m[12],
		v.X * m.m[1] + v.Y * m.m[5] + v.Z * m.m[9]  + m.m[13],
		v.X * m.m[2] + v.Y * m.m[6] + v.Z * m.m[10] + m.m[14],
	}
}