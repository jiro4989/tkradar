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
	X, Y int
}

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
		r := 250
		w, h := r*2, r*2
		paramPos, titlePos := PolygonXYs(c, r, w, h)
		paramPos = PolygonXYs3(c, float64(r), float64(w), float64(h))
		WriteSVG(os.Stdout, "test", w, h, paramPos, titlePos, paramNames)
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

// PolygonXYs3 はClassのパラメータからX,Y座標のスライスを返す
func PolygonXYs3(c Class, r, w, h float64) (paramPos PolygonPosition) {
	var (
		cx           = w / 2
		cy           = h / 2
		radian       = math.Pi / 180
		polygonCount = len(c.Params)
	)
	for i := 0; i < polygonCount; i++ {
		// 座標計算に必要な半径はパラメータの値で都度異なるため、rを都度更新
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
		nr := r * float64(last) / float64(max)

		var (
			n     = float64(360 / polygonCount * i)
			theta = n * radian
			x     = math.Cos(theta)
			y     = math.Sin(theta)
		)
		// 90度半時計回りに回転させる
		// See. https://mathwords.net/heimenkaiten
		rad := 270 * radian
		nx := nr*(x*math.Cos(rad)-y*math.Sin(rad)) + cx
		ny := nr*(x*math.Sin(rad)+y*math.Cos(rad)) + cy

		paramPos.X = append(paramPos.X, int(nx))
		paramPos.Y = append(paramPos.Y, int(ny))
	}
	return
}

func WriteSVG(wr io.Writer, title string, w, h int, paramPos, titlePos PolygonPosition, paramNames []string) {
	canvas := svg.New(wr)
	canvas.Start(w, h)
	canvas.Circle(w/2, h/2, 100)
	// 外枠の描画
	canvas.Polygon(titlePos.X, titlePos.Y, "fill:#FAFAFA; stroke:#BDBDBD; ")
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
