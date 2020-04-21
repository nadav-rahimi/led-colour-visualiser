package main

import (
	"../../led-colour-visualiser"
	"github.com/andlabs/ui"
	"math/rand"
)

// Main
func main() {
	*lcv.Current_colour_hex = rand.Uint32()
	_ = ui.Main(lcv.SetupUI)
}
