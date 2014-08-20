package complex

import (
	pkgscene "github.com/Opioid/scout/core/scene"
)

type sphereBlock struct {

}

func (c *sphereBlock) Init(scene *pkgscene.Scene) {
}

type SphereBlockFactory struct {

}

func (f *SphereBlockFactory) New() pkgscene.Complex {
	return new(sphereBlock)
}