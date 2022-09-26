package util

func Filter[T any](arr []T, f func(T) bool) []T {
	j := 0
	for _, v := range arr {
		if f(v) {
			arr[j] = v
			j++
		}
	}
	return arr[:j]
}
