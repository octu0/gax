package gax

type (
	XYFunc[T Number]  func(x, y Var) T
	XYCFunc[T Number] func(x, y, ch Var) T
)

type Reader[T Number] interface {
	XY(x, y Var) T
	XYC(x, y, ch Var) T
}

type XYReader[T Number] struct {
	xy XYFunc[T]
}

func (r *XYReader[T]) XY(x, y Var) T {
	return r.xy(x, y)
}

func (r *XYReader[T]) XYC(x, y, ch Var) T {
	return r.xy(x, y)
}

type XYCReader[T Number] struct {
	xyc XYCFunc[T]
}

func (r *XYCReader[T]) XY(x, y Var) T {
	return r.xyc(x, y, 0)
}

func (r *XYCReader[T]) XYC(x, y, ch Var) T {
	return r.xyc(x, y, ch)
}

type Buffer[T Number] struct {
	in Reader[T]
}

func (b *Buffer[T]) XY(x, y Var) T {
	return b.in.XY(x, y)
}

func (b *Buffer[T]) XYC(x, y, ch Var) T {
	return b.in.XYC(x, y, ch)
}

func (b *Buffer[T]) SetXY(fn XYFunc[T]) {
	b.in = &XYReader[T]{fn}
}

func (b *Buffer[T]) SetXYC(fn XYCFunc[T]) {
	b.in = &XYCReader[T]{fn}
}

func (b *Buffer[T]) Realize(width, height, ch int) []T {
	// data layout...
	//stride := width * ch

	out := make([]T, 0, width*height*ch)
	for h := 0; h < height; h += 1 {
		for w := 0; w < width; w += 1 {
			for c := 0; c < ch; c += 1 {
				out = append(out, b.in.XYC(Var(w), Var(h), Var(c)))
			}
		}
	}
	return out
}
