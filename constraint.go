package gax

type Signed interface {
	~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Float interface {
	~float32 | ~float64
}

type Number interface {
	Signed | Unsigned | Float
}
