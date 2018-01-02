//Package dbconnection establishes the connection to the database and reads SQL queries
//from the files in the "sqlQueries" sub directories.
package dbconnection

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strconv"
	mc "txoddsrush/myconfig"

	_ "github.com/lib/pq" //This is for the drivers for postgresql
)

//RunQuery is our connection to a database to run inserts or queries
//TODO: Right now everything that goes in is a string in vars ..string.
//I need to go back and do a case/switch to fix later
func RunQuery(cfg mc.Config, filename string, vars ...string) string {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBData.Host, cfg.DBData.Port, cfg.DBData.DBUser, cfg.DBData.DBPassword,
		cfg.DBData.DBName)
	args := make([]interface{}, len(vars))
	for i, v := range vars {
		args[i] = v
	}
	sqlRead, err := ioutil.ReadFile(filename)
	mc.HandleError(err)
	db, err := sql.Open("postgres", psqlConn)
	mc.HandleError(err)
	oddsid, err1 := strconv.Atoi(vars[0])
	mc.HandleError(err1)
	bookiekey, err2 := strconv.Atoi(vars[1])
	mc.HandleError(err2)
	matchkey, err3 := strconv.Atoi(vars[2])
	mc.HandleError(err3)
	orderofoffer, err4 := strconv.Atoi(vars[3])
	mc.HandleError(err4)
	hteam, err5 := strconv.ParseFloat(vars[4], 32)
	mc.HandleError(err5)
	ateam, err6 := strconv.ParseFloat(vars[5], 32)
	mc.HandleError(err6)
	draw, err7 := strconv.ParseFloat(vars[6], 32)
	mc.HandleError(err7)
	tableupdate, err8 := strconv.Atoi(vars[7])
	mc.HandleError(err8)
	lu, err9 := strconv.Atoi(vars[8])
	mc.HandleError(err9)
	_, execerr := db.Exec(string(sqlRead), oddsid, bookiekey, matchkey, orderofoffer,
		hteam, ateam, draw, tableupdate, lu) //args...)
	mc.HandleError(execerr)
	defer db.Close()
	return string(sqlRead)
}
