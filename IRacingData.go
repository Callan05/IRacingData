package IRacingData

// V1.1.0

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
)

var irapi string = "https://members-ng.iracing.com/"
var irsession *http.Client

type leagueResponseLink struct {
	Link    string `json:"link"`
	Expires string `json:"expires"`
}

type Track struct {
	Config_name string
	Track_name  string
	Track_id    int
}

type TrackAsset struct {
	Large_image string `json:"large_image"`
	Folder      string `json:"folder"`
}
type Session struct {
	Launch_at          string
	Race_length        int
	Status             int
	Track              Track
	Time_limit         int
	Private_session_id int
	Has_results        bool
	Session_id         int
	Results            RaceResult
	Subsession_id      int
}
type LeagueSeason struct {
	Sessions []Session
}

type Result struct {
	Display_name    string
	Time            string
	Finish_position int
	Interval        int64
	Car_name        string
	League_points   int
	Laps_complete   int
	Average_lap     int64
}

type SessionResult struct {
	Results []Result
}

type RaceResult struct {
	Subsession_id       int
	League_season_name  string
	Start_time          string
	End_time            string
	Session_results     []SessionResult
	Event_laps_complete int
	Track               Track
	League_season_id    int
	RaceID              int
}

type League_Season_Standings struct {
	Car_class_id int
	Success      bool
	Season_id    int
	Car_id       int
	Standings    Standings
	League_id    int
}

type Standings struct {
	Driver_standings []Driver_Standings
	Team_standings   []Team_Standings
}
type Driver_Standings struct {
	Rownum               int
	Position             int
	Driver               Driver
	Car_number           int
	Driver_nickname      string
	Wins                 int
	Average_start        int
	Average_finish       int
	Base_points          int
	Negative_adjustments int
	Positive_adjustments int
	Total_adjustments    int
	Total_points         int
}
type Driver struct {
	Cust_id      int
	Display_name string
}
type Team_Standings struct {
}

func auth(email string, hash string) (*http.Client, error) {
	var jar, err = cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	var url string = irapi + "auth"

	method := "POST"

	payload := strings.NewReader(`{
	  "email": "` + email + `",
	  "password": "` + hash + `"
  }`)

	client := &http.Client{
		Jar: jar,
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return client, nil
}
func Init(email string, hash string) {
	fmt.Println("iRacing Init() called")
	session, err := auth(email, hash)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("iRacing Session Started")
	irsession = session
}

func followLink(url string) (map[string]any, error) {

	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := irsession.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var jsondata map[string]any
	json.Unmarshal([]byte(body), &jsondata)

	return (jsondata), nil

}
func GetLeague(leagueID string) (map[string]any, error) {

	url := irapi + "data/league/get?league_id=" + leagueID
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := irsession.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var link leagueResponseLink
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(body, &link); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	data, err := followLink(link.Link)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	ret := data
	return ret, nil
}

func GetLeagueSessions(leagueID string, seasonID string, resultsOnly bool) ([]Session, error) {

	var url string = irapi + "data/league/season_sessions?league_id=" + leagueID + "&season_id=" + seasonID
	if resultsOnly {
		url += "&results_only=true"
	}
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := irsession.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var link leagueResponseLink
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(body, &link); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	req2, err := http.NewRequest(method, link.Link, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res2, err := irsession.Do(req2)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res2.Body.Close()

	body2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var season LeagueSeason
	json.Unmarshal([]byte(body2), &season)

	return (season.Sessions), nil

}

func GetLeagueSeasonStandings(leagueID string, seasonID string, resultsOnly bool) (League_Season_Standings, error) {
	var season League_Season_Standings
	var url string = irapi + "data/league/season_standings?league_id=" + leagueID + "&season_id=" + seasonID
	if resultsOnly {
		url += "&results_only=true"
	}
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return season, err
	}

	res, err := irsession.Do(req)
	if err != nil {
		fmt.Println(err)
		return season, err
	}
	defer res.Body.Close()

	var link leagueResponseLink
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return season, err
	}

	if err := json.Unmarshal(body, &link); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	req2, err := http.NewRequest(method, link.Link, nil)

	if err != nil {
		fmt.Println(err)
		return season, err
	}

	res2, err := irsession.Do(req2)
	if err != nil {
		fmt.Println(err)
		return season, err
	}
	defer res2.Body.Close()

	body2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		fmt.Println(err)
		return season, err
	}

	json.Unmarshal([]byte(body2), &season)

	return (season), nil

}

func GetLeagueSeasons(leagueID string) (map[string]any, error) {

	var url string = irapi + "data/league/seasons?league_id=" + leagueID
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	res, err := irsession.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	var link leagueResponseLink
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := json.Unmarshal(body, &link); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	data, err := followLink(link.Link)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	ret := data
	return ret, nil
}

func GetSession(sessionId int) (RaceResult, error) {
	var ret RaceResult
	var url string = irapi + "data/results/get?include_licenses=false&subsession_id=" + strconv.Itoa(sessionId)
	fmt.Println(url)
	method := "GET"
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ret, err
	}

	res, err := irsession.Do(req)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}
	defer res.Body.Close()

	var link leagueResponseLink
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}

	if err := json.Unmarshal(body, &link); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	data, err := followLink(link.Link)
	if err != nil {
		fmt.Println(err)
		return ret, err
	}

	jsonString, _ := json.Marshal(data)
	var s RaceResult
	json.Unmarshal(jsonString, &s)

	ret = s
	return ret, nil
}
