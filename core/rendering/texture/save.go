package texture

import (
	"io"
	"encoding/binary"
)

func Save(w io.Writer, t *Texture2D) error {
	w.Write([]byte{'S', 'U', 'I'})

	if err := binary.Write(w, binary.LittleEndian, t.description()); err != nil {
		return err
	}

	buffers := t.Image.Buffers
	for i := range buffers {
		if err := binary.Write(w, binary.LittleEndian, buffers[i].Data()); err != nil {
			return err
		}
	}

	return nil
}