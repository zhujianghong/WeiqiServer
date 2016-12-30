package main

import (
	"fmt"
	"math"
)

var (
	Back3D    Vector3D = Vector3D{0, 0, -1}
	Down3D    Vector3D = Vector3D{0, -1, 0}
	Forward3D Vector3D = Vector3D{0, 0, 1}
	Left3D    Vector3D = Vector3D{-1, 0, 0}
	One3D     Vector3D = Vector3D{1, 1, 1}
	Right3D   Vector3D = Vector3D{1, 0, 0}
	Up3D      Vector3D = Vector3D{0, 1, 0}
	zero3D    Vector3D = Vector3D{0, 0, 0}
)

type Vector3D struct {
	X, Y, Z float32
}

func (v3D *Vector3D) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v3D.X, v3D.Y, v3D.Z)
}

func (v3D *Vector3D) index3D(idx int) float32 {
	if idx == 0 {
		return v3D.X
	}
	if idx == 1 {
		return v3D.Y
	}
	if idx == 2 {
		return v3D.Z
	}
	panic("index out of range")
}

func (v3D *Vector3D) indexSet3D(idx int, value float32) {
	if idx == 0 {
		v3D.X = value
		return
	}
	if idx == 1 {
		v3D.Y = value
		return
	}
	if idx == 2 {
		v3D.Z = value
		return
	}
	panic("index out of range")
}

func (v3D *Vector3D) modulus3D() float32 {
	return float32(math.Sqrt(float64(v3D.modulusSquare3D())))
}

func (v3D *Vector3D) modulusSquare3D() float32 {
	return v3D.X*v3D.X + v3D.Y*v3D.Y
}

func GetVector3D(values ...float32) Vector3D {
	var result Vector3D
	size := len(values)
	if size >= 1 {
		result.X = values[0]
	}
	if size >= 2 {
		result.Y = values[1]
	}
	if size >= 3 {
		result.Z = values[2]
	}
	return result
}

func add3D(v3D1, v3D2 Vector3D) Vector3D {
	var reuslt Vector3D
	reuslt.X = v3D1.X + v3D2.X
	reuslt.Y = v3D1.Y + v3D2.Y
	reuslt.Z = v3D1.Z + v3D2.Z
	return reuslt
}

func sub3D(v3D1, v3D2 Vector3D) Vector3D {
	var reuslt Vector3D
	reuslt.X = v3D1.X - v3D2.X
	reuslt.Y = v3D1.Y - v3D2.Y
	reuslt.Z = v3D1.Z - v3D2.Z
	return reuslt
}

func scale3D(v3D Vector3D, factor float32) Vector3D {
	v3D.X *= factor
	v3D.Y *= factor
	v3D.Z *= factor
	return v3D
}

func vector3DNeg(v3D Vector3D) Vector3D {
	v3D.X = -v3D.X
	v3D.Y = -v3D.Y
	v3D.Z = -v3D.Z
	return v3D
}

func vector3DEqual(v3D1, v3D2 Vector3D) bool {
	reuslt := sub3D(v3D1, v3D2)
	return reuslt.modulusSquare3D() < 9.99999944E-11
}

func vector3DnotEqual(v3D1, v3D2 Vector3D) bool {
	return !vector3DEqual(v3D1, v3D2)
}

func cross3D(v3D1, v3D2 Vector3D) Vector3D {
	return Vector3D{
		v3D1.Y*v3D2.Z - v3D1.Z*v3D2.Y,
		v3D1.Z*v3D2.X - v3D1.X*v3D2.Z,
		v3D1.X*v3D2.Y - v3D1.Y*v3D2.X,
	}
}

func dotProduct3D(v3D1, v3D2 Vector3D) float32 {
	return v3D1.X*v3D2.X + v3D1.Y*v3D2.Y + v3D1.Z*v3D2.Z
}

func crossProduct3D(v3D1, v3D2 Vector3D) Vector3D {
	return Vector3D{
		v3D1.X * v3D2.X,
		v3D1.Y * v3D2.Y,
		v3D1.Z * v3D2.Z,
	}
}

func distance3D(v3D1, v3D2 Vector3D) float32 {
	result := Vector3D{v3D1.X - v3D2.X, v3D1.Y - v3D2.Y, v3D1.Z - v3D2.Z}
	return result.modulus3D()
}

func maxVector3D(v3D1, v3D2 Vector3D) Vector3D {
	result := v3D1
	if result.X < v3D2.X {
		result.X = v3D2.X
	}
	if result.Y < v3D2.Y {
		result.Y = v3D2.Y
	}
	if result.Z < v3D2.Z {
		result.Z = v3D2.Z
	}
	return result
}

func minVector3D(v3D1, v3D2 Vector3D) Vector3D {
	result := v3D1
	if result.X > v3D2.X {
		result.X = v3D2.X
	}
	if result.Y > v3D2.Y {
		result.Y = v3D2.Y
	}
	if result.Z > v3D2.Z {
		result.Z = v3D2.Z
	}
	return result
}
