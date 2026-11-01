// package matrix should help projecting a solid in 2d space using isometric projection
// https://en.wikipedia.org/wiki/Isometric_projection
package matrix

import "math"

type Matrix [3][3]float64

type Vector [3]float64

// cij = ai1 · b1j + ai2 · b2j + ai3 ·  b3j
func Multiply(a, b Matrix) Matrix {
	var c Matrix
	for k := 0; k < 3; k++ {
		for n := 0; n < 3; n++ {
			for i := 0; i < 3; i++ {
				c[k][n] = c[k][n] + a[k][i]*b[i][n]
			}
		}
	}
	return c
}

func Transform(a Matrix, v Vector) Vector {
	var x Vector
	for k := 0; k < 3; k++ {
		for i := 0; i < 3; i++ {
			x[k] = x[k] + a[k][i]*v[i]
		}
	}
	return x
}

func SMMultiply(s float64, a Matrix) Matrix {
	var c Matrix
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			c[i][j] = a[i][j] * s
		}
	}
	return c
}

func SVMultiply(s float64, v Vector) Vector {
	var c Vector
	for i := 0; i < 3; i++ {
		c[i] = v[i] * s
	}
	return c
}

// https://en.wikipedia.org/wiki/Isometric_projection#Mathematics
func Toisometric(v Vector) (x, y float64) {

	// rotation around the vertical (here y) axis by β = 45 deg
	beta := 45.0
	ry := Matrix{
		{1, 0, 0},
		{0, math.Cos(beta), -math.Sin(beta)},
		{0, math.Sin(beta), math.Cos(beta)},
	}
	// followed by a rotation around the horizontal (here x) axis by α  α = arcsin(tan 30°) ≈ 35.264°
	rx := Matrix{
		{math.Cos(beta), 0, -math.Sin(beta)},
		{0, 1, 0},
		{math.Sin(beta), 0, math.Cos(beta)},
	}

	// then followed by an orthographic projection to the xy-plane:
	p := Matrix{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 0},
	}

	var tv Vector

	tv = Transform(p, Transform(Multiply(ry, rx), v))
	return tv[0], tv[1]
}
