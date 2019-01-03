package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jiro4989/tkradar/point"
	"github.com/spf13/cobra"
)

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
