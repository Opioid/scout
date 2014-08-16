package scene

import (
	"github.com/Opioid/scout/core/scene/shape"
	"github.com/Opioid/scout/base/math/bounding"
)

type Prop struct {
	Shape shape.Shape
	Material Material
	AABB bounding.AABB
}
