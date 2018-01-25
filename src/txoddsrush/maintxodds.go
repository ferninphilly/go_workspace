package main

import (
	ai "txoddsrush/apiinterface"
	db "txoddsrush/dbconnection"
	mc "txoddsrush/myconfig"
)

func main() {
	var fo ai.CreateOdds
	fo.ReturnFeedOdds()
	var cfg = mc.ReturnConfig()
	db.LoadMatchData(fo, cfg)
	db.LoadOddsData(fo, cfg)
}
