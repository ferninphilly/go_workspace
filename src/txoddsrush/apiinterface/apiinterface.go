//Package apiinterface is the connectiont to the txOdds api
//It reads from the myconfig file and utilizes that data to build query strings
package apiinterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	mc "txoddsrush/myconfig"

	xj "github.com/basgys/goxml2json"
)

//So if we are going to get this data we are going to have to start by getting the matches...so let's get the matches.

/*CreateOdds is the struct that currently holds the data.
TODO on this is to break this apart. Right now the data won't unmarshal into
the smaller sections. I am keeping the structs around below anyway because I want to
come back to this rather annoying problem.
This is also the struct that can be used to unmarshal match data. */
type CreateOdds struct {
	Attributes struct {
		Time      string `json:"time"`
		Timestamp string `json:"timestamp"`
	} `json:"@attributes"`
	Match []struct {
		Attributes struct {
			ID   string `json:"id"`
			Xsid string `json:"xsid"`
		} `json:"@attributes"`
		Time    string `json:"time"`
		Group   string `json:"group"`
		Hteam   string `json:"hteam"`
		Ateam   string `json:"ateam"`
		Results struct {
			Num0 string `json:"0"`
		} `json:"results"`
		Bookmaker []struct {
			Attributes struct {
				Bid  string `json:"bid"`
				Name string `json:"name"`
			} `json:"@attributes"`
			Offer struct {
				Attributes struct {
					ID          string `json:"id"`
					N           string `json:"n"`
					Ot          string `json:"ot"`
					Otname      string `json:"otname"`
					LastUpdated string `json:"last_updated"`
					Flags       string `json:"flags"`
					Bmoid       string `json:"bmoid"`
				} `json:"@attributes"`
				Odds []struct {
					Attributes struct {
						I            string `json:"i"`
						Time         string `json:"time"`
						StartingTime string `json:"starting_time"`
					} `json:"@attributes"`
					O1 string `json:"o1"`
					O2 string `json:"o2"`
					O3 string `json:"o3"`
				} `json:"odds"`
			} `json:"offer"`
		} `json:"bookmaker"`
	} `json:"match"`
}

//Bookmaker contains Bookmaker->Offer->Odds
type Bookmaker []struct {
	Attributes struct {
		Bid  string `json:"bid"`
		Name string `json:"name"`
	} `json:"@attributes"`
	Offer
}

//Offer contains Offer->Odds
type Offer struct {
	Attributes struct {
		ID          string `json:"id"`
		N           string `json:"n"`
		Ot          string `json:"ot"`
		Otname      string `json:"otname"`
		LastUpdated string `json:"last_updated"`
		Flags       string `json:"flags"`
		Bmoid       string `json:"bmoid"`
	} `json:"@attributes"`
	Odds
}

//Mainstruct is the highest level struct for unmarshalling data
type Mainstruct struct {
	Attributes struct {
		Time      string
		Timestamp string
	}
	Match
}

//Match contains Bookmaker->Offer->Odds
type Match []struct {
	Attributes struct {
		ID   string `json:"id"`
		Xsid string `json:"xsid"`
	} `json:"@attributes"`
	Time    string `json:"time"`
	Group   string `json:"group"`
	Hteam   string `json:"hteam"`
	Ateam   string `json:"ateam"`
	Results struct {
		Num0 string `json:"0"`
	} `json:"results"`
	Bookmaker
}

//feedOdds The Order here is:
//feedOdds->Match->Bookmaker->Offer->Odds
type feedOdds struct {
	Attributes struct {
		Time      string `json:"time"`
		Timestamp string `json:"timestamp"`
	}
	Match
}

//Teams is the struct I am using to capture all teams data
type Teams struct {
	Competitors struct {
		Competitor []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Group     string `json:"group"`
			Countryid string `json:"countryid"`
			Sportid   string `json:"sportid"`
		} `json:"competitor"`
	} `json:"competitors"`
}

//Odds is the lowest level
type Odds []struct {
	Attributes struct {
		I            string `json:"i"`
		Time         string `json:"time"`
		StartingTime string `json:"starting_time"`
	} `json:"@attributes"`
	O1 string `json:"o1"`
	O2 string `json:"o2"`
	O3 string `json:"o3"`
}

//xmlToJSON exists because not everything in the api responds to "&json=1" so sometimes we
//need to convert from xml to json. This package does that really, really well.
func xmlToJSON(contents string) string {
	xml := strings.NewReader(contents) //Input needs to be a reader
	json, err := xj.Convert(xml)
	mc.HandleError(err)
	return json.String()
}

//URLEncode is basically to encode the json from config
func urlEncode(data map[string]string) string {
	var buf bytes.Buffer
	for k, v := range data {
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
		buf.WriteByte('&')
	}
	eurl := buf.String()
	fmt.Println(eurl[0 : len(eurl)-1])
	return eurl[0 : len(eurl)-1]
}

//createURL basically creates the URL from strings
func createURL() string {
	cfg := mc.ReturnConfig()
	bodstring := fmt.Sprintf("ident=%s&passwd=%s&", cfg.APIData.Username, cfg.APIData.Password)
	extras := urlEncode(cfg.APIData.URLOptions)
	var buf bytes.Buffer
	buf.WriteString(cfg.APIData.BaseURL)
	buf.WriteString(cfg.APIData.FeedType)
	buf.WriteString(cfg.APIData.EndURL)
	buf.WriteString(bodstring)
	buf.WriteString(extras)
	fmt.Println(buf.String())
	return buf.String()
}

func callAPI(url string) []byte {
	response, err := http.Get(url)
	mc.HandleError(err)
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	mc.HandleError(err)
	return contents //Return a string to make troubleshooting easier- string(contents)
}

//createFeedOdds is how we unmarshal data. This will be the struct to return
func (fo *CreateOdds) createFeedOdds(contents []byte) {
	err := json.Unmarshal([]byte(contents), &fo)
	mc.HandleError(err)
}

//CreateTeamList creates our list of British soccer teams
//There is no option to get JSON directly so we need to convert from XML
func (tn *Teams) createTeamList(contents []byte) {
	updjson := xmlToJSON(string(contents))
	error := json.Unmarshal([]byte(updjson), &tn)
	mc.HandleError(error)
}

//ReturnFeedOdds is the return function for the odds calculations.
func (fo *CreateOdds) ReturnFeedOdds() {
	url := createURL()
	contents := callAPI(url)
	fo.createFeedOdds(contents)
}

//ReturnTeamList is the return function for team list depending on what's in myconfig
func ReturnTeamList() Teams {
	var tn Teams
	url := createURL()
	contents := callAPI(url)
	tn.createTeamList(contents)
	return tn
}

//ReturnBookies is how to get the Unique Bookies
func ReturnBookies() map[string]string {
	var fo CreateOdds
	fo.ReturnFeedOdds()
	ub := make(map[string]string) //for "unique bookies"
	for _, v := range fo.Match {
		for _, val := range v.Bookmaker {
			if _, ok := ub[val.Attributes.Bid]; !ok {
				ub[val.Attributes.Bid] = val.Attributes.Name
			}
		}
	}
	return ub
}
