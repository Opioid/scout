package bvh

import (
	"github.com/Opioid/scout/core/scene/prop"
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/math/bounding"
)

type Builder struct {
//	nodes []miniNode
//	numNodes, currentNode uint32
}


func (b *Builder) Build(props []*prop.StaticProp, maxShapes int, tree *Tree, outProps *[]*prop.StaticProp) {
	*outProps = make([]*prop.StaticProp, 0, len(props))

	root := buildNode{}

	root.split(props, maxShapes, outProps)
/*
	b.numNodes = 1
	root.numSubNodes(&b.numNodes)

	fmt.Println(b.numNodes)

	b.nodes = tree.allocateNodes(b.numNodes)

	b.currentNode = 0
	b.serialize(&root)
*/
	tree.root = root

	tree.infinitePropsBegin = uint32(len(*outProps))
	tree.infinitePropsEnd   = uint32(len(props))

	for _, p := range props {
		if !p.Shape.IsFinite() {
			*outProps = append(*outProps, p)
		}
	}

}

func (b *Builder) serialize(node *buildNode) {
/*
	AABB_node* n = get_new_node();

	n->aabb_ = node->aabb;
	n->props_.swap(node->props);

	if (node->child0)
	{
		serialize(node->child0);

		serialize(node->child1);

		n->set_has_children(true);
	}

	n->set_skip_node(skip_node());
*/
/*
	n := b.newNode()
	n.aabb = node.aabb
	n.indices = node.indices

	if node.children[0] != nil {
		b.serialize(node.children[0])
		b.serialize(node.children[1])

		n.setHasChildren(true)
	}

	n.setSkipOffset(b.skipOffset())
	*/
}

/*
func (b *Builder) newNode() *miniNode {
	n := &b.nodes[b.currentNode]
	b.currentNode++
	return n
}

func (b *Builder) skipOffset() uint32 {
	if b.currentNode >= b.numNodes {
		return 0
	}

	return b.currentNode
}
*/

type buildNode struct {
	aabb bounding.AABB

	axis int

	offset   uint32
	propsEnd uint32

	children [2]*buildNode
}

func (n *buildNode) split(props []*prop.StaticProp, maxShapes int, outProps *[]*prop.StaticProp) {
	n.aabb = miniaabb(props)

	if len(props) <= maxShapes {
		n.assign(props, outProps)
	} else {
		n.children[0] = new(buildNode)
		n.children[1] = new(buildNode)

		numProps := len(props) / 2
		props0 := make([]*prop.StaticProp, 0, numProps)
		props1 := make([]*prop.StaticProp, 0, numProps)

		sp, axis := chooseSplittingPlane(&n.aabb)

		n.axis = axis

		for _, p := range props {
			if !p.Shape.IsFinite() {
				continue
			}

			mib := sp.Behind(p.AABB.Bounds[0])
			mab := sp.Behind(p.AABB.Bounds[1])
			if mib && mab {
				props0 = append(props0, p)
			} else {
				props1 = append(props1, p)
			}
		}		

		n.children[0].split(props0, maxShapes, outProps)
		n.children[1].split(props1, maxShapes, outProps)
	}
}

func (n *buildNode) assign(props []*prop.StaticProp, outProps *[]*prop.StaticProp) {
	n.offset = uint32(len(*outProps))

	for _, p := range props {
		*outProps = append(*outProps, p)
	}	

	n.propsEnd = uint32(len(*outProps))
}

func (n *buildNode) numSubNodes(num *uint32) {
	if n.children[0] != nil {
		*num += 2

		n.children[0].numSubNodes(num);
		n.children[1].numSubNodes(num);
	}
}

func (n *buildNode) intersect(ray *math.OptimizedRay, props []*prop.StaticProp, intersection *prop.Intersection) bool {
	if !n.aabb.Intersect(ray) {
		return false
	}

	hit := false

	if n.children[0] != nil {
		c0, c1 := 0, 1

		if ray.DirIsNeg[n.axis] == 1 {
			c0 = 1
			c1 = 0
		}

		if n.children[c0].intersect(ray, props, intersection) {
			hit = true
		} 

		if n.children[c1].intersect(ray, props, intersection) {
			hit = true
		}
	} else {
		for i := n.offset; i < n.propsEnd; i++ {
			p := props[i]
			if p.Intersect(ray, intersection) {
				intersection.Prop = &p.Prop
				hit = true
			}
		}
	}

	return hit
}

func (n *buildNode) intersectP(ray *math.OptimizedRay, props []*prop.StaticProp) bool {
	if !n.aabb.Intersect(ray) {
		return false
	}

	if n.children[0] != nil {
		if n.children[0].intersectP(ray, props) {
			return true
		} 

		return n.children[1].intersectP(ray, props)
	}

	for i := n.offset; i < n.propsEnd; i++ {
		if props[i].IntersectP(ray) {
			return true
		}
	}

	return false
}

func miniaabb(props []*prop.StaticProp) bounding.AABB {
	b := bounding.MakeEmptyAABB()

	for _, p := range props {
		b = b.Merge(&p.AABB)
	}

	return b
}

func chooseSplittingPlane(aabb *bounding.AABB) (math.Plane, int) {
	position := aabb.Position()
	halfsize := aabb.Halfsize()

	if halfsize.X >= halfsize.Y && halfsize.X >= halfsize.Z {
		return math.MakePlane(math.Vector3{1.0, 0.0, 0.0}, position), 0
	} else if halfsize.Y >= halfsize.X && halfsize.Y >= halfsize.Z {
		return math.MakePlane(math.Vector3{0.0, 1.0, 0.0}, position), 1
	} else {
		return math.MakePlane(math.Vector3{0.0, 0.0, 1.0}, position), 2
	}
}