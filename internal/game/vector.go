package game

import "math"

type Vector struct {
	X, Y float64
}

func (v *Vector) Add(other *Vector) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vector) Sub(other *Vector) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vector) Mul(scalar float64) {
	v.X *= scalar
	v.Y *= scalar
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector) Normalize() {
	length := v.Length()
	if length > 0 {
		v.X /= length
		v.Y /= length
	}
}
