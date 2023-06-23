package IRacingData

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
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
	fmt.Println("iRacing INIT")
	session, err := auth(email, hash)

	if err != nil {
		log.Fatal(err)
	}
	irsession = session
}
