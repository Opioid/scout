package texture

import (
	"github.com/Opioid/scout/base/math"
	"github.com/Opioid/scout/base/rendering/color"
	"github.com/Opioid/scout/base/file"
	"os"
	goimage "image"
	_ "image/jpeg"
	"github.com/Opioid/rgbe"
	"runtime"
	"sync"
	"fmt"
)

type Provider struct {

}

func (p *Provider) Load2D(filename string, treatAsLinear, encodedFloats bool) *Texture2D {
	fi, err := os.Open(filename)

	defer fi.Close()

	if err != nil {
		fmt.Printf("%s could not be loaded", filename)
		return nil
	}

	fileType := file.QueryFileType(fi)

	if file.RGBE == fileType {
		return textureFromRgbe(fi)
	} else {
		return textureFromGoSupportedFile(fi, treatAsLinear, encodedFloats)
	}
}

func textureFromRgbe(fi *os.File) *Texture2D {
	width, height, data, err := rgbe.Decode(fi)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	dimensions := math.MakeVector2i(int32(width), int32(height))

	texture := NewTexture2D(dimensions, 1)

	for y := int32(0); y < dimensions.Y; y++ {
		for x := int32(0); x < dimensions.X; x++ {
			o := (dimensions.X * y + x) * 3

			r := data[o + 0]
			g := data[o + 1]
			b := data[o + 2]

			texture.Image.Buffers[0].Set(x, y, math.MakeVector4(r, g, b, 1))
		}
	}

	return texture
}

func textureFromGoSupportedFile(fi *os.File, treatAsLinear, encodedFloats bool) *Texture2D {
	sourceImage, _, err := goimage.Decode(fi)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	dimensions := math.MakeVector2i(int32(sourceImage.Bounds().Max.X), int32(sourceImage.Bounds().Max.Y))

	texture := NewTexture2D(dimensions, 1)

	numTaks := int32(runtime.GOMAXPROCS(0))

	a := dimensions.Y / numTaks

	start := math.MakeVector2i(0, 0)
	end   := math.MakeVector2i(dimensions.X, a)

	wg := sync.WaitGroup{}

	for i := int32(0); i < numTaks; i++ {
		wg.Add(1)

		go func (start, end math.Vector2i) {
			if encodedFloats {
				processEncodedFloats(start, end, sourceImage, &texture.Image.Buffers[0])
			} else if treatAsLinear {
				processLinear(start, end, sourceImage, &texture.Image.Buffers[0])
			} else {
				processSrgbToLinear(start, end, sourceImage, &texture.Image.Buffers[0])
			}

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

func processSrgbToLinear(start, end math.Vector2i, source goimage.Image, target *Buffer) {
	max := float32(0xFFFF)

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			r, g, b, a := source.At(int(x), int(y)).RGBA()

			target.Set(x, y, math.MakeVector4(
				color.SrgbToLinear(float32(r) / max), 
				color.SrgbToLinear(float32(g) / max), 
				color.SrgbToLinear(float32(b) / max), 
				float32(a) / max,
			))
		}
	}
}

func processLinear(start, end math.Vector2i, source goimage.Image, target *Buffer) {
	max := float32(0xFFFF)

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			r, g, b, a := source.At(int(x), int(y)).RGBA()

			target.Set(x, y, math.MakeVector4(
				float32(r) / max, 
				float32(g) / max, 
				float32(b) / max, 
				float32(a) / max,
			))
		}
	}
}

func processEncodedFloats(start, end math.Vector2i, source goimage.Image, target *Buffer) {
	max := float32(0xFFFF)

	for y := start.Y; y < end.Y; y++ {
		for x := start.X; x < end.X; x++ {
			r, g, b, a := source.At(int(x), int(y)).RGBA()

			target.Set(x, y, math.MakeVector4(
				2 * (float32(r) / max - 0.5), 
				2 * (float32(g) / max - 0.5), 
				2 * (float32(b) / max - 0.5), 
				float32(a) / max,
			))
		}
	}
}