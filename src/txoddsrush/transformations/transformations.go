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
)

//Teams is the struct to contain our teams

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
