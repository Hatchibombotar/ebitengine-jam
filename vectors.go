package main

import "math"

type Vector struct {
	X, Y float64
}

func VectorAdd(v1 Vector, v2 Vector) Vector {
	return Vector{
		v1.X + v2.X,
		v1.Y + v2.Y,
	}
}
func VectorSubtract(v1, v2 Vector) Vector {
	return Vector{
		v1.X - v2.X,
		v1.Y - v2.Y,
	}
}
func VectorIs(v1, v2 Vector) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}
func VectorFloor(v1 Vector) Vector {
	return Vector{
		math.Floor(v1.X),
		math.Floor(v1.Y),
	}
}

func VectorScale(v1 Vector, s float64) Vector {
	return Vector{
		v1.X * s,
		v1.Y * s,
	}
}

func LerpVectors(v1, v2 Vector, p float64) Vector {
	return Vector{
		X: v1.X + ((v2.X - v1.X) * p),
		Y: v1.Y + ((v2.Y - v1.Y) * p),
	}
}

func VectorMagnitude(v1 Vector) float64 {
	return math.Sqrt(math.Pow(v1.X, 2) + math.Pow(v1.Y, 2))
}
func VectorNormalise(v1 Vector) Vector {
	magnitude := VectorMagnitude(v1)
	return Vector{
		X: v1.X / magnitude,
		Y: v1.Y / magnitude,
	}
}

func VectorAngle(v Vector) float64 {
	return math.Atan2(v.Y, v.X)
}
