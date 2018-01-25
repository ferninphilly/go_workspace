//Package dbconnection establishes the connection to the database and reads SQL queries
//from the files in the "sqlQueries" sub directories.
package dbconnection

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	ai "txoddsrush/apiinterface"
	mc "txoddsrush/myconfig"
	tx "txoddsrush/transformations"

	_ "github.com/lib/pq" //This is for the drivers for postgresql
)

//executeSQL is an abstracted method to execute any insert SQL query.
//It uses the variadiac interface function to take a variable number of args from anywhere
//Ultimately this is kind of the core DB connection function
func executeSQL(cfg mc.Config, filename string, queryType string, vals [][]interface{}) string {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBData.Host, cfg.DBData.Port, cfg.DBData.DBUser, cfg.DBData.DBPassword,
		cfg.DBData.DBName)
	sqlRead, err := ioutil.ReadFile(filename)
	mc.HandleError(err)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		log.Print(err)
		return "There was an issue with the connection. That issue is above"
	}
	for _, v := range vals {
		_, execerr := db.Exec(string(sqlRead), v...)
		if execerr != nil {
			log.Print(execerr)
			return "There was an issue with the insert. That issue is above"
		}
	}

	defer db.Close()
	return "Completed SQL query for " + queryType
}

//LoadMatchData is how we load data to the matches table.
func LoadMatchData(fo ai.CreateOdds, cfg mc.Config) {
	filename, pdata := tx.ReadyMatchData(fo)
	executeSQL(cfg, filename, "matches", pdata)
	totrows := strconv.Itoa(len(pdata))
	fmt.Println("Completed insert into MATCHES of " + totrows + " rows.")
}

//LoadOddsData is how we load Odds data into the Odds table
func LoadOddsData(fo ai.CreateOdds, cfg mc.Config) {
	filename, pdata := tx.ReadyOddsData(fo)
	executeSQL(cfg, filename, "odds", pdata)
	totrows := strconv.Itoa(len(pdata))
	fmt.Println("Completed insert into ODDS of " + totrows + " rows.")

}
