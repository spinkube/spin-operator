package generics

func MapList[A ~[]X, X, Y any](input A, mapper func(X) Y) []Y {
	result := make([]Y, len(input))
	for idx, value := range input {
		result[idx] = mapper(value)
	}
	return result
}

func MapMap[X ~map[A]B, Y ~map[A]C, A comparable, B, C any](input X, mapper func(A, B) C) Y {
	result := make(Y, len(input))
	for key, value := range input {
		result[key] = mapper(key, value)
	}
	return result
}

func AssociateBy[A ~[]X, X any, Y comparable](input A, assocBy func(X) Y) map[Y]X {
	result := make(map[Y]X, len(input))
	for _, elem := range input {
		result[assocBy(elem)] = elem
	}
	return result
}

func Ptr[T any](v T) *T {
	return &v
}
