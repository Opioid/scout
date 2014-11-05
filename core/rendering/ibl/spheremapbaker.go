package ibl

import (
	"github.com/Opioid/scout/core/rendering/texture"
	"github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/base/math"
	"os"
	"image/png"
)

func BakeSphereMap(s surrounding.Surrounding) {
	buffer := texture.Buffer{}

	buffer.Resize(math.MakeVector2i(512, 128))

	image := buffer.RGBA()

	fo, err := os.Create("sphere_map.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
}