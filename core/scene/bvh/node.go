package bvh

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type node struct {
	aabb bounding.AABB

	props []*prop.StaticProp

	children [2]*node
}

func (n *node) split(props []*prop.StaticProp) {
	n.aabb = aabb(props)

	if len(props) < 32 {
		n.assign(props)
	} else {
		n.children[0] = new(node)
		n.children[1] = new(node)

		numProps := len(props) / 2
		props0 := make([]*prop.StaticProp, 0, numProps)
		props1 := make([]*prop.StaticProp, 0, numProps)

		sp := chooseSplittingPlane(&n.aabb)

		for _, p := range props {
			mib := sp.Behind(p.AABB.Min)
			mab := sp.Behind(p.AABB.Max)
			if mib && mab {
				props0 = append(props0, p)
			} else if !mib && !mab {
				props1 = append(props1, p)
			} else {
				props0 = append(props0, p)
				props1 = append(props1, p)
			}
		}

		n.children[0].split(props0)
		n.children[1].split(props1)
	}
}

func (n *node) assign(props []*prop.StaticProp) {
	n.props = props
}

func (n *node) intersect(ray *math.OptimizedRay, intersection *prop.Intersection) bool {
	if !n.aabb.Intersect(ray) {
		return false
	}

	hit := false

	if n.children[0] != nil {
		sd0 := n.children[0].aabb.Position().SquaredDistance(ray.Origin)
		sd1 := n.children[1].aabb.Position().SquaredDistance(ray.Origin)

		var c0, c1 int

		if sd0 <= sd1 {
			c0 = 0
			c1 = 1
		} else {
			c0 = 1
			c1 = 0
		}

		if n.children[c0].intersect(ray, intersection) {
			hit = true
		} 

		if n.children[c1].intersect(ray, intersection) {
			hit = true
		}
	} else {
		for _, p := range n.props {
			if p.Intersect(ray, intersection) {
				intersection.Prop = &p.Prop
				hit = true
			}
		}
	}

	return hit
}

func (n *node) intersectP(ray *math.OptimizedRay) bool {
	if !n.aabb.Intersect(ray) {
		return false
	}

	if n.children[0] != nil {
		if n.children[0].intersectP(ray) {
			return true
		} 

		return n.children[1].intersectP(ray)
	}

	for _, p := range n.props {
		if p.IntersectP(ray) {
			return true
		}
	}

	return false
}

func aabb(props []*prop.StaticProp) bounding.AABB {
	b := bounding.MakeAABB()

	for _, p := range props {
		b = b.Merge(&p.AABB)
	}

	return b
}

func chooseSplittingPlane(aabb *bounding.AABB) math.Plane {
	position := aabb.Position()
	halfsize := aabb.Halfsize()

	if halfsize.X >= halfsize.Y && halfsize.X >= halfsize.Z {
		return math.MakePlane(math.Vector3{1.0, 0.0, 0.0}, position)
	} else if halfsize.Y >= halfsize.X && halfsize.Y >= halfsize.Z {
		return math.MakePlane(math.Vector3{0.0, 1.0, 0.0}, position)
	} else {
		return math.MakePlane(math.Vector3{0.0, 0.0, 1.0}, position)
	}
}