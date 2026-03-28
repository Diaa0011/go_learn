package main

// import (
// 	"math"
// 	"testing"
// )

// type Rectangle struct {
// 	Width  float64
// 	Height float64
// }

// func (r Rectangle) Area() float64 {
// 	return r.Width * r.Height
// }

// type Circle struct {
// 	Radius float64
// }

// func (c Circle) Area() float64 {
// 	return math.Pi * math.Pow(c.Radius, 2)
// }

// func TestArea(t *testing.T) {

// 	t.Run("rectangles", func(t *testing.T) {
// 		rectangle := Rectangle{12, 6}
// 		got := rectangle.Area()
// 		want := 72.0

// 		if got != want {
// 			t.Errorf("got %g want %g", got, want)
// 		}
// 	})

// 	t.Run("circles", func(t *testing.T) {
// 		circle := Circle{10}
// 		got := circle.Area()
// 		want := 314.1592653589793

// 		if got != want {
// 			t.Errorf("got %g want %g", got, want)
// 		}
// 	})

// }
