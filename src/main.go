package main

import (
	"CityPowerLoadShedding/src/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var stage string
var selectedBlock string
var loc *time.Location

func main() {
	// The block for your suburb can be found on city powers site
	selectedBlock = os.Getenv("SUBBLOCK")
	port := os.Getenv("PORT")
	tz := os.Getenv("TZ")
	loc, _ = time.LoadLocation(tz)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/get-schedule", getLoadSheddingSchedule)

	var err error
	go func() {
		for {
			stage, err = getCurrentStage()
			fmt.Printf("%s %s\n", "Stage: ", stage)
			if err != nil && stage == "" {
				fmt.Println("Error getting current stage")
				fmt.Println(err)
				continue
			}
			time.Sleep(10 * time.Minute)
		}
	}()

	// serve the app
	fmt.Printf("%s: %s\n", "Server started at port", port)
	port = ":" + port
	log.Fatal(http.ListenAndServe(port, router))
}

func getLoadSheddingSchedule(w http.ResponseWriter, r *http.Request) {
	var schedule types.Results
	var err error

	if stage != "" {
		schedule, err = fetchSchedule(stage, selectedBlock)
		if err != nil {
			fmt.Println("Error getting schedule from city power")
			fmt.Println(err)
		}

		var finalSchedule types.FinalSchedule
		finalSchedule.LoadSheddingTimes = getFinalSchedule(schedule, selectedBlock, loc)
		finalSchedule.CurrentStage = stage

		err = json.NewEncoder(w).Encode(finalSchedule)
		if err != nil {
			return
		}
	}
	return
}

// getFinalSchedule iterates through the data returned by fetchSchedule and processes the data to only return the time slots
// that load shedding will occur
func getFinalSchedule(schedule types.Results, selectedBlock string, loc *time.Location) []types.LoadSheddingTimes {

	currentTime := time.Now().In(loc)

	StageTimesJson := `[{"StartTime":"2022-07-13T22:00:00Z","EndTime":"2022-07-13T03:00:00Z","Stage":2},{"StartTime":"2022-07-13T03:00:00Z","EndTime":"2022-07-13T14:00:00Z","Stage":3},{"StartTime":"2022-07-13T14:00:00Z","EndTime":"2022-07-14T22:00:00Z","Stage":4},{"StartTime":"2022-07-14T22:00:00Z","EndTime":"2022-07-14T03:00:00Z","Stage":2},{"StartTime":"2022-07-14T03:00:00Z","EndTime":"2022-07-14T14:00:00Z","Stage":3},{"StartTime":"2022-07-14T14:00:00Z","EndTime":"2022-07-15T22:00:00Z","Stage":4},{"StartTime":"2022-07-15T22:00:00Z","EndTime":"2022-07-15T03:00:00Z","Stage":2},{"StartTime":"2022-07-15T03:00:00Z","EndTime":"2022-07-16T22:00:00Z","Stage":3}]`
	var stageTimes []types.StageTimes
	json.Unmarshal([]byte(StageTimesJson), &stageTimes)

	var loadShedToday []types.LoadSheddingTimes

	for k := range stageTimes {
		for _, result := range schedule {
			if isBlockMatch(result.SubBlock, selectedBlock) && (result.StartDateQuery.In(loc).Day() == currentTime.Day() || result.StartDateQuery.In(loc).Day() == time.Now().AddDate(0, 0, 1).In(loc).Day()) && result.StartDateQuery.In(loc).Month() == currentTime.Month() {
//				if (((result.StartDateQuery.In(loc).Equal(stageTimes[k].StartTime.In(loc)) || result.StartDateQuery.In(loc).After(stageTimes[k].StartTime.In(loc))) && (result.EndDateQuery.In(loc).Equal(stageTimes[k].EndTime.In(loc)) || result.EndDateQuery.In(loc).Before(stageTimes[k].EndTime.In(loc)))) && result.StageId == stageTimes[k].Stage) {
				if (result.StartDateQuery.In(loc).After(stageTimes[k].StartTime.In(loc)) && result.EndDateQuery.In(loc).Before(stageTimes[k].EndTime.In(loc)) && result.StageId == stageTimes[k].Stage) {
					var loadShedTimes types.LoadSheddingTimes
					loadShedTimes.StartTime = result.StartDateQuery.In(loc)
					loadShedTimes.EndTime = result.EndDateQuery.In(loc)
					loadShedTimes.Stage = result.StageId
					loadShedToday = append(loadShedToday, loadShedTimes)
				}
			}
	 	}

	}
/*
	for _, result := range schedule {
		if isBlockMatch(result.SubBlock, selectedBlock) && (result.StartDateQuery.In(loc).Day() == currentTime.Day() || result.StartDateQuery.In(loc).Day() == time.Now().AddDate(0, 0, 1).In(loc).Day()) && result.StartDateQuery.In(loc).Month() == currentTime.Month() {
			var loadShedTimes types.LoadSheddingTimes
			loadShedTimes.StartTime = result.StartDateQuery.In(loc)
			loadShedTimes.EndTime = result.EndDateQuery.In(loc)
			loadShedTimes.Stage = result.StageId
			loadShedToday = append(loadShedToday, loadShedTimes)
		}
	}
*/
	if len(loadShedToday) > 0 {
		return loadShedToday
	}

	return []types.LoadSheddingTimes{}
}

// fetchSchedule fetches the raw data from city power and parses it into the structs
func fetchSchedule(stage string, selectedBlock string) (types.Results, error) {
	url := "https://www.citypower.co.za/_api/web/lists/getByTitle('Loadshedding')/items?$select=*&$filter=substringof(%27" + selectedBlock + "%27,%20SubBlock)&$top=1000"
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("accept", "application/json;odata=verbose")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var jsonResponse types.JsonResponse
	err = json.Unmarshal([]byte(body), &jsonResponse)
	if err != nil {
		return types.Results{}, err
	}
	return jsonResponse.D.Results, nil
}

// getCurrentStage get the current load shedding stage from Eskom
func getCurrentStage() (string, error) {
	url := "https://loadshedding.eskom.co.za/LoadShedding/GetStatus"
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("accept", "application/json;odata=verbose")
	response, err := client.Do(request)
	if err != nil {
		return "", err // TODO: Improve error handling when Eskom fails
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	switch string(body) {
	case "1":
		return "0", nil
	case "2":
		return "1", nil
	case "3":
		return "2", nil
	case "4":
		return "3", nil
	case "5":
		return "4", nil
	case "6":
		return "5", nil
	case "7":
		return "6", nil
	default:
		return "", nil
	}
}

// isBlockMatch checks if the suburb block matches the one you need
func isBlockMatch(subBlock string, selectedBlock string) bool {
	split := strings.Split(subBlock, ";")
	for _, s := range split {
		if s == selectedBlock {
			return true
		}
	}
	return false
}
