package pointer

func String(s string) *string {
	return &s
}

func Float32(f float32) *float32 {
	return &f
}

func Ptr[T any](v T) *T {
	return &v
}
