package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"os"

	svg "github.com/ajstarks/svgo"
)

var (
	paramNames = []string{"MHP", "MMP", "ATK", "DEF", "MAT", "MDF", "AGI", "LUK"}
)

type Position struct {
	X []int
	Y []int
}

// Rotate はN度半時計回りに回転させる
func (p *Position) Rotate(n int) {
	var (
		rad  = math.Pi / 180
		nrad = float64(n) * rad
		cos  = math.Cos(nrad)
		sin  = math.Sin(nrad)
	)
	for i := 0; i < len(p.X); i++ {
		x, y := float64(p.X[i]), float64(p.Y[i])
		x = cos*x - sin*y
		y = sin*x + cos*y
		p.X[i] = int(x)
		p.Y[i] = int(y)
	}
}

func main() {
	b, err := ioutil.ReadFile("testdata/Classes.json")
	if err != nil {
		panic(err)
	}
	var classes Classes
	if err := json.Unmarshal(b, &classes); err != nil {
		panic(err)
	}
	for i, c := range classes {
		if i == 0 {
			// Classes.jsonの最初のデータは絶対にnullのためスキップ
			continue
		}
		r := 250
		w, h := r*2, r*2
		paramPos, titlePos := PolygonXYs(c, r, w, h)
		//paramPos.Rotate(90)
		//titlePos.Rotate(90)
		WriteSVG(os.Stdout, "test", w, h, paramPos, titlePos, paramNames)
		break
	}
}

// PolygonXYs はClassのパラメータからX,Y座標のスライスを返す
func PolygonXYs(c Class, r, w, h int) (paramPos Position, titlePos Position) {
	var (
		cx      = float64(w / 2) // 中心x座標
		cy      = float64(h / 2) // 中心y座標
		fr      = float64(r - 25)
		titleFr = float64(r)
		radian  = math.Pi / 180
	)
	for i := 0; i < len(c.Params); i++ {
		var (
			n      = float64(360 / len(c.Params) * i)
			theta  = n * radian
			x      = fr*math.Cos(theta) + cx
			y      = fr*math.Sin(theta) + cy
			titleX = titleFr*math.Cos(theta) + cx
			titleY = titleFr*math.Sin(theta) + cy
		)

		titlePos.X = append(titlePos.X, int(titleX))
		titlePos.Y = append(titlePos.Y, int(titleY))

		max := 255
		switch i {
		case 0: // MHP
			max = 9999
		case 1: // MMP
			max = 2000
		case 6, 7: // AGI, LUK
			max = 500
		}
		p := c.Params[i]
		last := p[len(p)-1]
		x = x * float64(last) / float64(max)
		y = y * float64(last) / float64(max)

		paramPos.X = append(paramPos.X, int(x))
		paramPos.Y = append(paramPos.Y, int(y))
	}
	return
}

func WriteSVG(wr io.Writer, title string, w, h int, paramPos, titlePos Position, paramNames []string) {
	canvas := svg.New(wr)
	canvas.Start(w, h)
	canvas.Circle(w/2, h/2, 100)
	canvas.Polygon(titlePos.X, titlePos.Y, "fill:white;")
	canvas.Text(w/2, h/2, title, "text-anchor:middle;font-size:30px;fill:white")
	for i := 0; i < len(titlePos.X); i++ {
		x := titlePos.X[i]
		y := titlePos.Y[i]
		t := paramNames[i]
		canvas.Text(x, y, t, "text-anchor:middle;font-size:30px;fill:black")
	}
	canvas.End()
}
