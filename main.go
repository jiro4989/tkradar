package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/jiro4989/tkradar/point"
)

var (
	paramNames = []string{"MHP", "MMP", "ATK", "DEF", "MAT", "MDF", "AGI", "LUK"}
)

type PolygonPosition struct {
	X, Y []int
}

func main() {
	b, err := ioutil.ReadFile("testdata/Classes2.json")
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
		var (
			r    = 250.0
			w, h = r * 2, r * 2
			cp   = point.Point{X: w / 2, Y: h / 2}
		)
		titlePP := point.RegularPolygonPoint(r, w, h, 8)
		paramPos := PolygonPoint(c.Params, r, cp)
		fmt.Println(titlePP)
		fmt.Println(paramPos)
		//paramPos, titlePos := PolygonXYs(c, r, w, h)
		//WriteSVG(os.Stdout, "test", w, h, paramPos, titlePos, paramNames, titlePP)
		break
	}
}

// PolygonXYs はClassのパラメータからX,Y座標のスライスを返す
func PolygonXYs(c Class, r, w, h int) (paramPos PolygonPosition, titlePos PolygonPosition) {
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
			titleX = math.Cos(theta)
			titleY = math.Sin(theta)
		)
		// 90度半時計回りに回転させる
		// See. https://mathwords.net/heimenkaiten
		rad := 270 * radian
		ntx := titleFr*(titleX*math.Cos(rad)-titleY*math.Sin(rad)) + cx
		nty := titleFr*(titleX*math.Sin(rad)+titleY*math.Cos(rad)) + cy

		titlePos.X = append(titlePos.X, int(ntx))
		titlePos.Y = append(titlePos.Y, int(nty))

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

// RegularPolygon は正多角形のポリゴンの座標を返す
func RegularPolygon(r, w, h float64, polygonCount int) (paramPos PolygonPosition) {
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
		paramPos.X = append(paramPos.X, int(x))
		paramPos.Y = append(paramPos.Y, int(y))
	}
	return
}

// PolygonPoint はdataから座標を返す
func PolygonPoint(data [][]float64, r float64, cp point.Point) (pp point.PolygonPoint) {
	var (
		radian       = math.Pi / 180
		polygonCount = len(data)
	)
	for i := 0; i < polygonCount; i++ {
		// 座標計算に必要な半径はパラメータの値で都度異なるため、rを都度更新
		max := 255.0
		switch i {
		case 0: // MHP
			max = 9999.0
		case 1: // MMP
			max = 2000.0
		case 6, 7: // AGI, LUK
			max = 500.0
		}
		p := data[i]
		last := p[len(p)-1]
		nr := r * last / max

		var (
			n     = float64(360 / polygonCount * i)
			theta = n * radian
			x     = nr*math.Cos(theta) + cp.X
			y     = nr*math.Sin(theta) + cp.Y
		)
		pp.Points = append(pp.Points, point.Point{X: x, Y: y})
	}
	return
}

func WriteSVG(wr io.Writer, title string, w, h int, paramPos, titlePos PolygonPosition, paramNames []string, titlePP point.PolygonPoint) {
	canvas := svg.New(wr)
	canvas.Start(w, h)
	canvas.Circle(w/2, h/2, 100)
	// 外枠の描画
	canvas.Polygon(titlePP.Xs().Int(), titlePP.Ys().Int(), "fill:#FAFAFA; stroke:#BDBDBD; ")
	// canvas.Polygon(titlePos.X, titlePos.Y, "fill:#FAFAFA; stroke:#BDBDBD; ")
	// パラメータ線の描画
	canvas.Polygon(paramPos.X, paramPos.Y, "fill:#BBD9E7; stroke:#91C0DA; stroke-width: 3px;")
	// 中央線の描画
	for i, x := range titlePos.X {
		cx, cy := w/2, h/2
		y := titlePos.Y[i]
		canvas.Line(cx, cy, x, y, "fill:#FAFAFA; stroke:#BDBDBD; ")
	}
	// 等間隔基準線の描画
	for i := 0; i < 5; i++ {
		r := w / 2 * i / 5
		p := RegularPolygon(float64(r), float64(w), float64(h), len(paramNames))
		canvas.Polygon(p.X, p.Y, "fill:none; stroke:#BDBDBD;")
	}
	// canvas.Text(w/2, h/2, title, "text-anchor:middle; font-size:30px; fill:white;")
	for i := 0; i < len(titlePos.X); i++ {
		x := titlePos.X[i]
		y := titlePos.Y[i]
		t := paramNames[i]
		canvas.Text(x, y, t, "text-anchor:middle; font-size:30px; fill:black;")
	}
	canvas.End()
}
