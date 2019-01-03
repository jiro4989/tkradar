package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
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
			pr   = r - 50.0
			w, h = r * 2, r * 2
			cp   = point.Point{X: w / 2, Y: h / 2}
		)
		paramPP := point.FetchPolygonPoint(c.FetchLastParams(), ParamMaxes, pr, cp)
		titlePP := point.RegularPolygonPoint(pr, w, h, 8, cp)
		textPP := point.RegularPolygonPoint(r-50, w, h, 8, cp)
		WriteSVG(os.Stdout, r*2, r*2, cp, paramPP, titlePP, textPP)
		break
	}
}

func WriteSVG(wr io.Writer, w, h float64, cp point.Point, paramPP, titlePP, textPP point.PolygonPoint) {
	var (
		wi    = int(w)
		hi    = int(h)
		angle = -90.0
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
		p := point.RegularPolygonPoint(r, w, h, len(paramNames), cp)
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
