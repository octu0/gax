package gax

import (
	"fmt"
	"image"
)

type imageRGBAReader[T Number] struct {
	img *image.RGBA
}

func (r *imageRGBAReader[T]) XY(x, y Var) T {
	pos := r.img.PixOffset(int(x), int(y))
	return T(r.img.Pix[pos])
}

func (r *imageRGBAReader[T]) XYC(x, y, ch Var) T {
	c := r.img.RGBAAt(int(x), int(y))
	switch ch {
	case 0:
		return T(c.R)
	case 1:
		return T(c.G)
	case 2:
		return T(c.B)
	case 3:
		return T(c.A)
	}
	panic(fmt.Sprintf("out of range ch:%d", int(ch)))
}

func ImageRGBA[T Number](img *image.RGBA) *Buffer[T] {
	return &Buffer[T]{
		in: &imageRGBAReader[T]{img},
	}
}
