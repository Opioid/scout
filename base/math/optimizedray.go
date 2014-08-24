package math

type OptimizedRay struct {
	Ray
	ReciprocalDirection Vector3
}

func MakeOptimizedRay(origin, direction Vector3, mint, maxt float32) OptimizedRay {
	r := OptimizedRay{}
	r.Origin = origin
	r.SetDirection(direction)
	r.MinT = mint
	r.MaxT = maxt
	return r
}

func (r *OptimizedRay) SetDirection(direction Vector3) {
	r.Direction = direction
	r.ReciprocalDirection = Vector3{1.0 / direction.X, 1.0 / direction.Y, 1.0 / direction.Z}
}