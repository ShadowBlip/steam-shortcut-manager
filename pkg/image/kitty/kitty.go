package kitty

import (
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/dolmen-go/kittyimg"
)

func Display(file string) error {
	img, err := readImageFile(file)
	if err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}

	err = kittyimg.Fprintln(os.Stdout, img)
	if err != nil {
		return err
	}
	return nil
}

func readImageFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func icat(w io.Writer, img image.Image) {
	const chunkEncSize = 4096
	// const chunkEncSize = 48
	const chunkRawSize = (chunkEncSize / 4) * 3

	bounds := img.Bounds()

	// f=32 => RGBA
	fmt.Fprintf(w, "\033_Gq=1,a=T,f=32,s=%d,v=%d,t=d,", bounds.Dx(), bounds.Dy())

	bufRaw := make([]byte, 0, chunkRawSize)
	bufEnc := make([]byte, chunkEncSize)

	flush := func(last bool) {
		if len(bufRaw) == 0 {
			w.Write([]byte("m=0;\033\\"))
			return
		}
		if last {
			w.Write([]byte("m=0;"))
		} else {
			w.Write([]byte("m=1;"))
		}

		// fmt.Fprintln(os.Stderr, len(bufRaw), "=>", (len(bufRaw)+2)/3*4)

		base64.StdEncoding.Encode(bufEnc, bufRaw)
		w.Write(bufEnc[:(len(bufRaw)+2)/3*4])

		if last {
			w.Write([]byte("\033\\"))
		} else {
			w.Write([]byte("\033\\\033_G"))
			bufRaw = bufRaw[:0]
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if len(bufRaw)+4 > chunkRawSize {
				flush(false)
			}
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			// Shifting by 8 reduces this to the range [0, 255].
			bufRaw = append(bufRaw, byte(r>>8), byte(g>>8), byte(b>>8), byte(a>>8))
		}
	}
	flush(true)
}
