package lcv

import (
	"errors"
	"github.com/lucasb-eyer/go-colorful"
	"sort"
)

// BOX COLOUR IMPLEMENTATION
// Box colour type which implements methods for converting the colorful type
// straight to uint32
type boxColour colorful.Color

// Converts the box colour type to a uint32 value
func (col boxColour) UINT32() uint32 {
	var r = uint32(col.R*255.0 + 0.5)
	var g = uint32(col.G*255.0 + 0.5)
	var b = uint32(col.B*255.0 + 0.5)

	r = r << 8
	r = r + g
	r = r << 8
	r = r + b

	return r
}

// GRADIENTS
// This table contains the "keypoints" of the colorgradient you want to generate.
// The position of each keypoint has to live in the range [0,1]
type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

// This is the meat of the gradient computation. It returns a HCL-blend between
// the two colors around `t`.
// Note: It relies heavily on the fact that the gradient keypoints are sorted.
func (self GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(self)-1; i++ {
		c1 := self[i]
		c2 := self[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return self[len(self)-1].Col
}

// This is a very nice thing Golang forces you to do!
// It is necessary so that we can write out the literal of the colortable below.
func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}

// Functions which returns a gradient function based on its name
func getGradientTable(s string) (*GradientTable, error) {
	if val, ok := gradients[s]; ok {
		return val, nil
	}

	return &GradientTable{}, errors.New("Gradient name incorrect")
}

// Gradients available to users hardcoded into the application
var gradients = map[string]*GradientTable{
	"starboy": &GradientTable{
		//{MustParseHex("#1a0406"), 0.0},
		{MustParseHex("#1152cb"), 0.0},
		{MustParseHex("#1152cb"), 0.05},
		{MustParseHex("#e4032f"), 0.1},
		{MustParseHex("#f6c507"), 0.55},
		//{MustParseHex("#faf4e6"), 0.7},
		{MustParseHex("#faf6cb"), 1.0},
		//{MustParseHex("#faf4e6"), 1.0},
	},
	"franklake": &GradientTable{
		//{MustParseHex("#007dfe"), 0},
		{MustParseHex("#ff7303"), 0},
		{MustParseHex("#ff7303"), 0.1},
		{MustParseHex("#ffa7e1"), 0.5},
		{MustParseHex("#faf4e6"), 1.0},
	},
	"smiths": &GradientTable{
		{MustParseHex("#ff0202"), 0},
		{MustParseHex("#ff0202"), 0.1},
		{MustParseHex("#ff8d00"), 0.3},
		{MustParseHex("#fff400"), 0.5},
		{MustParseHex("#f1ff00"), 0.8},
		{MustParseHex("#A4ff00"), 1.0},
	},
	"weeknd": &GradientTable{
		//{MustParseHex("#ff0202"), 0.0},
		//{MustParseHex("#ff0258"), 0.03},
		{MustParseHex("#5202fc"), 0.00},
		{MustParseHex("#ff0074"), 0.9},
		//{MustParseHex("#ff0000"), 1.0},
	},
	"shabjdeed": &GradientTable{
		{MustParseHex("#020024"), 0.0},
		{MustParseHex("#ad63f4"), 0.35},
		{MustParseHex("#00d4ff"), 1.0},
	},
}

// Returns a string slice of the names of the hardcoded gradients
func gradientList() []string {
	keys := make([]string, 0, len(gradients))
	for k := range gradients {
		keys = append(keys, k)
	}
	keys = append(keys, "default")
	sort.Strings(keys)

	return keys
}

// Returns the name of a gradient at a given index of the gradient list
func getGradientName(i int) string {
	return gradientList()[i]
}
