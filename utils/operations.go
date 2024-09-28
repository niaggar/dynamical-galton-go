package utils

type Point = [2]float64

func Add(a, b *Point) *Point {
	return &Point{a[0] + b[0], a[1] + b[1]}
}

func Sub(a, b *Point) *Point {
	return &Point{a[0] - b[0], a[1] - b[1]}
}

func Mul(a *Point, scalar float64) *Point {
	return &Point{a[0] * scalar, a[1] * scalar}
}

func DistanceSquare(a, b *Point) float64 {
	return (a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1])
}
