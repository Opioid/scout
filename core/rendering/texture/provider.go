package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/rendering/color"
	"github.com/Opioid/scout/base/file"
	"os"
	goimage "image"
	_ "image/jpeg"
	"runtime"
	"sync"
	"fmt"
)

type Provider struct {

}

func (p *Provider) Load2D(filename string) *Texture2D {
	fi, err := os.Open(filename)

	fmt.Println(file.QueryFileType(fi))

	if err != nil {
		fmt.Printf("%s could not be loaded", filename)
		return nil
	}

	defer fi.Close()

	sourceImage, _, err := goimage.Decode(fi)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	dimensions := math.MakeVector2i(sourceImage.Bounds().Max.X, sourceImage.Bounds().Max.Y)

	texture := NewTexture2D(dimensions, 1)

	numTaks := runtime.GOMAXPROCS(0)

	a := dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := 0; i < numTaks; i++ {
		wg.Add(1)

		go func (start, end math.Vector2i) {
			process(start, end, sourceImage, &texture.Image.Buffers[0])
			wg.Done()
		}(start, end)

		start.Y += a

		if i == numTaks - 2 {
			end.Y = dimensions.Y
		} else {
			end.Y += a
		}
	}

	wg.Wait()

	return texture
}

func process(start, end math.Vector2i, source goimage.Image, target *Buffer) {
	max := float32(0xFFFF)

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			r, g, b, a := source.At(x, y).RGBA()

			target.Set(x, y, math.MakeVector4(
				color.SrgbToLinear(float32(r) / max), 
				color.SrgbToLinear(float32(g) / max), 
				color.SrgbToLinear(float32(b) / max), 
				float32(a) / max,
			))
		}
	}
}


/*

	numTaks := runtime.GOMAXPROCS(0)

	a := f.dimensions.Y / numTaks

	start := math.Vector2i{0, 0}
	end   := math.Vector2i{f.dimensions.X, a}

	wg := sync.WaitGroup{}

	for i := 0; i < numTaks; i++ {
		wg.Add(1)

		go func (s, e math.Vector2i, t *image.RGBA) {
			f.process(s, e, t)
			wg.Done()
		}(start, end, target)

		start.Y += a

		if i == numTaks - 1 {
			end.Y = f.dimensions.Y
		} else {
			end.Y += a
		}
	}

	wg.Wait()

	return target
}

func (f *Unfiltered) process(start, end math.Vector2i, target *image.RGBA) {
	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			pixel := f.at(x, y)
			iw := 1.0 / pixel.weightSum
			r := uint8(255.0 * color.LinearToSrgb(pixel.color.X * iw))
			g := uint8(255.0 * color.LinearToSrgb(pixel.color.Y * iw))
			b := uint8(255.0 * color.LinearToSrgb(pixel.color.Z * iw))

			target.Set(x, y, gocolor.RGBA{r, g, b, 255})
		}
	}
}
*/