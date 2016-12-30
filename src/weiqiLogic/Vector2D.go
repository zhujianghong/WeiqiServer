package main

import (
	"fmt"
	"math"
)

var (
	Down2D  Vector2D = Vector2D{0, -1}
	Left2D  Vector2D = Vector2D{-1, 0}
	One2D   Vector2D = Vector2D{1, 1}
	Right2D Vector2D = Vector2D{1, 0}
	Up2D    Vector2D = Vector2D{0, 1}
	zero2D  Vector2D = Vector2D{0, 0}
)

type Vector2D struct {
	X, Y float32
}

func (v2D *Vector2D) index2D(idx int) float32 {
	if idx == 0 {
		return v2D.X
	}
	if idx == 1 {
		return v2D.Y
	}
	panic("index out of range")
}

func (v2D *Vector2D) indexSet2D(idx int, value float32) {
	if idx == 0 {
		v2D.X = value
		return
	}
	if idx == 1 {
		v2D.Y = value
		return
	}
	panic("index out of range")
}

func (v2D *Vector2D) modulus2D() float32 {
	return float32(math.Sqrt(float64(v2D.modulusSquare2D())))
}

func (v2D *Vector2D) modulusSquare2D() float32 {
	return v2D.X*v2D.X + v2D.Y*v2D.Y
}

func (v2D *Vector2D) normalize() {
	scale := v2D.modulus2D()
	if scale == 0 {
		return
	}
	v2D.X /= scale
	v2D.Y /= scale
	return
}

func (v2D *Vector2D) normalized() Vector2D {
	result := *v2D
	result.normalize()
	return result
}

func (v2D *Vector2D) String() string {
	return fmt.Sprintf("(%f, %f)", v2D.X, v2D.Y)
}

func add2D(v2D1, v2D2 Vector2D) Vector2D {
	var reuslt Vector2D
	reuslt.X = v2D1.X + v2D2.X
	reuslt.Y = v2D1.Y + v2D2.Y
	return reuslt
}

func sub2D(v2D1, v2D2 Vector2D) Vector2D {
	var reuslt Vector2D
	reuslt.X = v2D1.X - v2D2.X
	reuslt.Y = v2D1.Y - v2D2.Y
	return reuslt
}

func scale2D(v2D Vector2D, factor float32) Vector2D {
	v2D.X *= factor
	v2D.Y *= factor
	return v2D
}

func vector2DNeg(v2D Vector2D) Vector2D {
	v2D.X = -v2D.X
	v2D.Y = -v2D.Y
	return v2D
}

func vector2DEqual(v2D1, v2D2 Vector2D) bool {
	reuslt := sub2D(v2D1, v2D2)
	return reuslt.modulusSquare2D() < 9.99999944E-11
}

func vector2DnotEqual(v2D1, v2D2 Vector2D) bool {
	return !vector2DEqual(v2D1, v2D2)
}

func clamp(v2D Vector2D, maxLength float32) Vector2D {
	if v2D.modulusSquare2D() > maxLength*maxLength {
		result := v2D.normalized()
		result.X *= maxLength
		result.Y *= maxLength
		return result
	}
	return v2D
}

func distance2D(v2D1, v2D2 Vector2D) float32 {
	result := Vector2D{v2D1.X - v2D2.X, v2D1.Y - v2D2.Y}
	return result.modulus2D()
}

func dotProduct2D(v2D1, v2D2 Vector2D) float32 {
	return v2D1.X*v2D2.X + v2D1.Y*v2D2.Y
}

func crossProduct2D(v2D1, v2D2 Vector2D) Vector2D {
	return Vector2D{
		v2D1.X * v2D2.X,
		v2D1.Y * v2D2.Y,
	}
}

func lerpUnclamped2D(v2D1, v2D2 Vector2D, t float32) Vector2D {
	return Vector2D{
		v2D1.X + (v2D2.X-v2D1.X)*t,
		v2D1.Y + (v2D2.Y-v2D1.Y)*t,
	}
}

func maxVector2D(v2D1, v2D2 Vector2D) Vector2D {
	result := v2D1
	if result.X < v2D2.X {
		result.X = v2D2.X
	}
	if result.Y < v2D2.Y {
		result.Y = v2D2.Y
	}
	return result
}

func minVector2D(v2D1, v2D2 Vector2D) Vector2D {
	result := v2D1
	if result.X > v2D2.X {
		result.X = v2D2.X
	}
	if result.Y > v2D2.Y {
		result.Y = v2D2.Y
	}
	return result
}

func moveToward2D(pos, dst Vector2D, maxDistanceDelta float32) Vector2D {
	result := sub2D(dst, pos)
	mdls := result.modulus2D()
	if mdls < maxDistanceDelta || mdls == 0.0 {
		return result
	}
	return add2D(pos, scale2D(result, mdls/maxDistanceDelta))
}

func reflect2D(v2D1, v2D2 Vector2D) Vector2D {
	return add2D(scale2D(v2D2, -2.0*dotProduct2D(v2D1, v2D2)), v2D1)
}
