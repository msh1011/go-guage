package guage

import (
	"github.com/lucasb-eyer/go-colorful"
)

// Gradient is a list of color steps
type Gradient []struct {
	Col colorful.Color
	Pos float64
}

func (g Gradient) getColor(t float64) colorful.Color {
	for i := 0; i < len(g)-1; i++ {
		c1 := g[i]
		c2 := g[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			return c1.Col.BlendRgb(c2.Col, t)
		}
	}
	return g[len(g)-1].Col
}

// GetDefaultGradient returns the default gradient
func GetDefaultGradient() Gradient {
	return Gradient{
		{Col: colorful.Color{R: 1, G: 0, B: 0}, Pos: 0},
		{Col: colorful.Color{R: 0, G: 1, B: 0}, Pos: 1},
	}
}
