package texture

import (
	_ "github.com/Opioid/scout/base/math"
	gomath "math"
	"io"
	"bufio"
	"encoding/binary"
)

func Save(w io.Writer, t *Texture2D) error {
	bw := bufio.NewWriter(w)

	bw.Write([]byte{'S', 'U', 'I'})

	if err := binary.Write(bw, binary.LittleEndian, t.description()); err != nil {
		return err
	}

	buf := make([]byte, 4 * 4)

	buffers := t.Image.Buffers
	for i := range buffers {
		dimensions := buffers[i].Dimensions()
		for y := int32(0); y < dimensions.Y; y++ {
			for x := int32(0); x < dimensions.X; x++ {
				c := buffers[i].At(x, y)

				binary.LittleEndian.PutUint32(buf[ 0: 4], gomath.Float32bits(c.X))
				binary.LittleEndian.PutUint32(buf[ 4: 8], gomath.Float32bits(c.Y))
				binary.LittleEndian.PutUint32(buf[ 8:12], gomath.Float32bits(c.Z))
				binary.LittleEndian.PutUint32(buf[12:16], gomath.Float32bits(c.W))

				if _, err := bw.Write(buf); err != nil {
					return err
				}
			}
		}
	}

	bw.Flush()

	return nil
}