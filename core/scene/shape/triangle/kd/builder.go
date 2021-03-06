package kd
/*
import (
	"github.com/Opioid/scout/core/scene/shape/geometry"
	"github.com/Opioid/scout/core/scene/shape/triangle/primitive"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
	gomath "math"
	_ "fmt"
)

type Builder struct {

}

func (b *Builder) Build(indices []uint32, vertices []geometry.Vertex, maxPrimitives int, tree *Tree) {
	primitiveIndices := make([]uint32, len(indices) / 3)

	for i := range primitiveIndices {
		primitiveIndices[i] = uint32(i) * 3
	}

	root := buildNode{}
	root.split(primitiveIndices, indices, vertices, maxPrimitives, 0)

	tree.root = root
}

type buildNode struct {
	axis int32
	splitPos float32

//	indices []uint32
	triangles []primitive.Triangle

	children [2]*buildNode
}

func (n *buildNode) split(primitiveIndices, indices []uint32, vertices []geometry.Vertex, maxPrimitives, depth int) {
	if len(primitiveIndices) < maxPrimitives || depth > 10 {
		n.assign(primitiveIndices, indices, vertices)
	} else {
		b := subMeshAabb(primitiveIndices, indices, vertices)

		n.axis, n.splitPos = splittingPlane(&b)

		p := n.plane()

		numPids := len(primitiveIndices) / 2
		pids0 := make([]uint32, 0, numPids)
		pids1 := make([]uint32, 0, numPids)

		for _, pi := range primitiveIndices {
			s := triangleSide(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, p)
			
			if s == 0 {
				pids0 = append(pids0, pi)
			} else if s == 1 {
				pids1 = append(pids1, pi)
			} else {
				pids0 = append(pids0, pi)
				pids1 = append(pids1, pi)
			}
		}

		n.children[0] = new(buildNode)
		n.children[1] = new(buildNode)

		n.children[0].split(pids0, indices, vertices, maxPrimitives, depth + 1)
		n.children[1].split(pids1, indices, vertices, maxPrimitives, depth + 1)
	}
}

func (n *buildNode) assign(primitiveIndices []uint32, indices []uint32, vertices []geometry.Vertex) {
//	n.indices = primitiveIndices

	n.triangles = make([]primitive.Triangle, len(primitiveIndices))

	for i, pi := range primitiveIndices {
		n.triangles[i] = primitive.MakeTriangle(&vertices[indices[pi + 0]], &vertices[indices[pi + 1]], &vertices[indices[pi + 2]])
	}
}


var axis = [...]math.Vector3{ 
	math.MakeVector3(1, 0, 0),
	math.MakeVector3(0, 1, 0),
	math.MakeVector3(0, 0, 1), 
}

func (n *buildNode) plane() math.Plane {
	return math.Plane{A: axis[n.axis].X, B: axis[n.axis].Y, C: axis[n.axis].Z, D: n.splitPos}
}

func (n *buildNode) intersect(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, indices []uint32, vertices []geometry.Vertex, intersection *primitive.Intersection) bool {
	if intersection.T < boundingMinT || ray.MinT > boundingMaxT {
		return false
	}

	hit := false

	if n.children[0] != nil {
		oa := ray.Origin.At(n.axis)

		tplane := (n.splitPos + oa) * -ray.ReciprocalDirection.At(n.axis)

	//	c := ray.DirIsNeg[n.axis]

		c := 1

		if oa < -n.splitPos || (oa == -n.splitPos && ray.Direction.At(n.axis) >= 0.0) {
			c = 0
		}

		if tplane > boundingMaxT || tplane <= 0.0 {
			if n.children[c].intersect(ray, boundingMinT, boundingMaxT, indices, vertices, intersection) {
				hit = true
			} 
		} else if tplane < boundingMinT {
			if n.children[1 - c].intersect(ray, boundingMinT, boundingMaxT, indices, vertices, intersection) {
				hit = true
			}
		} else {
			if n.children[c].intersect(ray, boundingMinT, tplane, indices, vertices, intersection) {
				hit = true
			} 

			if n.children[1 - c].intersect(ray, tplane, boundingMaxT, indices, vertices, intersection) {
				hit = true
			}
		}
	} else {

		for _, pi := range n.indices {
			ti := Intersection{Index: pi}
			if intersectTriangle(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, ray, &ti.T, &ti.U, &ti.V) {
				if ti.T <= intersection.T {
					*intersection = ti
					hit = true
				}
			}
		}

		var ti primitive.Intersection
		var index int

		for i, t := range n.triangles {
			if t.Intersect(ray, &ti.T, &ti.U, &ti.V) {
				if ti.T <= intersection.T {
					*intersection = ti
					index = i
					hit = true
				}
			}
		}

		intersection.Triangle = &n.triangles[index]
		
	}

	return hit
}

func (n *buildNode) intersectP(ray *math.OptimizedRay, boundingMinT, boundingMaxT float32, indices []uint32, vertices []geometry.Vertex) bool {
	if ray.MaxT < boundingMinT || ray.MinT > boundingMaxT {
		return false
	}

	if n.children[0] != nil {
		oa := ray.Origin.At(n.axis)

		tplane := (n.splitPos + ray.Origin.At(n.axis)) * -ray.ReciprocalDirection.At(n.axis)

	//	c := ray.DirIsNeg[n.axis]

		c := 1

		if oa < -n.splitPos || (oa == -n.splitPos && ray.Direction.At(n.axis) >= 0.0) {
			c = 0
		}

		if tplane > boundingMaxT || tplane <= 0.0 {
			return n.children[c].intersectP(ray, boundingMinT, boundingMaxT, indices, vertices)
		} else if tplane < boundingMinT {
			return n.children[1 - c].intersectP(ray, boundingMinT, boundingMaxT, indices, vertices)
		} else {
			if n.children[c].intersectP(ray, boundingMinT, tplane, indices, vertices) {
				return true
			} 

			return n.children[1 - c].intersectP(ray, tplane, boundingMaxT, indices, vertices)
		}
	}

	for _, pi := range n.indices {
		if intersectTriangleP(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, ray) {
			return true
		}
	}

	for _, t := range n.triangles {
		if t.IntersectP(ray) {
			return true
		}
	}

	return false

}

func subMeshAabb(primitiveIndices, indices []uint32, vertices []geometry.Vertex) bounding.AABB {
	min := math.MakeVector3( gomath.MaxFloat32,  gomath.MaxFloat32,  gomath.MaxFloat32)
	max := math.MakeVector3(-gomath.MaxFloat32, -gomath.MaxFloat32, -gomath.MaxFloat32)
	
	for _, pi := range primitiveIndices {
		min = triangleMin(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, min)
		max = triangleMax(vertices[indices[pi + 0]].P, vertices[indices[pi + 1]].P, vertices[indices[pi + 2]].P, max)
	}

	return bounding.MakeAABB(min, max)
}

func triangleMin(a, b, c, x math.Vector3) math.Vector3 {
	return a.Min(b).Min(c).Min(x)
}

func triangleMax(a, b, c, x math.Vector3) math.Vector3 {
	return a.Max(b).Max(c).Max(x)
}

func triangleSide(a, b, c math.Vector3, p math.Plane) int {
	behind := 0

	if p.Behind(a) {
		behind++
	}

	if p.Behind(b) {
		behind++
	}

	if p.Behind(c) {
		behind++
	}

	if behind == 3 {
		return 0
	} else if behind == 0 {
		return 1
	} else {
		return 2
	}
}

func splittingPlane(aabb *bounding.AABB) (int32, float32) {
	position := aabb.Position()
	halfsize := aabb.Halfsize()

	if halfsize.X >= halfsize.Y && halfsize.X >= halfsize.Z {
		p := math.MakePlane(axis[0], position)
		return 0, p.D
	} else if halfsize.Y >= halfsize.X && halfsize.Y >= halfsize.Z {
		p := math.MakePlane(axis[1], position)
		return 1, p.D
	} else {
		p := math.MakePlane(axis[2], position)
		return 2, p.D
	}
}*/