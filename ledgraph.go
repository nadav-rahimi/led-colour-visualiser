package lcv

import (
	"github.com/cpmech/gosl/utl"
	"github.com/wcharczuk/go-chart"
	"log"
	"os"
	"time"
)

// Renders a graph given the names of each series, the elapsed time of
// streaming and each frequency series
func createGraph(seriesnames []string, elapsed_t time.Duration, freqseries ...*[]int) {

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
			Name:    seriesnames[i],
			XValues: xValues,
			YValues: yValues64,
		}
	}

	// Creates the graph object and populates it with the series
	graph := chart.Chart{
		XAxis: chart.XAxis{},
		YAxis: chart.YAxis{
			Name: "Frequency (Hz)",
		},
		Series: individualSeries,
	}

	// Styles the graph
	graph.Height = 400 * 3
	graph.Width = int(60 * elapsed_t.Seconds())
	graph.XAxis.TickStyle.Hidden = true
	graph.YAxisSecondary.Style.Hidden = true

	// Adding the legend to the graph
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	// Render the graph to an output file
	f, _ := os.Create("output.png")
	defer f.Close()
	err := graph.Render(chart.PNG, f)
	chk(err)

	log.Printf("Graph rendered to file successfully")

}
