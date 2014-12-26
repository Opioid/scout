package math

type OptimizedRay struct {
	Ray
	Time float32	
	Depth uint32	
	ReciprocalDirection Vector3
	DirIsNeg [3]int
}

func MakeOptimizedRay(origin, direction Vector3, mint, maxt float32, time float32, depth uint32) OptimizedRay {
	r := OptimizedRay{}
	r.Origin = origin
	r.SetDirection(direction)
	r.MinT = mint
	r.MaxT = maxt
	r.Time = time
	r.Depth = depth
	return r
}

func (r *OptimizedRay) SetDirection(direction Vector3) {
	r.Direction = direction
	r.ReciprocalDirection = MakeVector3(1.0 / direction.X, 1.0 / direction.Y, 1.0 / direction.Z)

	isNeg := func (x float32) int {
		if x < 0.0 {
			return 1
		} else {
			return 0
		}
	}

	r.DirIsNeg[0] = isNeg(direction.X)
	r.DirIsNeg[1] = isNeg(direction.Y)
	r.DirIsNeg[2] = isNeg(direction.Z)
}