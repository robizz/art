package main

import (
	"fmt"
	"math"
	"math/rand"
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
	return fmt.Sprintf(`<polygon points="%s %s %s" style="fill:%s;stroke:%s;stroke-width:%d" />`,
		t.a.print(), t.b.print(), t.c.print(), t.fill, t.stroke, t.strokeWidth)
}

type ellipse struct {
	center      point
	radius      point
	fill        string
	stroke      string
	strokeWidth int
}

func (e ellipse) print() string {
	return fmt.Sprintf(`<ellipse rx="%d" ry="%d" cx="%d" cy="%d" style="fill:%s;stroke:%s;stroke-width:%d" />`,
		e.radius.x, e.radius.y, e.center.x, e.center.y, e.fill, e.stroke, e.strokeWidth)
}

type rectangle struct {
	a           point
	b           point
	c           point
	d           point
	fill        string
	stroke      string
	strokeWidth int
}

func (r rectangle) print() string {
	return fmt.Sprintf(`<polygon points="%s %s %s %s" style="fill:%s;stroke:%s;stroke-width:%d" />`,
		r.a.print(), r.b.print(), r.c.print(), r.d.print(), r.fill, r.stroke, r.strokeWidth)
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
	rects := []svg{}
	// using brute forcing
	for d := 80.0; d < 300.0; d += 1.0 {
		for x := 0; x < 900; x++ {
			for y := 0; y < 900; y++ {
				density := rand.Intn(int(math.Pow(d/30, 2))) == 1
				if isEllipse(point{x, y}, point{450, 450}, 1.0+(d/2), 90.0+d, 0.01) && density {
					// density := rand.Intn(int(math.Pow(d/40, 2))) == 1
					// if isCircumference(point{x, y}, point{450, 450}, d, 90) && density {
					// t := triangleFrom(point{x, y}, 5)
					// rand.IntN(max+1-min) + min
					randx := rand.Intn(int(math.Pow(d/80, 3))+1+int(math.Pow(d/80, 3))) - int(math.Pow(d/80, 3))
					randy := rand.Intn(int(math.Pow(d/80, 3))+1+int(math.Pow(d/80, 3))) - int(math.Pow(d/80, 3))
					// if we are approaching outer ellipse so d value is in the last 50 points
					// and we are in the upper half of the image, let's remove some more triengles
					clean := rand.Intn(int(math.Pow(d/30, 2))) != 1
					if d > 80 && y < 405 && clean {
						continue
					}

					fill := fmt.Sprintf("rgb(%d,%d,%d)", y%200, int(d)%200, int(d)%200)

					a := point{x + randx, y + randy}
					r := rectangle{
						a:           a,
						b:           point{a.x + 2, a.y},
						c:           point{a.x + 2, a.y + 5 + randy},
						d:           point{a.x, a.y + 5 + randy},
						fill:        fill,
						stroke:      "#FFFFFF",
						strokeWidth: 1,
					}
					rects = append(rects, r)
				}

			}
		}
	}

	// transform it in svg triangle
	// inject in a svg file template
	c := canvas{
		elements: rects,
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
