package math

type Matrix3x3 struct {
/*	m[0] m00,	m[1] m01,	m[2] m02,
	m[3] m10, 	m[4] m11,	m[5] m12,
	m[6] m20,	m[7] m21,	m[8] m22 float32 */

	m [9]float32
}

func MakeIdentityMatrix3x3() Matrix3x3 {
	return Matrix3x3{[9]float32{
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}}
}

func NewMatrix3x3FromQuaternion(q Quaternion) *Matrix3x3 {
	d := q.Dot(q)

	s := 2.0 / d

	xs, ys, zs := q.X * s,  q.Y * s,  q.Z * s
	wx, wy, wz := q.W * xs, q.W * ys, q.W * zs
	xx, xy, xz := q.X * xs, q.X * ys, q.X * zs
	yy, yz, zz := q.Y * ys, q.Y * zs, q.Z * zs

	return &Matrix3x3{[9]float32{
		1.0 - (yy + zz), xy - wz,         xz + wy,
		xy + wz,         1.0 - (xx + zz), yz - wx,
		xz - wy,         yz + wx,         1.0 - (xx + yy),
	}}

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

func (m *Matrix3x3) SetFromQuaternion(q Quaternion) {
	d := q.Dot(q)

	s := 2.0 / d

	xs, ys, zs := q.X * s,  q.Y * s,  q.Z * s
	wx, wy, wz := q.W * xs, q.W * ys, q.W * zs
	xx, xy, xz := q.X * xs, q.X * ys, q.X * zs
	yy, yz, zz := q.Y * ys, q.Y * zs, q.Z * zs

	m.m[0] = 1.0 - (yy + zz); m.m[1] = xy - wz;         m.m[2] = xz + wy
	m.m[3] = xy + wz;         m.m[4] = 1.0 - (xx + zz); m.m[5] = yz - wx
	m.m[6] = xz - wy;         m.m[7] = yz + wx;         m.m[8] = 1.0 - (xx + yy)
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

	m.m[0] = r0.X; m.m[1] = r0.Y; m.m[2] = r0.Z
	m.m[3] = r1.X; m.m[4] = r1.Y; m.m[5] = r1.Z
	m.m[6] =  v.X; m.m[7] =  v.Y; m.m[8] =  v.Z
}


func (m *Matrix3x3) Row(i int) Vector3 {
	return MakeVector3(m.m[i * 3], m.m[i * 3 + 1], m.m[i * 3 + 2])
}

func (m *Matrix3x3) MuliplyAssign(o *Matrix3x3) {
	m.m[0] = m.m[0] * o.m[0] + m.m[1] * o.m[3] + m.m[2] * o.m[6]
	m.m[1] = m.m[0] * o.m[1] + m.m[1] * o.m[4] + m.m[2] * o.m[7]
	m.m[2] = m.m[0] * o.m[2] + m.m[1] * o.m[5] + m.m[2] * o.m[8]

	m.m[3] = m.m[3] * o.m[0] + m.m[4] * o.m[3] + m.m[5] * o.m[6]
	m.m[4] = m.m[3] * o.m[1] + m.m[4] * o.m[4] + m.m[5] * o.m[7]
	m.m[5] = m.m[3] * o.m[2] + m.m[4] * o.m[5] + m.m[5] * o.m[8]

	m.m[6] = m.m[6] * o.m[0] + m.m[7] * o.m[3] + m.m[8] * o.m[6]
	m.m[7] = m.m[6] * o.m[1] + m.m[7] * o.m[4] + m.m[8] * o.m[7]
	m.m[8] = m.m[6] * o.m[2] + m.m[7] * o.m[5] + m.m[8] * o.m[8]
}

func (m *Matrix3x3) Multiply(o *Matrix3x3) Matrix3x3 {
	return Matrix3x3{[9]float32{
		m.m[0] * o.m[0] + m.m[1] * o.m[3] + m.m[2] * o.m[6],
		m.m[0] * o.m[1] + m.m[1] * o.m[4] + m.m[2] * o.m[7],
		m.m[0] * o.m[2] + m.m[1] * o.m[5] + m.m[2] * o.m[8],

		m.m[3] * o.m[0] + m.m[4] * o.m[3] + m.m[5] * o.m[6],
		m.m[3] * o.m[1] + m.m[4] * o.m[4] + m.m[5] * o.m[7],
		m.m[3] * o.m[2] + m.m[4] * o.m[5] + m.m[5] * o.m[8],

		m.m[6] * o.m[0] + m.m[7] * o.m[3] + m.m[8] * o.m[6],
		m.m[6] * o.m[1] + m.m[7] * o.m[4] + m.m[8] * o.m[7],
		m.m[6] * o.m[2] + m.m[7] * o.m[5] + m.m[8] * o.m[8],
	}}
}

func (m *Matrix3x3) SetRotationX(a float32) {
	c, s := Cos(a), Sin(a)

	m.m[0] = 1.0; m.m[1] = 0.0; m.m[2] = 0.0
	m.m[3] = 0.0; m.m[4] = c;   m.m[5] = -s
	m.m[6] = 0.0; m.m[7] = s;   m.m[8] =  c
}

func (m *Matrix3x3) SetRotationY(a float32) {
	c, s := Cos(a), Sin(a)

	m.m[0] =  c;   m.m[1] = 0.0; m.m[2] = s;
	m.m[3] =  0.0; m.m[4] = 1.0; m.m[5] = 0.0;
	m.m[6] = -s;   m.m[7] = 0.0; m.m[8] = c;
}

func (m *Matrix3x3) SetRotationZ(a float32) {
	c, s := Cos(a), Sin(a)

	m.m[0] = c;   m.m[1] = -s;   m.m[2] = 0.0;
	m.m[3] = s;   m.m[4] =  c;   m.m[5] = 0.0;
	m.m[6] = 0.0; m.m[7] =  0.0; m.m[8] = 1.0;
}

func (m *Matrix3x3) TransformVector3(v Vector3) Vector3 {
	/*
	return Vector3<T>(v.x * m.m00 + v.y * m.m10 + v.z * m.m20,
					  v.x * m.m01 + v.y * m.m11 + v.z * m.m21,
					  v.x * m.m02 + v.y * m.m12 + v.z * m.m22);
	*/

	return MakeVector3(
		v.X * m.m[0] + v.Y * m.m[3] + v.Z * m.m[6],
		v.X * m.m[1] + v.Y * m.m[4] + v.Z * m.m[7],
		v.X * m.m[2] + v.Y * m.m[5] + v.Z * m.m[8],
	)
}