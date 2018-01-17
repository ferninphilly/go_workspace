//Package visualizations is where I will store the web server and all visualizations
package visualizations

import (
	"fmt"
	"net/http"
	mc "txoddsrush/myconfig"
	tx "txoddsrush/transformations"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

type Page struct {
	MatchName string
	Message   string
	Chart     chart.Chart
}

//rootdir is the root directory for the pages to be served.
const viewdir = "./visualizations/webpages/"

var matchname, chartData = tx.TransformChartData()

//var templates = template.Must(template.ParseFiles(viewdir+"edit.html", viewdir+"view.html", viewdir+"index.html"))

//var validPath = regexp.MustCompile("^/(index|edit|save|view)/([a-zA-Z0-9]+)$")

/*func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}*/
/*
func indexHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p Page
	p.drawChart(w, "First try dude", "A match")
	pretty.Print(p)
	renderTemplate(w, "index", p)
} */

func DrawChart() []chart.Chart {
	allchrt := make([]chart.Chart, len(matchname))
	for k, v := range matchname {
		graph := chart.Chart{
			Title:      v,
			TitleStyle: chart.StyleShow(),
			Background: chart.Style{
				Padding: chart.Box{
					Top:    50,
					Left:   25,
					Right:  25,
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
				Name:      "Odds of Home Team Winning",
				NameStyle: chart.StyleShow(),
				Style:     chart.StyleShow(),
			},

			Series: chartData[k],
		}

		graph.Elements = []chart.Renderable{
			chart.LegendLeft(&graph),
		}
		allchrt[k] = graph
	}
	return allchrt
}

func renderChart(w http.ResponseWriter, r *http.Request) {
	allchrt := DrawChart()
	collector := &chart.ImageWriter{}
	allchrt[4].Render(chart.PNG, collector)
	/*	for _, v := range allchrt {
		err := v.Render(chart.PNG, collector)
		mc.HandleError(err)
	} */
	img, err := collector.Image()
	mc.HandleError(err)
	fmt.Fprintf(w, img)
}

/*
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
} */

/*func (*p PloadPage(title string, msg string, chart chart.Chart) (*Page, error) {
	return &Page{Title: title, Message: msg, }, nil
}*/

/*func renderTemplate(w http.ResponseWriter, tmpl string, p Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/*func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
} */

//Wrapper for our handlers to simplify things downstream
/*func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	fmt.Println("Made it here")
	return func(w http.ResponseWriter, r *http.Request) {
			m := validPath.FindStringSubmatch(r.URL.Path)
				fmt.Println(m[2])
				if m == nil {
					http.NotFound(w, r)
					return
				}
		fn(w, r, "index")
	}
} */

//RunServer is basically how we start and run our server.
func RunServer() {
	/*	http.HandleFunc("/view/", makeHandler(viewHandler))
		http.HandleFunc("/edit/", makeHandler(editHandler)) */
	http.HandleFunc("/", renderChart)

	//http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":9090", nil)
	fmt.Println("Server is up and running on 9090")
}
