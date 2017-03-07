// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Vector2 struct {
	X float32
	Y float32
}

func NewVector2(x, y float32) *Vector2 {

	return &Vector2{X: x, Y: y}
}

func (this *Vector2) Set(x, y float32) *Vector2 {

	this.X = x
	this.Y = y
	return this
}

func (this *Vector2) SetX(x float32) *Vector2 {

	this.X = x
	return this
}

func (this *Vector2) SetY(y float32) *Vector2 {

	this.Y = y
	return this
}

func (this *Vector2) SetComponent(index int, value float32) {

	switch index {
	case 0:
		this.X = value
	case 1:
		this.Y = value
	default:
		panic("index is out of range")
	}
}

func (this *Vector2) GetComponent(index int) float32 {

	switch index {
	case 0:
		return this.X
	case 1:
		return this.Y
	default:
		panic("index is out of range: ")
	}
}

func (this *Vector2) Copy(v *Vector2) *Vector2 {

	this.X = v.X
	this.Y = v.Y
	return this
}

func (this *Vector2) Add(v *Vector2) *Vector2 {

	this.X += v.X
	this.Y += v.Y
	return this
}

func (this *Vector2) AddScalar(s float32) *Vector2 {

	this.X += s
	this.Y += s
	return this
}

func (this *Vector2) AddVectors(a, b *Vector2) *Vector2 {

	this.X = a.X + b.X
	this.Y = a.Y + b.Y
	return this
}

func (this *Vector2) Sub(v *Vector2) *Vector2 {

	this.X -= v.X
	this.Y -= v.Y
	return this
}

func (this *Vector2) SubScalar(s float32) *Vector2 {

	this.X -= s
	this.Y -= s
	return this
}

func (this *Vector2) SubVectors(a, b *Vector2) *Vector2 {

	this.X = a.X - b.X
	this.Y = a.Y - b.Y
	return this

}

func (this *Vector2) Multiply(v *Vector2) *Vector2 {

	this.X *= v.X
	this.Y *= v.Y
	return this
}

func (this *Vector2) MultiplyScalar(s float32) *Vector2 {

	this.X *= s
	this.Y *= s
	return this
}

func (this *Vector2) Divide(v *Vector2) *Vector2 {

	this.X /= v.X
	this.Y /= v.Y
	return this
}

func (this *Vector2) DivideScalar(scalar float32) *Vector2 {

	if scalar != 0 {
		invScalar := 1 / scalar
		this.X *= invScalar
		this.Y *= invScalar
	} else {
		this.X = 0
		this.Y = 0
	}
	return this
}

func (this *Vector2) Min(v *Vector2) *Vector2 {

	if this.X > v.X {
		this.X = v.X
	}
	if this.Y > v.Y {
		this.Y = v.Y
	}
	return this
}

func (this *Vector2) Max(v *Vector2) *Vector2 {

	if this.X < v.X {
		this.X = v.X
	}
	if this.Y < v.Y {
		this.Y = v.Y
	}
	return this
}

func (this *Vector2) Clamp(min, max *Vector2) *Vector2 {

	// This function assumes min < max, if this assumption isn't true it will not operate correctly
	if this.X < min.X {
		this.X = min.X
	} else if this.X > max.X {
		this.X = max.X
	}

	if this.Y < min.Y {
		this.Y = min.Y
	} else if this.Y > max.Y {
		this.Y = max.Y
	}
	return this
}

func (this *Vector2) ClampScalar(minVal, maxVal float32) *Vector2 {

	min := NewVector2(0, 0)
	max := NewVector2(0, 0)
	min.Set(minVal, minVal)
	max.Set(maxVal, maxVal)
	return this.Clamp(min, max)
}

func (this *Vector2) Floor() *Vector2 {

	this.X = Floor(this.X)
	this.Y = Floor(this.Y)
	return this
}

func (this *Vector2) Ceil() *Vector2 {

	this.X = Ceil(this.X)
	this.Y = Ceil(this.Y)
	return this
}

func (this *Vector2) Round() *Vector2 {

	// TODO NEED CHECK
	this.X = Floor(this.X + 0.5)
	this.Y = Floor(this.Y + 0.5)
	return this
}

func (this *Vector2) RoundToZero() *Vector2 {

	if this.X < 0 {
		this.X = Ceil(this.X)
	} else {
		this.X = Floor(this.X)
	}
	if this.Y < 0 {
		this.Y = Ceil(this.Y)
	} else {
		this.Y = Floor(this.Y)
	}
	return this
}

func (this *Vector2) Negate() *Vector2 {

	this.X = -this.X
	this.Y = -this.Y
	return this
}

func (this *Vector2) Dot(v *Vector2) float32 {

	return this.X*v.X + this.Y*v.Y
}

func (this *Vector2) LengthSq() float32 {

	return this.X*this.X + this.Y*this.Y
}

func (this *Vector2) Length() float32 {

	return Sqrt(this.X*this.X + this.Y*this.Y)
}

func (this *Vector2) Normalize() *Vector2 {

	return this.DivideScalar(this.Length())
}

func (this *Vector2) DistanceTo(v *Vector2) float32 {

	return Sqrt(this.DistanceToSquared(v))
}

func (this *Vector2) DistanceToSquared(v *Vector2) float32 {

	dx := this.X - v.X
	dy := this.Y - v.Y
	return dx*dx + dy*dy
}

func (this *Vector2) SetLength(l float32) *Vector2 {

	oldLength := this.Length()
	if oldLength != 0 && l != oldLength {
		this.MultiplyScalar(l / oldLength)
	}
	return this
}

func (this *Vector2) Lerp(v *Vector2, alpha float32) *Vector2 {

	this.X += (v.X - this.X) * alpha
	this.Y += (v.Y - this.Y) * alpha
	return this
}

func (this *Vector2) LerpVectors(v1, v2 *Vector2, alpha float32) *Vector2 {

	this.SubVectors(v2, v1).MultiplyScalar(alpha).Add(v1)
	return this
}

func (this *Vector2) Equals(v *Vector2) bool {

	return (v.X == this.X) && (v.Y == this.Y)
}

func (this *Vector2) FromArray(array []float32, offset int) *Vector2 {

	this.X = array[offset]
	this.Y = array[offset+1]
	return this
}

func (this *Vector2) ToArray(array []float32, offset int) []float32 {

	array[offset] = this.X
	array[offset+1] = this.Y
	return array
}

// TODO attribute ???
//func (this *Vector2) FromAttribute(attribute, index, offset) *Vector2 {
//
//
//
//
//}

func (this *Vector2) Close() *Vector2 {

	return NewVector2(this.X, this.Y)
}
