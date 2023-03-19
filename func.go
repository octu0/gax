package gax

func Func[T Number]() *Buffer[T] {
	return &Buffer[T]{}
}
