package texture

import (
	"github.com/Opioid/scout/base/math"
	gomath "math"
	"io"
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	_ "fmt"
)

func Load(r io.Reader) (*Texture2D, error) {
	br := bufio.NewReader(r)

	header := make([]byte, 3)

	if _, err := io.ReadFull(br, header); err != nil {
		return nil, err
	}

	if bytes.Compare(header, []byte{'S', 'U', 'I'}) != 0 {
		return nil, errors.New("Header does not match SUI.")
	}

	description := description{}
	if err := binary.Read(br, binary.LittleEndian, &description); err != nil {
		return nil, err
	}

	texture := NewTexture2DFromDescription(&description)

	buf := make([]byte, 4 * 4)

	buffers := texture.Image.Buffers
	for i := range buffers {
		dimensions := buffers[i].Dimensions()
		for y := int32(0); y < dimensions.Y; y++ {
			for x := int32(0); x < dimensions.X; x++ {
				if _, err := io.ReadFull(br, buf); err != nil {
					return nil, err
				}

				r := gomath.Float32frombits(binary.LittleEndian.Uint32(buf[ 0: 4]))
				g := gomath.Float32frombits(binary.LittleEndian.Uint32(buf[ 4: 8]))
				b := gomath.Float32frombits(binary.LittleEndian.Uint32(buf[ 8:12]))
				a := gomath.Float32frombits(binary.LittleEndian.Uint32(buf[12:16]))

				buffers[i].Set(x, y, math.MakeVector4(r, g, b, a))
			}
		}
	}

	return texture, nil
}