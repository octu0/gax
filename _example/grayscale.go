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

const (
	grayR = 76
	grayG = 152
	grayB = 28
)

func grayscale[T uint8](in gax.FunctionXYC[T]) gax.FunctionXYC[T] {
	grayscale := gax.FuncXYC[T]()
	grayscale.SetXYC(func(x, y, ch gax.Var) T {
		if int(ch) == 3 {
			return 255
		}

		// https://github.com/octu0/blurry/blob/f8ab3409480850181ae06f66db2d01e9db21340d/blurry.cpp#L1695-L1704
		r := int(in.XYC(x, y, gax.Var(0))) * grayR
		g := int(in.XYC(x, y, gax.Var(1))) * grayG
		b := int(in.XYC(x, y, gax.Var(2))) * grayB
		return T((r + g + b) >> 8)
	})

	return grayscale
}

func main() {
	img, err := readImage("./src.png")
	if err != nil {
		panic(err)
	}

	in := gax.ImageRGBA[uint8](img)

	gray := grayscale(in)

	out := gray.Realize(img.Rect.Dx(), img.Rect.Dy(), 4) // 4 = r,g,b,a

	outImg := &image.RGBA{
		Pix:    out.Data(),
		Stride: 4 * img.Rect.Dx(),
		Rect:   image.Rect(0, 0, img.Rect.Dx(), img.Rect.Dy()),
	}

	path, err := saveImage(outImg)
	if err != nil {
		panic(err)
	}
	println("gray", path)
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
