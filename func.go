package gax

type (
	Var int
)

type (
	Output[T Number] interface {
		Data() []T
	}
)

type todoOutput[T Number] struct {
	data []T
}

func (o *todoOutput[T]) Data() []T {
	return o.data
}

type (
	funcXY[T Number] func(x, y Var) T

	ReaderXY[T Number] interface {
		XY(x, y Var) T
	}

	FunctionXY[T Number] interface {
		ReaderXY[T]
		SetXY(funcXY[T])
		Realize(width, height int) Output[T]
	}
)

type (
	funcXYC[T Number] func(x, y, ch Var) T

	ReaderXYC[T Number] interface {
		XYC(x, y, ch Var) T
	}

	FunctionXYC[T Number] interface {
		ReaderXYC[T]
		SetXYC(funcXYC[T])
		Realize(width, height, ch int) Output[T]
	}
)

type xyFunc[T Number] struct {
	FunctionXY[T]

	inout funcXY[T]
}

func (f *xyFunc[T]) XY(x, y Var) T {
	return f.inout(x, y)
}

func (f *xyFunc[T]) SetXY(fn funcXY[T]) {
	f.inout = fn
}

func (f *xyFunc[T]) Realize(width, height int) Output[T] {
	out := make([]T, 0, width*height)
	for h := 0; h < height; h += 1 {
		for w := 0; w < width; w += 1 {
			out = append(out, f.inout(Var(w), Var(h)))
		}
	}
	return &todoOutput[T]{out}
}

func FuncXY[T Number]() FunctionXY[T] {
	return new(xyFunc[T])
}

type xycFunc[T Number] struct {
	FunctionXYC[T]

	inout funcXYC[T]
}

func (f *xycFunc[T]) XYC(x, y, ch Var) T {
	return f.inout(x, y, ch)
}

func (f *xycFunc[T]) SetXYC(fn funcXYC[T]) {
	f.inout = fn
}

func (f *xycFunc[T]) Realize(width, height, ch int) Output[T] {
	// data layout...
	//stride := width * ch

	out := make([]T, 0, width*height*ch)
	for h := 0; h < height; h += 1 {
		for w := 0; w < width; w += 1 {
			for c := 0; c < ch; c += 1 {
				out = append(out, f.inout(Var(w), Var(h), Var(c)))
			}
		}
	}
	return &todoOutput[T]{out}
}

func FuncXYC[T Number]() FunctionXYC[T] {
	return new(xycFunc[T])
}
