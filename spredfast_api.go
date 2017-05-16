package main

import (
        "fmt"
        "io/ioutil"
        "gopkg.in/yaml.v2"
//        "log"
        "net/http"
        "strings"
        "encoding/json"
        "os"
        "encoding/csv"
)

var resp *http.Response

type highlevel struct {
    SfEntityType string
    Reports []reports
    }

type reports struct {
    Name string
    Title string
    Description string
}

type yamld struct {
    Username string
    Password string
    Client_id string
    Secret string
    Base_url string
    Production_url string
    Company_url string
    Version string
    Report string
}

type keyd struct {
    SfEntityType string
    AccessToken string
    ExpiryTime int
    Status []string
    }

type companyData struct {
    SfEntityType string
    Email string
    Companies []companyID
    }

type companyID struct {
    SfEntityType string
    DirectoryId string
    Id string
    Name string
    Environment string
    EnvironmentUrl string
}

type apiData struct {
    Id string
    Environment string
    Version string
    Production_url string
}

type reportID struct {
    EntityType string
    Id string
}

/*func getParameter() string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
    text, _ := reader.Readstring('\n')
}*/

func (y *yamld) getYaml() *yamld {
    yamlFile, err := ioutil.ReadFile("spred.yaml")
    errorHandler("yamlFile.Get err", err)
    err = yaml.Unmarshal(yamlFile, y)
    errorHandler("Unmarshal Error", err)
    return y
}

func makeAPIcall(url, apikey, calltype, body, reason string) map[string]*json.RawMessage{
    var objmap map[string]*json.RawMessage
    switch calltype {
        case "POST":
//            fmt.Println("I got a Post!")
            resp, err := http.Post(url, "", strings.NewReader(body))
            errorHandler("Post Error ", err)
            defer resp.Body.Close()
            if resp.StatusCode == 200 {
                fmt.Println(reason)
                bodyBytes, err2 := ioutil.ReadAll(resp.Body)
            errorHandler("Read Error Post ", err2)
            jsondata := []byte(bodyBytes)
            json.Unmarshal(jsondata, &objmap)
            }
        case "GET":
//            fmt.Println("I got a Get")
            req, err := http.NewRequest("GET", url, nil)
            errorHandler("Request Error ", err)
           req.Header.Set("Accept", "application/json")
           req.Header.Set("Authorization", "Bearer " + apikey)
           resp, err := http.DefaultClient.Do(req)
           errorHandler("http error ", err)
           defer resp.Body.Close()
           if resp.StatusCode == 200 {
               fmt.Println(reason)
               bodyBytes, err2 := ioutil.ReadAll(resp.Body)
               errorHandler("ioutil Readall Error", err2)
               jsondata := []byte(bodyBytes)
               json.Unmarshal(jsondata, &objmap)
                }
        }   
    return objmap
}

func getAPIKey(y yamld) string{
    body := "email=" + y.Username +
            "&password=" + y.Password +
            "&client_id=" + y.Client_id +
            "&client_secret=" + y.Secret
    api_json := makeAPIcall(y.Base_url, "", "POST", body, "Getting Api Key")
    var d keyd  
    err := json.Unmarshal(*api_json["data"], &d)
    errorHandler("JSON Unmarshal Error ", err)
    return d.AccessToken 
}

func pyformat(format string, args ...string)  string{
    r := strings.NewReplacer(args...)
    return r.Replace(format)
}

func (id_env *apiData) getCompanyID(y yamld, apikey string) *apiData {
    company_url := pyformat(y.Company_url, "{version}", y.Version)
    comdata := makeAPIcall(company_url, apikey, "GET", "nothing", "Getting Company Data")
    var cd companyData
    err := json.Unmarshal(*comdata["data"], &cd)
    errorHandler("Unmarshal Error in getCompanyID ", err)
    id_env.Id = cd.Companies[0].Id
    id_env.Environment = cd.Companies[0].Environment
    id_env.Version = y.Version
    id_env.Production_url = y.Production_url 
    return  id_env
}

func write_to_file(responseString string, report string) {
    writtendata := strings.Split(responseString, "\n")
    file,err := os.Create(report + ".csv")
    errorHandler("Error creating file ", err)
    defer file.Close()
    writer := csv.NewWriter(file)
    defer writer.Flush()
    writer.UseCRLF = true
    for _, value := range writtendata {
        err2 := writer.Write(strings.Split(value, ","))
        errorHandler("Write error on ", err2)
    }
    
}

func getListOfReport(y yamld, baseurl string, apikey string) []string {
    report_js := makeAPIcall(baseurl, apikey, "GET", "nothing",
                             "Getting list of reports")
    var hl highlevel
    err4 := json.Unmarshal(*report_js["data"], &hl)
    errorHandler("json Unmarshal @ report list", err4)
    var rpts_lst []string
    for _, value := range hl.Reports {
        rpts_lst = append(rpts_lst, value.Name)
    }
    return rpts_lst
    }

func main() {
    var y yamld
    y.getYaml()
    apikey := getAPIKey(y)
    var co apiData
    co.getCompanyID(y,apikey)
    baseurl := pyformat(y.Production_url, "{environment}", co.Environment,
            "{version}", y.Version, "{companyID}", co.Id)
    rpt := getListOfReport(y, baseurl, apikey)
    for _, value := range rpt {
            fmt.Println("Now I'm creating a report for " + string(value))
            baseurl := pyformat(y.Production_url, "{environment}", co.Environment,
            "{version}", y.Version, "{companyID}", co.Id)
            baseurl += "/" + string(value)
            getReportId := makeAPIcall(baseurl, apikey, "GET", "nothing",
                                       "Getting Report Name")
            var newrpt map[string]interface{}
            err := json.Unmarshal(*getReportId["data"], &newrpt)
            errorHandler("jsonUnMarshal error in main ", err)
            reporturl := baseurl + "/instance/" + newrpt["id"].(string)
            req, err := http.NewRequest("GET", reporturl, nil)
            errorHandler("http error in main ", err)
            req.Header.Set("Accept", "text/plain")
            req.Header.Set("Authorization", "Bearer " + apikey)
            resp, err := http.DefaultClient.Do(req)
            errorHandler("Error in main DefaultClient", err)
            responseData, err := ioutil.ReadAll(resp.Body)
            errorHandler("ioutil error in main ", err)
            if 
            write_to_file(string(responseData), string(value))
    }

}

func errorHandler(message string, err error) {
    if err != nil {
        fmt.Println(message, err)
    }
}



