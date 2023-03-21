package gax

import (
	"fmt"
	"image"
)

type imageRGBAFunc[T Number] struct {
	FunctionXYC[T]

	img *image.RGBA
	fn  funcXYC[T]
}

func (r *imageRGBAFunc[T]) XYC(x, y, ch Var) T {
	if r.fn != nil {
		return r.fn(x, y, ch)
	}

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

func (r *imageRGBAFunc[T]) SetXYC(fn funcXYC[T]) {
	r.fn = fn
}

func ImageRGBA[T Number](img *image.RGBA) FunctionXYC[T] {
	return &imageRGBAFunc[T]{
		img: img,
	}
}
