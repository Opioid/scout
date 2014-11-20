package file

import (
	"os"
	"bytes"
	_ "fmt"
)

const (
	Unknown = iota
	JPG     = iota
	PNG     = iota
	SUI     = iota
)

func QueryFileType(fi *os.File) int {
	header := make([]byte, 4)

	_, err := fi.ReadAt(header, 0)

	if err != nil {
		return Unknown
	}

	if bytes.HasPrefix(header, []byte{255, 216}) {
		return JPG
	} else if bytes.Compare(header, []byte{137, 'P', 'N', 'G'}) == 0 {
		return PNG
	} else if bytes.Compare(header, []byte{'S', 'U', 'I'}) == 0 {
		return SUI
	}

	return Unknown
}