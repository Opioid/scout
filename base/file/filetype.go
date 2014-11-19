package file

import (
	"os"
	"bytes"
	"fmt"
)

func QueryFileType(fi *os.File) string {
	header := make([]byte, 3)

	n, err := fi.ReadAt(header, 0)

	if err != nil {
		return "Unknown"
	}

	if bytes.HasPrefix(header, []byte{255, 216}) {
		return "JPG"
	} else if bytes.Compare(header, []byte{'P', 'N', 'G'}) == 0 {
		return "PNG"
	}

	fmt.Println(header)

	return string(header[:n])
}