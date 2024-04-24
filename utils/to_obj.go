package utils

func ToObjArray[T any](pointers []*T) []T {
	var result []T
	for _, p := range pointers {
		result = append(result, *p)
	}
	return result
}
