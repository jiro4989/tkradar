package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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
		fmt.Println(c.Name)
		fmt.Println(c.Params)
		WriteSVG(os.Stdout, c)
	}
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
