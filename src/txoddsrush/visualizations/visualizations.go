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
}

func indexHandler(w http.ResponseWriter, r *http.Request, title string) {
	var p Page
	p.drawChart(w, "First try dude", "A match")
	pretty.Print(p)
	renderTemplate(w, "index", p)
}
*/
func handlerForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./visualizations/webpages/index.html")
	mc.HandleError(err)
	err2 := t.Execute(w, matchname)
	mc.HandleError(err2)
}

func drawChart() []chart.Chart {
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
	buffer := new(bytes.Buffer)
	err := allchrt[chrtn].Render(chart.PNG, buffer)
	mc.HandleError(err)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		mc.HandleError(err)
	}
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
	http.HandleFunc("/", handlerForm)
	http.HandleFunc("/matchHandler", matchHandler)

	//http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":9090", nil)
	fmt.Println("Server is up and running on 9090")
}
