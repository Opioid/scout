package texture

import (
	"io"
	"bytes"
	"encoding/binary"
	"errors"
	_ "fmt"
)

func Load(r io.Reader) (*Texture2D, error) {
	header := make([]byte, 3)

	if _, err := r.Read(header); err != nil {
		return nil, err
	}

	if bytes.Compare(header, []byte{'S', 'U', 'I'}) != 0 {
		return nil, errors.New("Header does not match SUI.")
	}

	description := description{}
	if err := binary.Read(r, binary.LittleEndian, &description); err != nil {
		return nil, err
	}

	texture := NewTexture2DFromDescription(&description)

	buffers := texture.Image.Buffers
	for i := range texture.Image.Buffers {
		if err := binary.Read(r, binary.LittleEndian, buffers[i].Data()); err != nil {
			return nil, err
		}
	}

/*
	if err := binary.Write(w, binary.LittleEndian, t.description()); err != nil {
		return err
	}

	buffers := t.Image.Buffers
	for i := range buffers {
		if err := binary.Write(w, binary.LittleEndian, buffers[i].Data()); err != nil {
			return err
		}
	}
*/
	return texture, nil
}