package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/jiro4989/tkradar/point"
	"github.com/spf13/cobra"
)

var (
	paramNamesEN = []string{"MHP", "MMP", "ATK", "DEF", "MAT", "MDF", "AGI", "LUK"}
	paramNamesJA = []string{"最大HP", "最大MP", "攻撃力", "防御力", "魔法攻撃", "魔法防御", "敏捷性", "運"}
	paramNames   = paramNamesEN
)

type PolygonPosition struct {
	X, Y []int
}

func init() {
	cobra.OnInitialize()
	RootCommand.Flags().Float64P("width", "w", 500.0, "SVG width")
	RootCommand.Flags().Float64P("text-interval", "t", 50, "Interval")
	RootCommand.Flags().Float64P("frame-interval", "f", 90, "Interval")
	RootCommand.Flags().StringP("language", "l", "en", "Parameters language [en|ja]")
	RootCommand.Flags().StringP("out-filename-format", "o", "class%03d.svg", "Out filename format")
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

		lang, err := flags.GetString("language")
		checkErr(err)

		outfmt, err := flags.GetString("out-filename-format")
		checkErr(err)

		if strings.ToLower(lang) == "ja" {
			paramNames = paramNamesJA
		}

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
				framePP := point.RegularPolygonPoint(paramR, w, h, 8, cp)
				textPP := point.RegularPolygonPoint(textR, w, h, 8, cp)

				func() {
					fn := fmt.Sprintf(outfmt, i)
					wr, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
					checkErr(err)
					defer wr.Close()
					WriteSVG(wr, w, h, cp, paramR, paramPP, framePP, textPP)
				}()
			}
		}
	},
}

func WriteSVG(wr io.Writer, w, h float64, cp point.Point, paramR float64, paramPP, framePP, textPP point.PolygonPoint) {
	var (
		wi    = int(w)
		hi    = int(h)
		angle = -90.0
	)
	canvas := svg.New(wr)
	canvas.Start(wi, hi)

	// 外枠の描画
	canvas.Polygon(framePP.Xs().Int(), framePP.Ys().Int(), "fill:#FAFAFA; stroke:#BDBDBD; ")

	// パラメータ線の描画
	canvas.Polygon(paramPP.Rotate(angle, cp).Xs().Int(), paramPP.Rotate(angle, cp).Ys().Int(), "fill:#BBD9E7; stroke:#91C0DA; stroke-width: 3px;")

	// 中央線の描画
	for _, p := range framePP.Points {
		var (
			x, y   = int(p.X), int(p.Y)
			cx, cy = wi / 2, hi / 2
		)
		canvas.Line(cx, cy, x, y, "fill:#FAFAFA; stroke:#BDBDBD; ")
	}

	// 等間隔基準線の描画
	for i := 0; i < 5; i++ {
		r := paramR * float64(i) / 5
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
