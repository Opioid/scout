package ibl

import (
	"github.com/Opioid/scout/core/rendering/texture"
	_ "github.com/Opioid/scout/core/rendering/material/ggx"
	_ "github.com/Opioid/scout/core/scene/surrounding"
	"github.com/Opioid/scout/base/math"
	_ "github.com/Opioid/scout/base/math/random"
	_ "math"
	"os"
	"image/png"
	_ "runtime"
	_ "sync"	
	_ "strconv"
	_ "fmt"
)

func IntegrateGgxBrdf(numSamples uint32, buffer *texture.Buffer) {
	dimensions := buffer.Dimensions()


	for y := 0; y < dimensions.Y; y++ {

		for x := 0; x < dimensions.X; x++ {


			buffer.Set(x, y, math.MakeVector4(1, 1, 1, 1))
		}

		
	}

	image := buffer.RGBA()

	fo, err := os.Create("ggx_brdf.png")

	if err != nil {
		panic(err)
	}

	defer fo.Close()

	png.Encode(fo, image)
}