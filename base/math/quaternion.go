package math

type Quaternion struct {
	X, Y, Z, W float32
}

func MakeIdentityQuaterion() Quaternion {
	return Quaternion{0.0, 0.0, 0.0, 1.0}
}

func MakeQuaternionFromMatrix3x3(m *Matrix3x3) Quaternion {
	trace := m.m[0] + m.m[4] + m.m[8]
	var temp [4]float32

	if trace > 0 {
		s := Sqrt(trace + 1.0)
		temp[3] = s * 0.5
		s = 0.5 / s

		temp[0] = (m.m[7] - m.m[5]) * s
		temp[1] = (m.m[2] - m.m[6]) * s
		temp[2] = (m.m[3] - m.m[1]) * s
	} else {
		var i int
		if m.m[0] < m.m[4] {
			if m.m[4] < m.m[8] {
				i = 2
			} else {
				i = 1
			}
		} else {
			if m.m[0] < m.m[8] {
				i = 2
			} else {
				i = 0
			}
		}

		j := (i + 1) % 3
		k := (i + 2) % 3

		s := Sqrt(m.m[i * 3 + i] - m.m[j * 3 + j] - m.m[k * 3 + k] + 1.0)
		temp[i] = s * 0.5
		s = 0.5 / s

		temp[3] = (m.m[k * 3 + j] - m.m[j * 3 + k]) * s
		temp[j] = (m.m[j * 3 + i] + m.m[i * 3 + j]) * s
		temp[k] = (m.m[k * 3 + i] + m.m[i * 3 + k]) * s
	}
	
/*
	if (trace > T(0))
	{
		T s = sqrt(trace + T(1));
		temp[3] = s * T(0.5);
		s = T(0.5) / s;

		temp[0] = (m.m21 - m.m12) * s;
		temp[1] = (m.m02 - m.m20) * s;
		temp[2] = (m.m10 - m.m01) * s;
	}
	else
	{
		int i = m.m00 < m.m11 ? (m.m11 < m.m22 ? 2 : 1) :(m.m00 < m.m22 ? 2 : 0);
		int j = (i + 1) % 3;
		int k = (i + 2) % 3;

		T s = sqrt(m.m[i * 3 + i] - m.m[j * 3 + j] - m.m[k * 3 + k] + T(1));
		temp[i] = s * T(0.5);
		s = T(0.5) / s;

		temp[3] = (m.m[k * 3 + j] - m.m[j * 3 + k]) * s;
		temp[j] = (m.m[j * 3 + i] + m.m[i * 3 + j]) * s;
		temp[k] = (m.m[k * 3 + i] + m.m[i * 3 + k]) * s;
	}

	x = temp[0];
	y = temp[1];
	z = temp[2];
	w = temp[3];
	*/

	return Quaternion{temp[0], temp[1], temp[2], temp[3]}
}

func (a Quaternion) Dot(b Quaternion) float32 {
	return a.X * b.X + a.Y * b.Y + a.Z * b.Z + a.W * b.W
}