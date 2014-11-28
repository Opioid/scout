package rgbe

// http://www.graphics.cornell.edu/~bjw/rgbe

import (
	"io"
	"bufio"
	_ "bytes"
	_ "encoding/binary"
	"errors"
	"fmt"
)

func Decode(r io.Reader) (int, int, []float32, error) {
	br := bufio.NewReader(r)

	width, height, err := readHeader(br) 

	if err != nil {
		return 0, 0, nil, err
	} 

	data := make([]float32, width * height * 3)

	if err := readPixels_RLE(br, width, height, &data); err != nil {
		return 0, 0, nil, err
	}

	return width, height, data, nil
}

func readHeader(r *bufio.Reader) (int, int, error) {
	line, err := r.ReadString('\n')

	if err != nil {
		return 0, 0, err
	}

	if line[0] != '#' || line[1] != '?' {
		return 0, 0, errors.New("Bad initial token")
	}

	for {
		if line[0] == 0 || line[0] == '\n' {
			return 0, 0, errors.New("No FORMAT specifier found")
		} else if line == "FORMAT=32-bit_rle_rgbe\n" {
			break
		}

		line, err = r.ReadString('\n')

		if err != nil {
			return 0, 0, err
		}
	}	

	line, err = r.ReadString('\n')

	if err != nil {
		return 0, 0, err
	}

	if line[0] != '\n' {
		errors.New("Missing blank line after FORMAT specifier")
	}

	line, err = r.ReadString('\n')

	if err != nil {
		return 0, 0, err
	}

	var width, height int
	if n, err := fmt.Sscanf(line, "-Y %d +X %d", &height, &width); n < 2 || err != nil {
		return 0, 0, errors.New("Missing image size specifier")
	}

	return width, height, nil
}

func readPixels_RLE(r *bufio.Reader, scanlineWidth, numScanlines int, data *[]float32) error {
	fmt.Println("readPixels_RLE()")

	if scanlineWidth < 8 || scanlineWidth > 0x7fff {
		// run length encoding is not allowed so read flat
		return readPixels(r, 0, scanlineWidth * numScanlines, data)
	} 

	offset := 0
	rgbe := make([]byte, 4)
	scanlineBuffer := make([]byte, 4 * scanlineWidth)
	buf := make([]byte, 2)

	for numScanlines > 0 {
		if _, err := r.Read(rgbe); err != nil {
			return err
		}

		if rgbe[0] != 2 || rgbe[1] != 2 || (rgbe[2] & 0x80) != 0 {
			// this file is not run length encoded
			return readPixels(r, offset + 3, scanlineWidth * numScanlines - 1, data)
		}

		if int(rgbe[2]) << 8 | int(rgbe[3]) != scanlineWidth {
			return errors.New("Wrong scanline width")
		}

		// read each of the four channels for the scanline into the buffer 
		index := 0
		for i := 0; i < 4; i++ {
			end := (i + 1) * scanlineWidth

			for index < end {
				if _, err := r.Read(buf); err != nil {
					return err
				}

				if buf[0] > 128 {
					 // a run of the same value
					count := int(buf[0]) - 128

					if count == 0 || count > end - index {
						return errors.New("Bad scanline data 0")
					}

					for ; count > 0; count-- {
						scanlineBuffer[index] = buf[1]
						index++
					}

				} else {
					// a non-run
					count := int(buf[0])

					if count == 0 || count > end - index {
						return errors.New("Bad scanline data 1")
					}

					scanlineBuffer[index] = buf[1]
					index++

					count--
					if count > 0 {
						tmp := make([]byte, count)
						if _, err := r.Read(/*scanlineBuffer[index:index+count]*/tmp); err != nil {
							return err
						}

						index += count
					}
				}
			}
		}

		numScanlines--
	}


	return nil
}

func readPixels(r *bufio.Reader, offset, numPixels int, data *[]float32) error {
	fmt.Println("readPixels()")

	rgbe := make([]byte, 4)

	for ; numPixels > 0; numPixels-- {
		if _, err := r.Read(rgbe); err != nil {
			return err
		}

		(*data)[offset] = 1
		offset++

		(*data)[offset] = 0
		offset++

		(*data)[offset] = 1
		offset++

	}

	return nil
}