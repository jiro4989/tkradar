package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"

	svg "github.com/ajstarks/svgo"
)

func main() {
	b, err := ioutil.ReadFile("testdata/Classes.json")
	if err != nil {
		panic(err)
	}
	var classes Classes
	if err := json.Unmarshal(b, &classes); err != nil {
		panic(err)
	}
	for _, c := range classes {
		// fmt.Println(c.Name)
		// fmt.Println(c.Params)
		// WriteSVG(os.Stdout, c)

		r, w, h := 50, 100, 100
		xs, ys := PolygonXYs(c, r, w, h)
		fmt.Println("x:", xs)
		fmt.Println("y:", ys)
		fmt.Println("------------")
	}
}

// PolygonXYs はClassのパラメータからX,Y座標のスライスを返す
func PolygonXYs(c Class, r, w, h int) (xs []int, ys []int) {
	var (
		cx = float64(w / 2) // 中心x座標
		cy = float64(h / 2) // 中心y座標
		fr = float64(r)     // sin/cos計算のためにfloatに型を合わせるため
	)
	for i := 0; i < len(c.Params); i++ {
		n := float64(360 / len(c.Params) * i)
		x := fr*math.Cos(n*math.Pi/180) + cx
		y := fr*math.Sin(n*math.Pi/180) + cy
		xs = append(xs, int(x))
		ys = append(ys, int(y))
	}
	return
}

func WriteSVG(w io.Writer, c Class) {
	var (
		width  = 100
		height = 100
		title  = "test"
	)
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, 100)
	canvas.Text(width/2, height/2, title, "text-anchor:middle;font-size:30px;fill:white")
	canvas.End()
}
