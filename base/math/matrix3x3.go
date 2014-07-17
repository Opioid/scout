package math

type Matrix3x3 struct {
	m00, m01, m02,
	m10, m11, m12,
	m20, m21, m22 float32 
}

func (m *Matrix3x3) MuliplyAssign(o *Matrix3x3) {
	m.m00 = m.m00 * o.m00 + m.m01 * o.m10 + m.m02 * o.m20
	m.m01 = m.m00 * o.m01 + m.m01 * o.m11 + m.m02 * o.m21
	m.m02 = m.m00 * o.m02 + m.m01 * o.m12 + m.m02 * o.m22

	m.m10 = m.m10 * o.m00 + m.m11 * o.m10 + m.m12 * o.m20
	m.m11 = m.m10 * o.m01 + m.m11 * o.m11 + m.m12 * o.m21
	m.m12 = m.m10 * o.m02 + m.m11 * o.m12 + m.m12 * o.m22

	m.m20 = m.m20 * o.m00 + m.m21 * o.m10 + m.m22 * o.m20
	m.m21 = m.m20 * o.m01 + m.m21 * o.m11 + m.m22 * o.m21
	m.m21 = m.m20 * o.m02 + m.m21 * o.m12 + m.m22 * o.m22
}

func (m *Matrix3x3) Multiply(o *Matrix3x3) *Matrix3x3 {
	return &Matrix3x3{ 
		m.m00 * o.m00 + m.m01 * o.m10 + m.m02 * o.m20,
		m.m00 * o.m01 + m.m01 * o.m11 + m.m02 * o.m21,
		m.m00 * o.m02 + m.m01 * o.m12 + m.m02 * o.m22,

		m.m10 * o.m00 + m.m11 * o.m10 + m.m12 * o.m20,
		m.m10 * o.m01 + m.m11 * o.m11 + m.m12 * o.m21,
		m.m10 * o.m02 + m.m11 * o.m12 + m.m12 * o.m22,

		m.m20 * o.m00 + m.m21 * o.m10 + m.m22 * o.m20,
		m.m20 * o.m01 + m.m21 * o.m11 + m.m22 * o.m21,
		m.m20 * o.m02 + m.m21 * o.m12 + m.m22 * o.m22}
}

func (m *Matrix3x3) SetRotationX(a float32) {
	c, s := Cos(a), Sin(a)

	m.m00 = 1.0; m.m01 = 0.0; m.m02 = 0.0
	m.m10 = 0.0; m.m11 = c;   m.m12 = -s
	m.m20 = 0.0; m.m21 = s;   m.m22 =  c
}

func (m *Matrix3x3) SetRotationY(a float32) {
	c, s := Cos(a), Sin(a)

	m.m00 =  c;   m.m01 = 0.0; m.m02 = s;
	m.m10 =  0.0; m.m11 = 1.0; m.m12 = 0.0;
	m.m20 = -s;   m.m21 = 0.0; m.m22 = c;
}

func (m *Matrix3x3) SetRotationZ(a float32) {
	c, s := Cos(a), Sin(a)

	m.m00 = c;   m.m01 = -s;   m.m02 = 0.0;
	m.m10 = s;   m.m11 =  c;   m.m12 = 0.0;
	m.m20 = 0.0; m.m21 =  0.0; m.m22 = 1.0;
}