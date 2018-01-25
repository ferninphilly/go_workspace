//Package visualizations is where I will store the web server and all visualizations
package visualizations

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	tx "txoddsrush/transformations"

	mc "txoddsrush/myconfig"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

const viewdir = "./visualizations/webpages/"

var hometeam, matchname, chartData = tx.TransformChartData()

func handlerForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(viewdir + "index.html")
	mc.HandleError(err)
	err2 := t.Execute(w, matchname)
	mc.HandleError(err2)
}

func drawChart() [][]chart.Chart {
	allchrt := make([][]chart.Chart, len(matchname))
	for k, v := range matchname {
		homeGraph := chart.Chart{
			Title:      v,
			TitleStyle: chart.StyleShow(),
			Background: chart.Style{
				Padding: chart.Box{
					Top:    40,
					Left:   145,
					Right:  10,
					Bottom: 10,
				},
				FillColor: drawing.ColorFromHex("efefef"),
			},
			XAxis: chart.XAxis{
				Name:           "Times odds have changed",
				NameStyle:      chart.StyleShow(),
				Style:          chart.StyleShow(),
				ValueFormatter: chart.TimeHourValueFormatter,
			},

			YAxis: chart.YAxis{
				Name:      "Odds of " + hometeam[k] + " Winning",
				NameStyle: chart.StyleShow(),
				Style:     chart.StyleShow(),
			},
			Series: chartData[k],
		}

		homeGraph.Elements = []chart.Renderable{
			chart.LegendLeft(&homeGraph),
		}

		//Now we make a second chart called awayGraph for the away teams
		awayGraph := chart.Chart{
			Title:      v,
			TitleStyle: chart.StyleShow(),
			Background: chart.Style{
				Padding: chart.Box{
					Top:    40,
					Left:   145,
					Right:  10,
					Bottom: 10,
				},
				FillColor: drawing.ColorFromHex("efefef"),
			},
			XAxis: chart.XAxis{
				Name:           "Times odds have changed",
				NameStyle:      chart.StyleShow(),
				Style:          chart.StyleShow(),
				ValueFormatter: chart.TimeHourValueFormatter,
			},

			YAxis: chart.YAxis{
				Name:      "Odds of " + hometeam[k] + " Winning",
				NameStyle: chart.StyleShow(),
				Style:     chart.StyleShow(),
			},
			Series: chartData[k],
		}

		awayGraph.Elements = []chart.Renderable{
			chart.LegendLeft(&awayGraph),
		}
		allchrt[k] = []chart.Chart{homeGraph, awayGraph}
	}
	return allchrt
}

func matchHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	mc.HandleError(err)
	renderChart(w, r, r.Form["matchid"][0])
}

func renderChart(w http.ResponseWriter, r *http.Request, matchid string) {
	allchrt := drawChart()
	var chrtn int
	chrtn, errco := strconv.Atoi(matchid)
	mc.HandleError(errco)
	w.Header().Set("Content-Type", "image/png")
	var buff bytes.Buffer
	allchrt[chrtn][0].Render(chart.PNG, &buff)
	allchrt[chrtn][1].Render(chart.PNG, &buff)
	if _, err := w.Write(buff.Bytes()); err != nil {
		mc.HandleError(err)
	}
}

//RunServer is basically how we start and run our server.
func RunServer() {
	http.HandleFunc("/", handlerForm)
	http.HandleFunc("/matchHandler", matchHandler)

	//http.HandleFunc("/save/", makeHandler(saveHandler))
	fmt.Println("Server is up and running on localhost:9090")
	http.ListenAndServe(":9090", nil)

}
