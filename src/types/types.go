package types

import "time"

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

type LoadSheddingTimes struct {
	StartTime time.Time
	EndTime   time.Time
	Stage     int
}

type StageTimes struct {
	StartTime time.Time
	EndTime   time.Time
	Stage     int
}

type FinalSchedule struct {
	LoadSheddingTimes []LoadSheddingTimes
	CurrentStage      string
}
