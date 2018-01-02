package main

import (
	"fmt"
	"strconv"
	"time"
	ai "txoddsrush/apiinterface"
	db "txoddsrush/dbconnection"
	mc "txoddsrush/myconfig"
)

func main() {
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

//a := db.DBConn(cfg)
//fmt.Println(a)

/*	url := createURL(&cfg)
	fmt.Println(url)
	var afo afo
	afo.createFeedOdds(url)
	for _, v := range afo.Match {
		if v.Hteam == "Tottenham" || v.Ateam == "Tottenham" {
			fmt.Println(v.Hteam + " vs " + v.Ateam + "\nTime: " + v.Time)
			for _, val := range v.Bookmaker {
				fmt.Println("According to: " + val.Attributes.Name + " the odds are:")
				fmt.Println("\t" + v.Hteam + " Odds to win are: " + val.Offer.Odds.O1)
				fmt.Println("\t" + v.Ateam + " Odds to win are: " + val.Offer.Odds.O2)
				fmt.Println("\tOdds of a draw are: " + val.Offer.Odds.O3 + "\n")
				fmt.Println("--------------------------------------")
			}
		}
	}
	//fmt.Println(fo.match[0].Hteam)
	/*
			for _, val := range v.Bookmaker {
		fmt.Println("According to " + v.Bookmaker.Attributes.Name + ":\n")
		fmt.Println("Odds of " + v.Hteam + " winning: " + val.O1 + "\n")
		fmt.Println("Odds of " + v.Ateam + " winning: " + val.O2 + "\n")
		fmt.Println("Odds of a draw: " + val.O3)
*/
//pretty.Print(mtc.Matches.Match[0].Bookmaker[0].Name)

/*	for _, v := range mtc.Matches.Match {
	if v.Hteam.Content == "Tottenham" || v.Ateam.Content == "Tottenham" {
		pretty.Print("Home Team: " + v.Hteam.Content + "\n" +
			"Away Team: " + v.Ateam.Content + " Time: " + v.Time.String() + "\n\t")

		for _, book := range v.Bookmaker {
			if book.Name != "" {
				pretty.Print("Bookmaker: " + book.Name + "\nFlags: " + book.Offer.Flags + "\n\t")
				for _, odds := range book.Offer.Odds {
					pretty.Print(v.Hteam.Content + "ID: " + v.Hteam.ID + " Wins: " + odds.O1 + "\n" +
						v.Ateam.Content + " Wins: " + odds.O3 + "\n" +
						"Draw: " + odds.O2 + "\n" + "I: " + odds.I + "\n")

				}
			}
		}
		pretty.Print("\n-------------------\n\n")
	}
} */
