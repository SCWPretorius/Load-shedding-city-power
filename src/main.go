package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type JsonResponse struct {
	D D
}

type D struct {
	Results Results
}

type Results []struct {
	Metadata struct {
		Id   string `json:"id"`
		Uri  string `json:"uri"`
		Etag string `json:"etag"`
		Type string `json:"type"`
	} `json:"__metadata"`
	FirstUniqueAncestorSecurableObject struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"FirstUniqueAncestorSecurableObject"`
	RoleAssignments struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"RoleAssignments"`
	AttachmentFiles struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"AttachmentFiles"`
	ContentType struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"ContentType"`
	GetDlpPolicyTip struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"GetDlpPolicyTip"`
	FieldValuesAsHtml struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"FieldValuesAsHtml"`
	FieldValuesAsText struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"FieldValuesAsText"`
	FieldValuesForEdit struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"FieldValuesForEdit"`
	File struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"File"`
	Folder struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"Folder"`
	ParentList struct {
		Deferred struct {
			Uri string `json:"uri"`
		} `json:"__deferred"`
	} `json:"ParentList"`
	FileSystemObjectType int       `json:"FileSystemObjectType"`
	Id                   int       `json:"Id"`
	ContentTypeId        string    `json:"ContentTypeId"`
	Title                string    `json:"Title"`
	Location             string    `json:"Location"`
	EventDate            time.Time `json:"EventDate"`
	EndDate              time.Time `json:"EndDate"`
	Description          string    `json:"Description"`
	FAllDayEvent         bool      `json:"fAllDayEvent"`
	FRecurrence          bool      `json:"fRecurrence"`
	ParticipantsPickerId string    `json:"ParticipantsPickerId"`
	Category             string    `json:"Category"`
	FreeBusy             string    `json:"FreeBusy"`
	Overbook             string    `json:"Overbook"`
	SubBlock             string    `json:"SubBlock"`
	Reason               string    `json:"Reason"`
	StageId              int       `json:"StageId"`
	StartDateQuery       time.Time `json:"StartDateQuery"`
	EndDateQuery         time.Time `json:"EndDateQuery"`
	LoadShed             int       `json:"LoadShed"`
	ID                   int       `json:"ID"`
	Modified             time.Time `json:"Modified"`
	Created              time.Time `json:"Created"`
	AuthorId             int       `json:"AuthorId"`
	EditorId             int       `json:"EditorId"`
	ODataUIVersionString string    `json:"OData__UIVersionString"`
	Attachments          bool      `json:"Attachments"`
	GUID                 string    `json:"GUID"`
}

type LoadShedTimes struct {
	StartTime time.Time
	EndTime   time.Time
}

func main() {
	// The block for your suburb can be found on city powers site
	selectedBlock := os.Getenv("SUBBLOCK")
	tz := os.Getenv("TZ")
	loc, _ := time.LoadLocation(tz)

	var err error
	var stage string
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

	var schedule Results
	go func() {
		for {
			if stage != "" {
				schedule, err = fetchSchedule(stage, selectedBlock)
				if err != nil {
					fmt.Println("Error getting schedule from city power")
					fmt.Println(err)
					continue
				}
			}
			time.Sleep(12 * time.Hour) // TODO: Change to schedule for specific times
		}
	}()

	finalSchedule := getFinalSchedule(schedule, selectedBlock, loc)

	// TODO: Integrate with HA addon
	fmt.Println(finalSchedule)
}

// getFinalSchedule iterates through the data returned by fetchSchedule and processes the data to only return the time slots
// that load shedding will occur
func getFinalSchedule(schedule Results, selectedBlock string, loc *time.Location) map[int]LoadShedTimes {

	currentTime := time.Now().In(loc)

	loadShedToday := make(map[int]LoadShedTimes)
	for idx, result := range schedule {
		var loadShedTimes LoadShedTimes

		if isBlockMatch(result.SubBlock, selectedBlock) && result.StartDateQuery.In(loc).Day() == currentTime.Day() && result.StartDateQuery.In(loc).Month() == currentTime.Month() {
			loadShedTimes.StartTime = result.StartDateQuery.In(loc)
			loadShedTimes.EndTime = result.EndDateQuery.In(loc)
			loadShedToday[idx] = loadShedTimes
		}
	}
	return loadShedToday
}

// fetchSchedule fetches the raw data from city power and parses it into the structs
func fetchSchedule(stage string, selectedBlock string) (Results, error) {
	url := "https://www.citypower.co.za/_api/web/lists/getByTitle('Loadshedding')/items?$select=*&$filter=Title%20eq%20%27Stage" + stage + "%27%20and%20substringof(%27" + selectedBlock + "%27,%20SubBlock)&$top=1000"
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

	var jsonResponse JsonResponse
	err = json.Unmarshal([]byte(body), &jsonResponse)
	if err != nil {
		return Results{}, err
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
		return "", err
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
