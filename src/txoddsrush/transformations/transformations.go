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

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

//Teams is the struct to contain our teams

//ChartCreate will contain data for the charts
type ChartCreate struct {
	LastUpdated string
	MatchName   string
	HomeTeam    string
	AwayTeam    string
	Bookies     []ChartData
}

//ChartData will contain the specific structs of data for the charts
type ChartData struct {
	BookieName string
	OddsTime   []time.Time
	Ival       []float64
	HTeamWin   []float64
	ATeamWin   []float64
	Draw       []float64
}

var fo = ai.ReturnFeedOdds()

func getMaxNNelements(rns int) int {
	y := 0
	for _, v := range fo.Match[rns].Bookmaker {
		if v.Offer.Odds != nil {
			for k := range v.Offer.Odds {
				y = k
			}
		}
	}
	return int(y)
}

//CreateCharts is how we populate the ChartCreate struct with data from the api
func CreateCharts() []ChartCreate {
	alls := make([]ChartCreate, len(fo.Match))
	//Create our slice of chartdata pretty.Print(fo)
	for a, el := range fo.Match {
		alls[a].MatchName = "Match: " + el.Hteam + " (Home) vs " + el.Ateam + " (Away)" + "\n Match date: " + el.Time[0:10] + " @ " + el.Time[11:16]
		alls[a].HomeTeam = el.Hteam
		alls[a].AwayTeam = el.Ateam
		cd := make([]ChartData, len(el.Bookmaker))
		layout := "2006-01-02T15:04:05Z07:00"
		for k, sv := range el.Bookmaker {
			if sv.Offer.Odds != nil {
				cd[k].BookieName = sv.Attributes.Name
				for _, ssv := range sv.Offer.Odds {
					betime, err := time.Parse(layout, ssv.Attributes.Time)
					mc.HandleError(err)
					cd[k].OddsTime = append(cd[k].OddsTime, betime)
					ival, err := strconv.ParseFloat(ssv.Attributes.I, 64)
					mc.HandleError(err)
					cd[k].Ival = append(cd[k].Ival, ival)
					oddsH, err1 := strconv.ParseFloat(ssv.O1, 64)
					mc.HandleError(err1)
					cd[k].HTeamWin = append(cd[k].HTeamWin, oddsH)
					oddsA, err2 := strconv.ParseFloat(ssv.O2, 64)
					mc.HandleError(err2)
					cd[k].ATeamWin = append(cd[k].ATeamWin, oddsA)
					oddsD, err3 := strconv.ParseFloat(ssv.O3, 64)
					mc.HandleError(err3)
					cd[k].Draw = append(cd[k].Draw, oddsD)
					fmt.Println("Name: " + sv.Attributes.Name + " Ival: " + ssv.Attributes.I + " HteamWin: " + ssv.O1 + " as of: " + betime.String())
				}
				alls[a].Bookies = cd
			}
		}
	}
	return alls
}

//TransformChartData is the function that puts chart data into the Continuous Series.
func TransformChartData() ([]string, [][]chart.Series) {
	chrt := CreateCharts()
	colors := [12]drawing.Color{
		drawing.ColorBlack,
		drawing.ColorBlue,
		drawing.ColorGreen,
		drawing.ColorRed,
		drawing.ColorFromAlphaMixedRGBA(uint32(85), uint32(107), uint32(47), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(70), uint32(130), uint32(180), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(210), uint32(105), uint32(30), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(112), uint32(128), uint32(144), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(255), uint32(0), uint32(255), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(0), uint32(255), uint32(255), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(0), uint32(128), uint32(128), uint32(255)),
		drawing.ColorFromAlphaMixedRGBA(uint32(0), uint32(255), uint32(0), uint32(255)),
	}
	allsries := make([][]chart.Series, len(chrt))
	var matches = make([]string, len(chrt))
	//	layout := "Mon Jan 02 2006 15:04:05 GMT-0700"
	for k, v := range chrt {
		var sries = make([]chart.TimeSeries, len(v.Bookies))
		var retsries = make([]chart.Series, len(v.Bookies))
		for a, b := range v.Bookies {
			if b.OddsTime != nil {
				//				betime, err := time.Parse(layout, b.OddsTime[0])
				ser := chart.TimeSeries{
					Name: b.BookieName,
					Style: chart.Style{
						Show:        true,
						StrokeColor: colors[a],
						//FillColor:   chart.GetDefaultColor(0),
					},
					XValues: b.OddsTime,
					YValues: b.HTeamWin,
				}
				sries[a] = ser
				retsries[a] = sries[a]
				allsries[k] = append(allsries[k], retsries[a])
			}
		}
		matches[k] = v.MatchName

	}
	return matches, allsries
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
