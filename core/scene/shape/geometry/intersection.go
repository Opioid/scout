package geometry

type Intersection struct {
	/*Dg*/ Differential
	Epsilon float32	
	MaterialIndex uint32	// Material id (relative to the shape)
}