package main

import (
	"fmt"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/msh1011/go-guage"
)

func main() {
	g := guage.NewGuage("Value")
	g.SetSize(211, 52)
	g.Min = -5
	g.Max = 5
	p := perlin.NewPerlin(2, 2, 10, time.Now().UnixNano())

	x := .5
	for {
		fmt.Printf("\033[H\033[2J%s", g.Render(p.Noise1D(x)*5))
		time.Sleep(time.Second / 10)
		x += 0.05
	}
}
