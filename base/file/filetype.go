package file

import (
	"os"
	"bytes"
	"path/filepath"
	_ "fmt"
)

const (
	Unknown = iota
	JPG     = iota
	PNG     = iota
	RGBE    = iota
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
	} else if bytes.HasPrefix(header, []byte{137, 'P', 'N', 'G'}) {
		return PNG
	} else if bytes.HasPrefix(header, []byte{'#', '?' }) {
		return RGBE
	} else if bytes.HasPrefix(header, []byte{'S', 'U', 'I'}) {
		return SUI
	}

	return Unknown
}

func WithoutExt(path string) string {
	ext := filepath.Ext(path)
	return path[0:len(path) - len(ext)]
}