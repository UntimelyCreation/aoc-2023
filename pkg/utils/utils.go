package utils

func Mod(a, b int) int {
	return (a%b + b) % b
}

func Gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func Lcm(a, b int) int {
	return int(a * b / Gcd(a, b))
}

func Last[T any](slice []T) T {
	return slice[len(slice)-1]
}

func Reverse[T any](slice []T) []T {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
