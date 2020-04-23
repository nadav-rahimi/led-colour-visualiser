package lcv

import (
	"github.com/cpmech/gosl/utl"
	"github.com/wcharczuk/go-chart"
	"log"
	"os"
)

// TODO add ability to name f series
func createGraph(freqseries ...*[]int) {

	// Create the series object which will be plotted on the graph
	individualSeries := make([]chart.Series, len(freqseries))

	// Populates the series objects with the buffer data from each f series
	for i := 0; i < len(freqseries); i++ {
		xValues := utl.LinSpace(0, 1, len(*freqseries[i]))
		yValues := *freqseries[i]
		yValues64 := make([]float64, len(yValues))

		for j := 0; j < len(yValues); j++ {
			yValues64[j] = float64(yValues[j])
		}

		individualSeries[i] = chart.ContinuousSeries{
			XValues: xValues,
			YValues: yValues64,
		}
	}

	// Creates the graph object and populates it with the series
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "",
		},
		YAxis: chart.YAxis{
			Name: "Frequency (Hz)",
		},
		Series: individualSeries,
	}

	// Styles the graph
	graph.XAxis.TickStyle.Hidden = true
	graph.YAxisSecondary.Style.Hidden = true
	// TODO make height and width dynamic
	graph.Height = 400 * 3
	graph.Width = 1024 * 3

	// Render the graph to an output file
	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

	log.Printf("Printing the graph done")

}
