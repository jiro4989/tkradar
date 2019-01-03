package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"os"

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
	b, err := ioutil.ReadFile(os.Args[1])
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
			pr   = r - 25.0
			w, h = r * 2, r * 2
			cp   = point.Point{X: w / 2, Y: h / 2}
		)
		paramPP := PolygonPoint(c.Params, pr, cp)
		titlePP := point.RegularPolygonPoint(pr, w, h, 8)
		textPP := point.RegularPolygonPoint(r, w, h, 8)
		WriteSVG(os.Stdout, w, h, paramPP, titlePP, textPP)
		break
	}
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

func WriteSVG(wr io.Writer, w, h float64, paramPP, titlePP, textPP point.PolygonPoint) {
	var (
		wi    = int(w)
		hi    = int(h)
		angle = -90.0
		cp    = point.Point{X: w / 2, Y: h / 2}
	)
	canvas := svg.New(wr)
	canvas.Start(wi, hi)
	canvas.Circle(wi/2, hi/2, 100)
	// 外枠の描画
	canvas.Polygon(titlePP.Xs().Int(), titlePP.Ys().Int(), "fill:#FAFAFA; stroke:#BDBDBD; ")
	// パラメータ線の描画
	canvas.Polygon(paramPP.Rotate(angle, cp).Xs().Int(), paramPP.Rotate(angle, cp).Ys().Int(), "fill:#BBD9E7; stroke:#91C0DA; stroke-width: 3px;")
	// 中央線の描画
	for _, p := range titlePP.Points {
		var (
			x, y   = int(p.X), int(p.Y)
			cx, cy = wi / 2, hi / 2
		)
		canvas.Line(cx, cy, x, y, "fill:#FAFAFA; stroke:#BDBDBD; ")
	}
	// 等間隔基準線の描画
	for i := 0; i < 5; i++ {
		r := (w/2 - 25) * float64(i) / 5
		p := point.RegularPolygonPoint(r, w, h, len(paramNames))
		canvas.Polygon(p.Xs().Int(), p.Ys().Int(), "fill:none; stroke:#BDBDBD;")
	}
	// テキストの描画
	for i, v := range textPP.Rotate(angle, cp).Points {
		var (
			x, y = int(v.X), int(v.Y)
			n    = paramNames[i]
		)
		canvas.Text(x, y, n, "text-anchor:middle; font-size:30px; fill:black;")
	}
	canvas.End()
}
