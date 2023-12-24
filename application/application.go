package application

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type KeyVal struct {
	Key string
	Val string
}

type Application struct {
	Name        string   `json:"appname"`
	Useability  uint     `json:"useability"`
	Version     string   `json:"version"`
	RosVersion  string   `json:"rosversion"`
	Tags        []string `json:"tags"`
	Workarounds string   `json:"-"` //short workaround (eg copy ogwindows.dll somewhere)
	//
	ProgramDetails string   `json:"-"`
	KnownIssues    string   `json:"-"` //long workaround/instructions
	Downloads      []KeyVal `json:"-"` //key: desc, value:link
	MainScreenshot KeyVal   `json:"-"` //the one displayed
	Screenshots    []KeyVal `json:"-"` //the others. key: desc, value:filename
}

// Use strings.NewReader to make an io.Reader
func ApplicationFromCSV(source io.Reader) (*Application, error) {
	var app Application
	var errfin error
	csvreader := csv.NewReader(source)
	csvreader.Comma = ';'
	csvreader.FieldsPerRecord = -1 //set to variable fields
	csvreader.ReuseRecord = true   //we copy everything anyway
	for {
		record, err := csvreader.Read()
		if err == io.EOF {
			break //done
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse record: %w", err)
		}
		if len(record) < 2 {
			errfin = errors.Join(errfin, fmt.Errorf("line too short: %s", record[0])) //we must have at least 1 key and 1 value for all types at this point
			continue
		}

		switch record[0] {
		case "Name":
			app.Name = record[1]
		case "Version":
			app.Version = record[1]
		case "Useability":
			useability, err := strconv.ParseUint(record[1], 10, 0)
			if err != nil {
				errfin = errors.Join(errfin, fmt.Errorf("error converting number: %s reason: %w", record[1], err))
				continue
			}
			app.Useability = uint(useability)
		case "RosVersion":
			app.RosVersion = record[1]
		case "Tags":
			if record[1] != "" { //If s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.
				for _, tag := range strings.Split(record[1], ",") {
					app.Tags = append(app.Tags, tag)
				}
			}
		case "Workarounds":
			app.Workarounds = record[1]
		case "ProgramDetails":
			app.ProgramDetails = record[1]
		case "KnownIssues":
			app.KnownIssues = record[1]
		case "Downloads":
			if len(record) < 3 {
				errfin = errors.Join(errfin, fmt.Errorf("Downloads entry too short"))
				continue
			}
			app.Downloads = append(app.Downloads, KeyVal{record[1], record[2]})
		case "Screenshots":
			if len(record) < 3 {
				errfin = errors.Join(errfin, fmt.Errorf("Screenshots entry too short"))
				continue
			}
			if app.MainScreenshot.Val == "" {
				app.MainScreenshot.Key = record[1]
				app.MainScreenshot.Val = record[2]
				continue
			}
			app.Screenshots = append(app.Screenshots, KeyVal{record[1], record[2]})
		default:
			errfin = errors.Join(errfin, fmt.Errorf("unknown field: %s", record[0]))
			continue
		}
	}

	//pre validation
	if errfin != nil {
		return nil, errfin
	}

	//after all records, do some validation and add it to the list
	if app.Name == "" {
		errfin = errors.Join(errfin, fmt.Errorf("missing Name"))
	}
	if app.Version == "" {
		errfin = errors.Join(errfin, fmt.Errorf("missing Version"))
	}
	if app.Useability == 0 {
		errfin = errors.Join(errfin, fmt.Errorf("missing Useability"))
	}
	if app.RosVersion == "" {
		errfin = errors.Join(errfin, fmt.Errorf("missing RosVersion"))
	}
	if len(app.Tags) == 0 {
		errfin = errors.Join(errfin, fmt.Errorf("missing Tags"))
	}

	//explicitly return a nil-app when there was an error
	if errfin != nil {
		return nil, errfin
	}

	return &app, nil
}
