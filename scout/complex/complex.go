package complex

import (
	pkgscene "github.com/Opioid/scout/core/scene"
)

func Init(scene *pkgscene.Scene) {
	scene.ComplexProvider.Register("Sphere_block", &SphereBlockFactory{})
}