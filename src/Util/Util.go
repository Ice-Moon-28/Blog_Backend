package Util

func ArrayHasValue[T int | string | float32 | float64 | struct{}](value T, array []T) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
func GetValueIndexInArray[T int | string | float32 | float64](value T, array []T) int {
	for i, v := range array {
		if v == value {
			return i
		}
	}
	return -1
}
