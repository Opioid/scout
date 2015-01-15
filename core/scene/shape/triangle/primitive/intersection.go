package primitive

import (
	
)

type Coordinates struct {
	T, U, V float32
}

type Intersection struct {
	Coordinates
//	Triangle *Triangle
	Index uint32
}