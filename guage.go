package guage

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/acarl005/stripansi"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/qeesung/image2ascii/convert"
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/softwarebackend"
)

var (
	emptyLinesRe = regexp.MustCompile(`^\s+$`)
	emptyColor   = colorful.Color{R: 0.062, G: 0.062, B: 0.062}
)

// Guage holds the guage configuration
type Guage struct {
	// Witdth hold the max # of characters for one line
	// if -1, then the width is calculated based on the height
	Width int
	// Height hold the max # of lines
	// if -1, then the height is calculated based on the width
	Height int
	// Color determines if the guage should be colorized
	Color bool

	// ShowLabel determines if the label should be shown
	ShowLabel bool

	// Min hold the minimum value, values below this value will not be visible
	Min float64
	// Max hold the maximum value, values above this value will not be visible
	Max float64

	// InnerFillPct holds the percentage of the guage that
	// should be removed from the center
	InnerFillPct float64

	// Gradient holds the gradient to use, see GetDefaultGradient()
	Gradient Gradient

	label string
}

// NewGuage creates a new guage with the default options
func NewGuage(name string) *Guage {
	g := &Guage{
		Width:  50,
		Height: 20,
		Color:  true,

		ShowLabel:    true,
		label:        name,
		Min:          0.0,
		Max:          1.0,
		InnerFillPct: 0.4,

		Gradient: GetDefaultGradient(),
	}

	return g
}

// SetGradient sets the gradient to use
func (g *Guage) SetGradient(gr Gradient) {
	g.Gradient = gr
}

// SetSize sets the size of the guage
func (g *Guage) SetSize(w, h int) {
	g.Width = w
	g.Height = h
	if g.Height == -1 {
		g.Height = int(float64(g.Width) * 0.4)
	} else if g.Width == -1 {
		g.Width = int(float64(g.Height) / 0.4)
	}
}

// Render renders the guage from the given value
func (g *Guage) Render(val float64) string {
	ret := g.renderGuage(val)
	if g.ShowLabel {
		ret += "\n" + g.renderLabel(val)
	}

	return ret
}

func (g *Guage) renderLabel(val float64) string {
	tx := fmt.Sprintf("%s: %.2f", g.label, val)
	s := g.Width/2 - len(tx)/2
	if s < 0 {
		s = 0
	}
	return fmt.Sprintf("%s%s%s",
		strings.Repeat(" ", s), tx, strings.Repeat(" ", s))
}

func (g *Guage) renderGuage(val float64) (ret string) {
	convertOptions := convert.DefaultOptions

	convertOptions.Colored = g.Color
	convertOptions.FixedWidth = g.Width
	convertOptions.FixedHeight = g.Height

	if g.Height == -1 {
		convertOptions.FixedWidth = g.Width
		convertOptions.FixedHeight = int(float64(convertOptions.FixedWidth) * 0.4)
	} else if g.Width == -1 {
		convertOptions.FixedHeight = g.Height
		convertOptions.FixedWidth = int(float64(convertOptions.FixedHeight) / 0.4)
	}

	cvSize := int(math.Max(
		float64(convertOptions.FixedWidth),
		float64(convertOptions.FixedHeight),
	))
	backend := softwarebackend.New(cvSize, cvSize)
	cv := canvas.New(backend)

	w, h := float64(cv.Width()), float64(cv.Height())
	cv.FillRect(0, 0, w, h)

	pi := 3.14
	step := pi / 48
	start := pi - (step * 6)
	end := 2*3.14 + (step * 6)

	for r := start; r < end; r += step {
		pct := (r - start) / (end - start)
		c := g.Gradient.getColor(pct)
		if pct > (val-g.Min)/(g.Max-g.Min) {
			c = emptyColor
		}

		cv.SetFillStyle(c.R, c.G, c.B)
		cv.BeginPath()
		cv.MoveTo(w*0.5, h*0.5)
		cv.Arc(w*0.5, h*0.5, h*0.5, r, r+math.Pi/24, false)
		cv.ClosePath()
		cv.Fill()
		cv.SetFillStyle(0, 0, 0)
		cv.BeginPath()
		cv.MoveTo(w*0.5, h*0.5)
		cv.Arc(w*0.5, h*0.5, math.Min(w, h)*0.5*g.InnerFillPct, r, r+math.Pi/24, false)
		cv.ClosePath()
		cv.Fill()
	}

	// Create the image converter
	converter := convert.NewImageConverter()
	s := converter.Image2ASCIIString(
		backend.Image.SubImage(backend.Image.Rect),
		&convertOptions,
	)

	for _, r := range strings.Split(s, "\n") {
		r1 := stripansi.Strip(r)
		if emptyLinesRe.MatchString(r1) {
			continue
		}
		ret += r + "\n"
	}
	return strings.TrimRight(ret, "\n")
}
