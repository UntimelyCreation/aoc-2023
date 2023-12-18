package grid

import "fmt"

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

func TurnLeft(dir Direction) Direction {
	switch dir {
	case Up:
		return Left
	case Right:
		return Up
	case Down:
		return Right
	case Left:
		return Down
	default:
		panic("Not valid direction")
	}
}

func TurnRight(dir Direction) Direction {
	switch dir {
	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up
	default:
		panic("Not valid direction")
	}
}

func TurnBack(dir Direction) Direction {
	switch dir {
	case Up:
		return Down
	case Right:
		return Left
	case Down:
		return Up
	case Left:
		return Right
	default:
		panic("Not valid direction")
	}
}

type Position struct {
	Row int
	Col int
}

func (pos Position) Move(dir Direction) Position {
	x, y := 0, 0
	switch dir {
	case Up:
		x, y = -1, 0
	case Right:
		x, y = 0, 1
	case Down:
		x, y = 1, 0
	case Left:
		x, y = 0, -1
	default:
		panic("Not  direction")
	}
	return Position{Row: pos.Row + x, Col: pos.Col + y}
}

type Grid[T any] map[Position]T

func (g Grid[T]) XRange() (int, int) {
	xMin, xMax := 0, 0
	for pos := range g {
		xMin, xMax = min(xMin, pos.Row), max(xMax, pos.Row)
	}
	return xMin, xMax
}

func (g Grid[T]) YRange() (int, int) {
	yMin, yMax := 0, 0
	for pos := range g {
		yMin, yMax = min(yMin, pos.Col), max(yMax, pos.Col)
	}
	return yMin, yMax
}

func (g Grid[T]) Dimensions() (int, int) {
	rows, cols := 0, 0
	for pos := range g {
		rows, cols = max(rows, pos.Row), max(cols, pos.Col)
	}
	return rows + 1, cols + 1
}

func (g Grid[T]) Transpose() Grid[T] {
	transposedG := make(Grid[T])
	for k, v := range g {
		transposedG[Position{Row: k.Col, Col: k.Row}] = v
	}
	return transposedG
}

func (g Grid[T]) Rotate(r complex128) Grid[T] {
	rotatedG := make(Grid[T])

	for k, v := range g {
		z := complex(float64(k.Row), float64(k.Col))
		z *= r
		rotatedG[Position{Row: int(real(z)), Col: int(imag(z))}] = v
	}
	return rotatedG
}

func (g Grid[T]) PrettyPrint() {
	xMin, xMax := g.XRange()
	yMin, yMax := g.YRange()

	for i := xMin; i <= xMax; i++ {
		for j := yMin; j <= yMax; j++ {
			fmt.Print(g[Position{Row: i, Col: j}], " ")
		}
		fmt.Print("\n")
	}
}
