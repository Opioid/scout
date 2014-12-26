package math

type Matrix3x3 struct {
/*	m[0] m00,	m[1] m01,	m[2] m02,
	m[3] m10, 	m[4] m11,	m[5] m12,
	m[6] m20,	m[7] m21,	m[8] m22 float32 */

//	m [9]float32

	m00, m01, m02,
	m10, m11, m12,
	m20, m21, m22 float32
}

func MakeIdentityMatrix3x3() Matrix3x3 {
	return Matrix3x3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
}

func MakeMatrix3x3FromAxes(x, y, z Vector3) Matrix3x3 {
	return Matrix3x3{
		x.X, x.Y, x.Z,
		y.X, y.Y, y.Z,
		z.X, z.Y, z.Z,
	}
}

func NewMatrix3x3FromQuaternion(q Quaternion) *Matrix3x3 {
	d := q.Dot(q)

	s := 2.0 / d

	xs, ys, zs := q.X * s,  q.Y * s,  q.Z * s
	wx, wy, wz := q.W * xs, q.W * ys, q.W * zs
	xx, xy, xz := q.X * xs, q.X * ys, q.X * zs
	yy, yz, zz := q.Y * ys, q.Y * zs, q.Z * zs

	return &Matrix3x3{
		1.0 - (yy + zz), xy - wz,         xz + wy,
		xy + wz,         1.0 - (xx + zz), yz - wx,
		xz - wy,         yz + wx,         1.0 - (xx + yy),
	}

/*	T d = dot(q, q);

	T s = T(2) / d;

	T xs = q.x * s,  ys = q.y * s,  zs = q.z * s;
	T wx = q.w * xs, wy = q.w * ys, wz = q.w * zs;
	T xx = q.x * xs, xy = q.x * ys, xz = q.x * zs;
	T yy = q.y * ys, yz = q.y * zs, zz = q.z * zs;

	m00 = T(1) - (yy + zz); m01 = xy - wz;          m02 = xz + wy;
	m10 = xy + wz;          m11 = T(1) - (xx + zz); m12 = yz - wx;
	m20 = xz - wy;          m21 = yz + wx,          m22 = T(1) - (xx + yy);
	*/
}

func (m *Matrix3x3) Row(i int32) Vector3 {
//	return MakeVector3(m.m[i * 3], m.m[i * 3 + 1], m.m[i * 3 + 2])
	switch i {
	case 0:
		return MakeVector3(m.m00, m.m01, m.m02)
	case 1:
		return MakeVector3(m.m10, m.m11, m.m12)
	default:
		return MakeVector3(m.m20, m.m21, m.m22)
	}	
}

func (m *Matrix3x3) Right() Vector3 {
	return MakeVector3(m.m00, m.m01, m.m02)
}

func (m *Matrix3x3) Up() Vector3 {
	return MakeVector3(m.m10, m.m11, m.m12)
}

func (m *Matrix3x3) Direction() Vector3 {
	return MakeVector3(m.m20, m.m21, m.m22)
}

func (m *Matrix3x3) At(i, j int32) float32 {
	return m.Row(i).At(j)
}

func (m *Matrix3x3) SetFromQuaternion(q Quaternion) {
	d := q.Dot(q)

	s := 2.0 / d

	xs, ys, zs := q.X * s,  q.Y * s,  q.Z * s
	wx, wy, wz := q.W * xs, q.W * ys, q.W * zs
	xx, xy, xz := q.X * xs, q.X * ys, q.X * zs
	yy, yz, zz := q.Y * ys, q.Y * zs, q.Z * zs

	m.m00 = 1.0 - (yy + zz); m.m01 = xy - wz;         m.m02 = xz + wy
	m.m10 = xy + wz;         m.m11 = 1.0 - (xx + zz); m.m12 = yz - wx
	m.m20 = xz - wy;         m.m21 = yz + wx;         m.m22 = 1.0 - (xx + yy)
}

/*
template<typename T>
inline void setBasis(Matrix3x3_t<T> &m, const Vector3_t<T> &v)
{
	m.rows[2] = v;

	if (v.x < T(0.6) && v.x > -T(0.6)) 
		m.rows[1] = Vector3_t<T>(T(1), T(0), T(0));
	else if (v.y < T(0.6) && v.y > T(0.6)) 
		m.rows[1] = Vector3_t<T>(T(0), T(1), T(0));
	else 
		m.rows[1] = Vector3_t<T>(T(0), T(0), T(1));
	
	m.rows[0] = normalize(cross(v, m.rows[1]));
	m.rows[1] = cross(m.rows[0], m.rows[2]);
}
*/

func (m *Matrix3x3) SetBasis(v Vector3) {
	var r1 Vector3

	if v.X < 0.6 && v.X > -0.6 {
		r1 = MakeVector3(1.0, 0.0, 0.0)
	} else if v.Y < 0.6 && v.Y > -0.6 {
		r1 = MakeVector3(0.0, 1.0, 0.0)
	} else {
		r1 = MakeVector3(0.0, 0.0, 1.0)
	}

	r0 := v.Cross(r1).Normalized()
	r1 = r0.Cross(v)

	m.m00 = r0.X; m.m01 = r0.Y; m.m02 = r0.Z
	m.m10 = r1.X; m.m11 = r1.Y; m.m12 = r1.Z
	m.m20 =  v.X; m.m21 =  v.Y; m.m22 =  v.Z
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
	m.m22 = m.m20 * o.m02 + m.m21 * o.m12 + m.m22 * o.m22
}

func (m *Matrix3x3) Multiply(o *Matrix3x3) Matrix3x3 {
	return Matrix3x3{
		m.m00 * o.m00 + m.m01 * o.m10 + m.m02 * o.m20,
		m.m00 * o.m01 + m.m01 * o.m11 + m.m02 * o.m21,
		m.m00 * o.m02 + m.m01 * o.m12 + m.m02 * o.m22,

		m.m10 * o.m00 + m.m11 * o.m10 + m.m12 * o.m20,
		m.m10 * o.m01 + m.m11 * o.m11 + m.m12 * o.m21,
		m.m10 * o.m02 + m.m11 * o.m12 + m.m12 * o.m22,

		m.m20 * o.m00 + m.m21 * o.m10 + m.m22 * o.m20,
		m.m20 * o.m01 + m.m21 * o.m11 + m.m22 * o.m21,
		m.m20 * o.m02 + m.m21 * o.m12 + m.m22 * o.m22,
	}
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

func (m *Matrix3x3) TransformVector3(v Vector3) Vector3 {
	return MakeVector3(
		v.X * m.m00 + v.Y * m.m10 + v.Z * m.m20,
		v.X * m.m01 + v.Y * m.m11 + v.Z * m.m21,
		v.X * m.m02 + v.Y * m.m12 + v.Z * m.m22,
	)
}

func (m *Matrix3x3) TransposedTransformVector3(v Vector3) Vector3 {
	return MakeVector3(
		v.X * m.m00 + v.Y * m.m01 + v.Z * m.m02,
		v.X * m.m10 + v.Y * m.m11 + v.Z * m.m12,
		v.X * m.m20 + v.Y * m.m21 + v.Z * m.m22,
	)
}