package point

import "math"

type Point struct {
	X, Y float64
}
type Points []Point

// Rotate は任意の角度回転した座標を返す。
// cpは原点座標
func (p Point) Rotate(angle float64, cp Point) Point {
	p.X -= cp.X
	p.Y -= cp.Y

	var (
		sin, cos = math.Sincos(angle * math.Pi / 180)
		nx       = cos*p.X - sin*p.Y
		ny       = sin*p.X + cos*p.Y
	)

	p.X = nx + cp.X
	p.Y = ny + cp.Y
	return p
}

type PolygonPoint struct {
	Points
}

func (p PolygonPoint) Rotate(angle float64, cp Point) (np PolygonPoint) {
	l := len(p.Points)
	np.Points = make(Points, len(p.Points))
	for i := 0; i < l; i++ {
		np.Points[i] = p.Points[i].Rotate(angle, cp)
	}
	return
}

func (p PolygonPoint) Xs() (x float64slice) {
	l := len(p.Points)
	x = make([]float64, l)
	for i := 0; i < l; i++ {
		x[i] = p.Points[i].X
	}
	return
}

func (p PolygonPoint) Ys() (y float64slice) {
	l := len(p.Points)
	y = make([]float64, l)
	for i := 0; i < l; i++ {
		y[i] = p.Points[i].Y
	}
	return
}

type float64slice []float64

func (f float64slice) Int() (i []int) {
	i = make([]int, len(f))
	for n, v := range f {
		i[n] = int(v)
	}
	return
}

// RegularPolygonPoint は正多角形のポリゴンの座標を返す
func RegularPolygonPoint(r, w, h float64, polygonCount int) (pp PolygonPoint) {
	var (
		cx     = w / 2
		cy     = h / 2
		radian = math.Pi / 180
	)
	for i := 0; i < polygonCount; i++ {
		var (
			n     = float64(360 / polygonCount * i)
			theta = n * radian
			x     = r*math.Cos(theta) + cx
			y     = r*math.Sin(theta) + cy
		)
		pp.Points = append(pp.Points, Point{X: x, Y: y})
	}
	return
}
