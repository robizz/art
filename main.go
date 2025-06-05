package main

import (
	"fmt"
	"math"
	"os"
)

type svg interface {
	print() string
}

type canvas struct {
	template string
	elements []svg
}

func (c *canvas) print() string {

	c.template = `<?xml version="1.0" encoding="UTF-8" ?>
<svg version="1.1" width="900" height="900" xmlns="http://www.w3.org/2000/svg">
<rect width="900" height="900" x="0" y="0" fill="#ffffff" />
%s
</svg>`

	var content string
	for _, e := range c.elements {
		content = content + e.print() + "\n"
	}
	return fmt.Sprintf(c.template, content)
}

type point struct {
	x int
	y int
}

func (p point) print() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

// <polygon points="100,10 150,190 50,190" style="fill:lime;stroke:purple;stroke-width:3" />
type triangle struct {
	a           point
	b           point
	c           point
	fill        string
	stroke      string
	strokeWidth int
}

func (t triangle) print() string {
	return fmt.Sprintf(`<polygon points="%s %s  %s" style="fill:%s;stroke:%s;stroke-width:%d" />`,
		t.a.print(), t.b.print(), t.c.print(), t.fill, t.stroke, t.strokeWidth)
}

func triangleFrom(center point, side int) triangle {
	// https://math.stackexchange.com/a/1344707
	return triangle{
		a: point{center.x, int(float64(center.y) + ((math.Sqrt(3) / 3) * float64(side)))},
		b: point{
			x: int(float64(center.x) - (float64(side) / 2)),
			y: int(float64(center.y) - ((math.Sqrt(3) / 6) * float64(side))),
		},
		c: point{
			x: int(float64(center.x) + (float64(side) / 2)),
			y: int(float64(center.y) - ((math.Sqrt(3) / 6) * float64(side))),
		},
	}
}

type polyFour struct {
	a           point
	b           point
	c           point
	d           point
	fill        string
	stroke      string
	strokeWidth int
}

func (p polyFour) print() string {
	return fmt.Sprintf(`<polygon points="%s %s %s %s" style="fill:%s;stroke:%s;stroke-width:%d" />`,
		p.a.print(), p.b.print(), p.c.print(), p.d.print(), p.fill, p.stroke, p.strokeWidth)
}

type isometricCube struct {
	top   polyFour
	left  polyFour
	right polyFour
}

func (i isometricCube) print() string {
	return fmt.Sprintf("\n%s\n%s\n%s\n", i.top.print(), i.left.print(), i.right.print())
}

func rotate(a, b point, theta float64) (point, point) {
	// https://math.stackexchange.com/a/4287500
	//cx=cos(θ) (bx−ax)−sin(θ) (by−ay)+ax
	//cy=sin(θ) (bx−ax)+cos(θ) (by−ay)+ay

	cx := int(math.Cos(theta)*float64(b.x-a.x) - math.Sin(theta)*float64(b.y-a.y) + float64(a.x))
	cy := int(math.Sin(theta)*float64(b.x-a.x) + math.Cos(theta)*float64(b.y-a.y) + float64(a.y))
	return a, point{cx, cy}
}

func isometricCubeFrom(origin point, side int) isometricCube {

	right := polyFour{
		a: point{origin.x, origin.y},
		b: point{origin.x, origin.y + side},
		c: point{},
		d: point{},
	}
	_, rd := rotate(right.a, point{right.a.x + side, right.a.y}, 30)
	right.d = rd
	_, rc := rotate(right.b, point{right.b.x + side, right.b.y}, 30)
	right.c = rc

	left := polyFour{
		a: point{},
		b: point{},
		c: right.b,
		d: right.a,
	}

	top := polyFour{
		a: right.b,
		b: left.b,
		c: point{},
		d: right.c,
	}

	return isometricCube{
		right: right,
		left:  left,
		top:   top,
	}
}

func isCircumference(p, center point, radius, tolerance float64) bool {
	// https://www.quora.com/What-is-the-Cartesian-equation-of-a-circle/answer/Abhay-Roy-51
	t := math.Pow(float64(p.x-center.x), 2) + math.Pow(float64(p.y-center.y), 2)
	rr := float64(math.Pow(radius, 2))
	// here we are comparing againsta the circle formula,
	// but we are in a discrete grid so we need to add some tolerance otherwise the points that
	// will be detected in the circumference will be very low in number
	return t > rr-tolerance && t < rr+tolerance
}

func isEllipse(p, center point, a, b, tolerance float64) bool {
	// a(x-h)^2 + b(y-k)^2 = 1
	t := (math.Pow(float64(p.x-center.x), 2) / math.Pow(b, 2)) + (math.Pow(float64(p.y-center.y), 2) / math.Pow(a, 2))
	rr := 1.0
	// here we are comparing againsta the circle formula,
	// but we are in a discrete grid so we need to add some tolerance otherwise the points that
	// will be detected in the circumference will be very low in number
	return t > rr-tolerance && t < rr+tolerance
}

func main() {

	// define a "ghost circle"
	// x^2+y^2 = 38
	// for all the points in the grid determine if you are in the circumference or not
	nablas := []svg{}
	// // using brute forcing
	// for d := 80.0; d < 300.0; d += 1.0 {
	// 	for x := 0; x < 900; x++ {
	// 		for y := 0; y < 900; y++ {
	// 			density := rand.Intn(int(math.Pow(d/30, 2))) == 1
	// 			if isEllipse(point{x, y}, point{450, 450}, 1.0+(d/2), 90.0+d, 0.01) && density {
	// 				// density := rand.Intn(int(math.Pow(d/40, 2))) == 1
	// 				// if isCircumference(point{x, y}, point{450, 450}, d, 90) && density {
	// 				// t := triangleFrom(point{x, y}, 5)
	// 				// rand.IntN(max+1-min) + min
	// 				randx := rand.Intn(int(math.Pow(d/80, 3))+1+int(math.Pow(d/80, 3))) - int(math.Pow(d/80, 3))
	// 				randy := rand.Intn(int(math.Pow(d/80, 3))+1+int(math.Pow(d/80, 3))) - int(math.Pow(d/80, 3))
	// 				// if we are approaching outer ellipse so d value is in the last 50 points
	// 				// and we are in the upper half of the image, let's remove some more triengles
	// 				clean := rand.Intn(int(math.Pow(d/30, 2))) != 1
	// 				if d > 80 && y < 405 && clean {
	// 					continue
	// 				}
	// 				t := triangleFrom(point{x + randx, y + randy}, 5)
	// 				t.fill = "#ffffff"
	// 				t.stroke = "#000000"
	// 				t.strokeWidth = 1
	// 				nablas = append(nablas, t)
	// 			}

	// 		}
	// 	}
	// }

	nablas = append(nablas, isometricCubeFrom(point{100, 100}, 50).right)
	// transform it in svg triangle
	// inject in a svg file template
	c := canvas{
		elements: nablas,
	}

	// write file
	f, err := os.Create("art.svg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	svgfile := c.print()
	f.WriteString(svgfile)
}
