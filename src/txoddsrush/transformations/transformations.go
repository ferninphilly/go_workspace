//Package transformations alters the data as necessary prior to inserting into the pg database
package transformations

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	ai "txoddsrush/apiinterface"
	db "txoddsrush/dbconnection"
	mc "txoddsrush/myconfig"

	"github.com/kr/pretty"
	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

//Teams is the struct to contain our teams

//ChartCreate will contain data for the charts
type ChartCreate struct {
	LastUpdated string
	MatchName   string
	Bookies     []ChartData
}

type ChartData struct {
	BookieName string
	Ival       []float64
	HTeamWin   []float64
	ATeamWin   []float64
	Draw       []float64
}

var fo = ai.ReturnFeedOdds()
var cc ChartCreate

func getMaxNNelements() int {
	y := 0
	for _, v := range fo.Match[0].Bookmaker {
		if v.Offer.Odds != nil {
			y++
		}
	}
	return y
}

//createCharts is how we populate the ChartCreate struct with data from the api
func (cc *ChartCreate) createCharts() {
	i := 0
	z := 0
	cd := make([]ChartData, getMaxNNelements()) //Create our slice of chartdata
	cc.MatchName = "Match: " + fo.Match[0].Hteam + " (Home) vs " + fo.Match[0].Ateam + " (Away)"
	for _, sv := range fo.Match[0].Bookmaker {
		if sv.Offer.Odds != nil {
			cd[i].BookieName = sv.Attributes.Name
			for _, ssv := range sv.Offer.Odds {
				ival, err := strconv.ParseFloat(ssv.Attributes.I, 64)
				mc.HandleError(err)
				cd[i].Ival = append(cd[i].Ival, ival)
				oddsH, err1 := strconv.ParseFloat(ssv.O1, 64)
				mc.HandleError(err1)
				cd[i].HTeamWin = append(cd[i].HTeamWin, oddsH)
				oddsA, err2 := strconv.ParseFloat(ssv.O2, 64)
				mc.HandleError(err2)
				cd[i].ATeamWin = append(cd[i].ATeamWin, oddsA)
				oddsD, err3 := strconv.ParseFloat(ssv.O3, 64)
				mc.HandleError(err3)
				cd[i].Draw = append(cd[i].Draw, oddsD)
			}
			cc.Bookies = cd
			if z > len(cd[i].Ival) {
				z = len(cd[i].Ival)
			}
			i++
		}
	}
	fmt.Println(z)

}

func TransformChartData() []chart.ContinuousSeries {
	cc.createCharts()
	fmt.Println(len(cc.Bookies))
	for _, v := range cc.Bookies {
		pretty.Print(v)
	}
	colors := [9]drawing.Color{drawing.ColorBlack,
		drawing.ColorBlue,
		drawing.ColorGreen,
		drawing.ColorRed,
		drawing.ColorFromAlphaMixedRGBA(uint32(85), uint32(107), uint32(47), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(70), uint32(130), uint32(180), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(210), uint32(105), uint32(30), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(112), uint32(128), uint32(144), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(255), uint32(0), uint32(255), uint32(255)),
	}
	var sries = make([]chart.ContinuousSeries, len(cc.Bookies))
	i := 0
	for _, v := range cc.Bookies {
		if v.Ival != nil {
			ser := chart.ContinuousSeries{
				Name: v.BookieName,
				Style: chart.Style{
					Show:        true,
					StrokeColor: colors[i],
					//FillColor:   chart.GetDefaultColor(0),
				},
				XValues: v.Ival,
				YValues: v.HTeamWin,
			}
			sries[i] = ser
			i++
		}
	}
	return sries
}

//InsertBookies is basically how I created the Bookies table.
//It will error out for primary key constraints if you try to do it twice.
func InsertBookies() string {
	query := "dbconnection/sqlQueries/InsertBookies.sql"
	a := ai.ReturnFeedOdds()
	ub := make(map[string]string) //for "unique bookies"
	cfg := mc.ReturnConfig()
	for _, v := range a.Match {
		for _, val := range v.Bookmaker {
			if _, ok := ub[val.Attributes.Bid]; !ok {
				ub[val.Attributes.Bid] = val.Attributes.Name
			}
		}
	}
	for k, v := range ub {
		db.RunQuery(cfg, query, k, v)
		result := fmt.Sprintf("Entered bookie %s into table", v)
		fmt.Println(result)
	}
	return "Completed entry for all bookies"
}

//InsertMatches is how I inserted into the "dim_matches" table
func InsertMatches() string {
	filename := "dbconnection/sqlQueries/InsertMatches.sql"
	results := ai.ReturnFeedOdds()
	cfg := mc.ReturnConfig()
	lu := results.Attributes.Timestamp
	for _, v := range results.Match {
		gt := strings.Split(v.Time, "T")
		db.RunQuery(cfg, filename, v.Attributes.ID, v.Group, v.Hteam, v.Ateam, gt[0], gt[1][:5], lu)
	}
	return fmt.Sprintf("Completed entries for all matches for the next %s days", cfg.APIData.URLOptions["days"])
}

//InsertTeams populates the dim_teams table
func InsertTeams() string {
	now := int32(time.Now().UTC().Unix())
	query := "dbconnection/sqlQueries/InsertTeamNames.sql"
	results := ai.ReturnTeamList()
	cfg := mc.ReturnConfig()
	for _, v := range results.Competitors.Competitor {
		db.RunQuery(cfg, query, v.ID, v.Name, v.Group, v.Countryid, v.Sportid, fmt.Sprint(now))
		fmt.Println("Completed insert for " + v.Name)
	}
	return "Successfully populated the dim_teams table"
}

//InsertOdds is how I am updating the odds settings.
//I will obviously need to update to do upserts because I don't want dupe data.
func InsertOdds() {
	now := int32(time.Now().UTC().Unix())
	cfg := mc.ReturnConfig()
	d := ai.ReturnFeedOdds()
	i := 0
	for _, v := range d.Match {
		for _, val := range v.Bookmaker {
			for _, sv := range val.Offer.Odds {
				db.RunQuery(cfg, "dbconnection/sqlQueries/InsertOdds.sql", val.Offer.Attributes.ID,
					val.Attributes.Bid, v.Attributes.ID, sv.Attributes.I, sv.O1, sv.O3, sv.O2,
					fmt.Sprint(now), fmt.Sprint(now))
				fmt.Println("Completed insert for row: " + strconv.Itoa(i))
				i++
			}
		}

	}
}
