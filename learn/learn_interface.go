package main

import "fmt"

type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func CalculateArea(s Shape) float64 {
	return s.Area()
}

func main() {
	r := Rectangle{Width: 5, Height: 4}
	c := Circle{Radius: 3}

	fmt.Println("矩形面积:", CalculateArea(r))
	fmt.Println("圆形面积:", CalculateArea(c))
}
