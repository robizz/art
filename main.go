package main

import (
	"fmt"
	"math"
	"os"

	"github.com/robizz/art/matrix"
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
	return fmt.Sprintf(`<polygon points="%s %s %s %s" style="fill:%s;stroke:%s;stroke-width:%d; opacity:0.5" />`,
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

	a := point{300, 800}

	// cube, 6 faces
	// one (front)
	one := rectangle{
		a:           a,
		b:           point{a.x + 100, a.y},
		c:           point{a.x + 100, a.y + 100},
		d:           point{a.x, a.y + 100},
		fill:        "#c41515",
		stroke:      "#c41515",
		strokeWidth: 1,
	}
	rax, ray := matrix.Toisometric(matrix.Vector{float64(one.a.x), float64(one.a.y), 0})
	one.a = point{int(rax), int(ray)}
	rbx, rby := matrix.Toisometric(matrix.Vector{float64(one.b.x), float64(one.b.y), 0})
	one.b = point{int(rbx), int(rby)}
	rcx, rcy := matrix.Toisometric(matrix.Vector{float64(one.c.x), float64(one.c.y), 0})
	one.c = point{int(rcx), int(rcy)}
	rdx, rdy := matrix.Toisometric(matrix.Vector{float64(one.d.x), float64(one.d.y), 0})
	one.d = point{int(rdx), int(rdy)}

	// two (down) plus z
	two := rectangle{
		a:           point{a.x, a.y + 100},
		b:           point{a.x + 100, a.y + 100},
		c:           point{a.x, a.y + 100},       // z
		d:           point{a.x + 100, a.y + 100}, // z
		fill:        "#229d68",
		stroke:      "#229d68",
		strokeWidth: 1,
	}
	rax, ray = matrix.Toisometric(matrix.Vector{float64(two.a.x), float64(two.a.y), 0})
	two.a = point{int(rax), int(ray)}
	rbx, rby = matrix.Toisometric(matrix.Vector{float64(two.b.x), float64(two.b.y), 0})
	two.b = point{int(rbx), int(rby)}
	rcx, rcy = matrix.Toisometric(matrix.Vector{float64(two.c.x), float64(two.c.y), 100})
	two.c = point{int(rcx), int(rcy)}
	rdx, rdy = matrix.Toisometric(matrix.Vector{float64(two.d.x), float64(two.d.y), 100})
	two.d = point{int(rdx), int(rdy)}

	// three (up) plus z
	three := rectangle{
		a:           a,
		b:           point{a.x + 100, a.y},
		c:           a,                     // z
		d:           point{a.x, a.y + 100}, // z
		fill:        "#9d229b",
		stroke:      "#9d229b",
		strokeWidth: 1,
	}
	rax, ray = matrix.Toisometric(matrix.Vector{float64(three.a.x), float64(three.a.y), 0})
	three.a = point{int(rax), int(ray)}
	rbx, rby = matrix.Toisometric(matrix.Vector{float64(three.b.x), float64(three.b.y), 0})
	three.b = point{int(rbx), int(rby)}
	rcx, rcy = matrix.Toisometric(matrix.Vector{float64(three.c.x), float64(three.c.y), 100})
	three.c = point{int(rcx), int(rcy)}
	rdx, rdy = matrix.Toisometric(matrix.Vector{float64(three.d.x), float64(three.d.y), 100})
	three.d = point{int(rdx), int(rdy)}

	//six is one (front plus z)
	six := rectangle{
		a:           a,
		b:           point{a.x + 100, a.y},
		c:           point{a.x + 100, a.y + 100},
		d:           point{a.x, a.y + 100},
		fill:        "#c1c415",
		stroke:      "#c1c415",
		strokeWidth: 1,
	}
	rax, ray = matrix.Toisometric(matrix.Vector{float64(six.a.x), float64(six.a.y), 100})
	six.a = point{int(rax), int(ray)}
	rbx, rby = matrix.Toisometric(matrix.Vector{float64(six.b.x), float64(six.b.y), 100})
	six.b = point{int(rbx), int(rby)}
	rcx, rcy = matrix.Toisometric(matrix.Vector{float64(six.c.x), float64(six.c.y), 100})
	six.c = point{int(rcx), int(rcy)}
	rdx, rdy = matrix.Toisometric(matrix.Vector{float64(six.d.x), float64(six.d.y), 100})
	six.d = point{int(rdx), int(rdy)}

	rects = append(rects, one, two, three, six)

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
