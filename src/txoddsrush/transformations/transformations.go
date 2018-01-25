//Package transformations alters the data as necessary prior to inserting into the pg database
package transformations

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	ai "txoddsrush/apiinterface"
	mc "txoddsrush/myconfig"
)

//Teams is the struct to contain our teams

//ChartCreate will contain data for the charts
/*
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
} */

//CreateCharts is how we populate the ChartCreate struct with data from the api
/*func CreateCharts() []ChartCreate {
	alls := make([]ChartCreate, len(fo.Match))
	//Create our slice of chartdata pretty.Print(fo)
	for a, el := range fo.Match {
		alls[a].MatchName = el.Hteam + " vs " + el.Ateam + " on " + el.Time[0:10] + " @ " + el.Time[11:16]
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
func TransformChartData() ([]string, []string, [][]chart.Series) {
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
	var hteam = make([]string, len(chrt))
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
		hteam[k] = v.HomeTeam

	}
	return hteam, matches, allsries
}

*/

func stringToTimestamp(strtime string) int64 {
	layout := "2006-01-02T15:04:05+00:00"
	t, err := time.Parse(layout, strtime)
	mc.HandleError(err)
	return t.Unix()
}

//ReadyMatchData prepares interfaces to be inserted in to the fact_matches table
func ReadyMatchData(fo ai.CreateOdds) (string, [][]interface{}) {
	now := int32(time.Now().UTC().Unix())
	cells := make([][]interface{}, len(fo.Match))
	filename := "dbconnection/sqlQueries/InsertMatches.sql"
	if len(fo.Match) > 0 {
		lu, err := strconv.Atoi(fo.Attributes.Timestamp)
		mc.HandleError(err)
		for k, v := range fo.Match {
			gt := strings.Split(v.Time, "T")
			matchTime := stringToTimestamp(v.Time)
			cells[k] = []interface{}{v.Attributes.ID, v.Group, v.Hteam, v.Ateam, gt[0], gt[1][:5], lu, matchTime, now}
		}
		return filename, cells
	}
	fmt.Println("There is an error- the length of the Match struct is 0- which means NO matches. Check config file- maybe increase days?")
	return filename, cells
}

//ReadyBookiesData is basically how I created the Bookies table.
//It will error out for primary key constraints if you try to do it twice.
func ReadyBookiesData(fo ai.CreateOdds) (string, [][]interface{}) {
	now := int32(time.Now().UTC().Unix())
	filename := "dbconnection/sqlQueries/InsertBookies.sql"
	var cells [][]interface{}
	ub := make(map[string]string) //for "unique bookies"
	for _, v := range fo.Match {
		for _, val := range v.Bookmaker {
			if _, ok := ub[val.Attributes.Bid]; !ok {
				ub[val.Attributes.Bid] = val.Attributes.Name
			}
		}
	}
	for k, v := range ub {
		cells = append(cells, []interface{}{k, v, now})
	}
	return filename, cells
}

func calculateMargins(strOdds []string) [7]float64 {
	win, err := strconv.ParseFloat(strOdds[0], 64)
	mc.HandleError(err)
	lose, errl := strconv.ParseFloat(strOdds[1], 64)
	mc.HandleError(errl)
	draw, errd := strconv.ParseFloat(strOdds[2], 64)
	mc.HandleError(errd)
	margin := ((1 / win) + (1 / lose) + (1 / draw)) - 1
	actualWin := (3 * win) / (3 - (margin * win))
	actualLose := (3 * lose) / (3 - (margin * lose))
	actualDraw := (3 * draw) / (3 - (margin * draw))
	return [7]float64{win, lose, draw, margin, actualWin, actualLose, actualDraw}

}

//ReadyOddsData is how I am updating the odds settings.
//I will obviously need to update to do upserts because I don't want dupe data.
func ReadyOddsData(fo ai.CreateOdds) (string, [][]interface{}) {
	now := int32(time.Now().UTC().Unix())
	filename := "dbconnection/sqlQueries/InsertOdds.sql"
	var cells [][]interface{}
	for _, v := range fo.Match {
		for _, val := range v.Bookmaker {
			for _, sv := range val.Offer.Odds {
				strOdds := []string{sv.O1, sv.O3, sv.O2}
				insOdds := calculateMargins(strOdds)
				oddsTime := stringToTimestamp(sv.Attributes.Time)
				cells = append(cells, []interface{}{
					val.Offer.Attributes.ID,
					val.Attributes.Bid,
					v.Attributes.ID,
					sv.Attributes.I,
					insOdds[0],  //Win
					insOdds[1],  //Lose
					insOdds[2],  //Draw
					oddsTime,    //Time
					now,         //Last Table Update
					insOdds[3],  //Margin
					insOdds[4],  //actual win
					insOdds[5],  //Actual lose
					insOdds[6]}) //Actual Draw
			}
		}
	}
	return filename, cells
}
