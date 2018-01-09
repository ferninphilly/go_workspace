//Package visualizations is where I will store the web server and all visualizations
package visualizations

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	mc "txoddsrush/myconfig"
	tx "txoddsrush/transformations"

	chart "github.com/wcharczuk/go-chart"
)

//rootdir is the root directory for the pages to be served.
const viewdir = "./visualizations/webpages/"

var cc tx.ChartCreate
var templates = template.Must(template.ParseFiles(viewdir+"edit.html", viewdir+"view.html", viewdir+"index.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

//Page struct holds the data for page data
type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func indexHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "index", p)

}

func drawChart(w http.ResponseWriter, req *http.Request) {
	var x = tx.TransformChartData()
	graph := chart.Chart{

		XAxis: chart.XAxis{
			Name:      "The XAxis",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},

		YAxis: chart.YAxis{
			Name:      "The YAxis",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(),
		},

		Series: []chart.Series{
			x[0],
			x[1],
			x[2],
			x[3],
			x[4],
		},
	}
	w.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, w)
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

func loadPage(title string) (*Page, error) {
	filename := viewdir + title + ".html"
	body, err := ioutil.ReadFile(filename)
	mc.HandleError(err)
	return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
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
}

//Wrapper for our handlers to simplify things downstream
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

//RunServer is basically how we start and run our server.
func RunServer() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/", drawChart)
	//http.HandleFunc("/save/", makeHandler(saveHandler))
	http.ListenAndServe(":9090", nil)
}
