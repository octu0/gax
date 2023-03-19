package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	"github.com/octu0/gax"
)

func main() {
	img, err := readImage("./src.png")
	if err != nil {
		panic(err)
	}

	in := gax.ImageRGBA[uint8](img)

	blurX := gax.Func[uint8]()
	blurX.SetXYC(func(x, y, ch gax.Var) uint8 {
		if int(ch) == 3 {
			return 255
		}

		val1 := int(in.XYC(x-1, y, ch))
		val2 := int(in.XYC(x, y, ch))
		val3 := int(in.XYC(x+1, y, ch))
		return uint8((val1 + val2 + val3) / 3)
	})

	blurY := gax.Func[uint8]()
	blurY.SetXYC(func(x, y, ch gax.Var) uint8 {
		if int(ch) == 3 {
			return 255
		}

		val1 := int(blurX.XYC(x, y-1, ch))
		val2 := int(blurX.XYC(x, y, ch))
		val3 := int(blurX.XYC(x, y+1, ch))
		return uint8((val1 + val2 + val3) / 3)
	})

	out := blurY.Realize(img.Rect.Dx(), img.Rect.Dy(), 4) // 4 = r,g,b,a
	outImg := &image.RGBA{
		Pix:    out,
		Stride: 4 * img.Rect.Dx(),
		Rect:   image.Rect(0, 0, img.Rect.Dx(), img.Rect.Dy()),
	}

	path, err := saveImage(outImg)
	if err != nil {
		panic(err)
	}
	println("blur", path)
}

func readImage(path string) (*image.RGBA, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	rgba, err := pngToRGBA(data)
	if err != nil {
		return nil, err
	}
	return rgba, nil
}

func saveImage(img *image.RGBA) (string, error) {
	out, err := os.CreateTemp("/tmp", "out*.png")
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		return "", err
	}
	return out.Name(), nil
}

func pngToRGBA(data []byte) (*image.RGBA, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if i, ok := img.(*image.RGBA); ok {
		return i, nil
	}

	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y += 1 {
		for x := b.Min.X; x < b.Max.X; x += 1 {
			c := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)
			rgba.Set(x, y, c)
		}
	}
	return rgba, nil
}
