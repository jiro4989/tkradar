package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/jiro4989/tkradar/point"
	"github.com/spf13/cobra"
)

var (
	paramNames = []string{"MHP", "MMP", "ATK", "DEF", "MAT", "MDF", "AGI", "LUK"}
)

type PolygonPosition struct {
	X, Y []int
}

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().Float64P("width", "w", 500.0, "SVG width")
	RootCommand.Flags().Float64P("text-interval", "t", 50, "Interval")
	RootCommand.Flags().Float64P("frame-interval", "f", 70, "Interval")
}

var RootCommand = &cobra.Command{
	Use:   "tkradar",
	Short: "tkradar",
	Long:  "tkradar",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Need Classes.json.\nSee 'tkradar -h'.")
			os.Exit(1)
		}

		checkErr := func(err error) {
			if err != nil {
				panic(err)
			}
		}

		flags := cmd.Flags()
		w, err := flags.GetFloat64("width")
		checkErr(err)

		ti, err := flags.GetFloat64("text-interval")
		checkErr(err)

		fi, err := flags.GetFloat64("frame-interval")
		checkErr(err)

		var (
			r      = w / 2
			h      = w
			cp     = point.Point{X: w / 2, Y: h / 2}
			paramR = r - fi
			textR  = r - ti
		)

		for _, f := range args {
			b, err := ioutil.ReadFile(f)
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
				paramPP := point.FetchPolygonPoint(c.FetchLastParams(), ParamMaxes, paramR, cp)
				titlePP := point.RegularPolygonPoint(paramR, w, h, 8, cp)
				textPP := point.RegularPolygonPoint(textR, w, h, 8, cp)
				WriteSVG(os.Stdout, w, h, cp, paramPP, titlePP, textPP)
				break
			}
		}
	},
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
