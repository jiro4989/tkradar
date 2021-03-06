package point

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	type TD struct {
		self   Point
		angle  float64
		cp     Point
		expect Point
		desc   string
	}
	tds := []TD{
		TD{
			self:   Point{X: 10, Y: 10},
			angle:  45,
			cp:     Point{X: 0, Y: 0},
			expect: Point{X: 0, Y: 14},
			desc:   "P(10, 10), rotate 45",
		},
		TD{
			self:   Point{X: 10, Y: 10},
			angle:  90,
			cp:     Point{X: 0, Y: 0},
			expect: Point{X: -10, Y: 10},
			desc:   "P(10, 10), rotate 90",
		},
		TD{
			self:   Point{X: 20, Y: 20},
			angle:  90,
			cp:     Point{X: 10, Y: 10},
			expect: Point{X: 0, Y: 20},
			desc:   "P(20, 20), cp(10, 10), rotate 90",
		},
		TD{
			self:   Point{X: 10, Y: 10},
			angle:  180,
			cp:     Point{X: 0, Y: 0},
			expect: Point{X: -10, Y: -10},
			desc:   "P(10, 10), rotate 180",
		},
		TD{
			self:   Point{X: 10, Y: 10},
			angle:  270,
			cp:     Point{X: 0, Y: 0},
			expect: Point{X: 10, Y: -10},
			desc:   "P(10, 10), rotate 270",
		},
		TD{
			self:   Point{X: 10, Y: 0},
			angle:  30,
			cp:     Point{X: 0, Y: 0},
			expect: Point{X: math.Round(5 * math.Sqrt(3)), Y: 5},
			desc:   "P(10, 0), rotate 30",
		},
		TD{
			self:   Point{X: 10, Y: 0},
			angle:  60,
			cp:     Point{X: 0, Y: 0},
			expect: Point{X: 5, Y: math.Round(5 * math.Sqrt(3))},
			desc:   "P(10, 0), rotate 60",
		},
	}
	for _, v := range tds {
		angle, cp, expect := v.angle, v.cp, v.expect
		got1 := v.self.Rotate(angle, cp)
		assert.Equal(t, expect.X, math.Round(got1.X), v.desc)
		assert.Equal(t, expect.Y, math.Round(got1.Y), v.desc)
	}
}

func TestFetchPolygonPoint(t *testing.T) {
	type TD struct {
		data  []float64
		maxes []float64
		r     float64
		cp    Point
		pp    PolygonPoint
	}
	tds := []TD{
		TD{
			data:  []float64{10.0, 5.0, 0.0, 10.0},
			maxes: []float64{10.0, 10.0, 10.0, 10.0},
			r:     10.0,
			cp:    Point{X: 0.0, Y: 0.0},
			pp: PolygonPoint{
				Points: Points{
					Point{X: 10.0, Y: 0.0},
					Point{X: 0.0, Y: -5.0},
					Point{X: 0.0, Y: 0.0},
					Point{X: 0.0, Y: 10.0},
				},
			},
		},
	}
	for _, v := range tds {
		data, maxes, r, cp, pp := v.data, v.maxes, v.r, v.cp, v.pp
		got1 := FetchPolygonPoint(data, maxes, r, cp)
		// assert.Equal(t, pp, got1)
		fmt.Println(pp, got1)
	}
}
