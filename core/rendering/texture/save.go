package texture

import (
	"io"
	"encoding/binary"
)

func Save(w io.Writer, t *Texture2D) error {
	w.Write([]byte{'S', 'U', 'I'})

	if err := binary.Write(w, binary.LittleEndian, t.texture.description); err != nil {
		return err
	}

	return nil
}