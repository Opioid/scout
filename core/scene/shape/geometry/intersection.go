package geometry

type Intersection struct {
	/*Dg*/ Differential
	Epsilon float32	
	MaterialId uint32	// Material id (relative to the shape)
}