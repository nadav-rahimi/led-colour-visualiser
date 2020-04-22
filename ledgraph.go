package lcv

import (
	"fmt"
	"github.com/cpmech/gosl/utl"
	"github.com/wcharczuk/go-chart"
	"log"
	"os"
)

func createGraph(freqseries ...*[]int) {

	individualSeries := make([]chart.Series, len(freqseries))

	for i := 0; i < len(freqseries); i++ {
		xValues := utl.LinSpace(0, 1, len(*freqseries[i]))
		yValues := *freqseries[i]
		yValues64 := make([]float64, len(yValues))

		for j := 0; j < len(yValues); j++ {
			yValues64[j] = float64(yValues[j])
		}

		individualSeries[i] = chart.ContinuousSeries{
			Name:    fmt.Sprintf("%vhuuuuuu", i),
			XValues: xValues,
			YValues: yValues64,
		}
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "",
		},
		YAxis: chart.YAxis{
			Name: "Frequency",
		},
		Series: individualSeries,
	}

	graph.XAxis.TickStyle.Hidden = true
	graph.YAxisSecondary.Style.Hidden = true
	graph.Height = 400 * 20
	graph.Width = 1024 * 20

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)

	log.Printf("Printing the graph done")

}
