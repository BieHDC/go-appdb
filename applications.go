package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "biehdc.webapp.applister/application"
	"biehdc.webapp.applister/paginate"
)

type Applications struct {
	Title string

	NumApps uint //only needed to display num Apps without re-len()-ing
	Apps    []*Application

	//Needed for searches for loading and storing
	Cons Constraints

	//Just for the stat
	Searchtime time.Duration

	//Pagination
	PG paginate.Paginate
}

type Constraints struct {
	Name       string
	Useability uint
	RosVersion string
	Tags       []string
}

func (apps *Applications) ConsToString() string {
	var final string
	for _, tag := range apps.Cons.Tags {
		final += tag + ", "
	}
	if final != "" {
		final = final[:len(final)-2]
	}
	return final
}

func NewApplications(title string) *Applications {
	return &Applications{Title: title}
}

func (apps *Applications) FindEntry(name string) *Application {
	for _, app := range apps.Apps {
		if app.Name == name {
			return app
		}
	}

	return nil
}

func (apps *Applications) GetSlice(fulllist *Applications) {
	low := (apps.PG.NextPage - 1) * apps.PG.EntriesPerPage
	if low < 0 {
		low = 0
	}

	high := int(low + apps.PG.EntriesPerPage)
	if high < 0 || high > len(fulllist.Apps) {
		high = len(fulllist.Apps)
	}

	apps.Apps = fulllist.Apps[low:high]

	apps.NumApps = uint(len(apps.Apps))
}

func NewApplicationsFromPath(path string, title string) (*Applications, error) {
	apps := NewApplications(title)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var errfin error

	for _, direntry := range files {
		if direntry.IsDir() {
			//ignore
			continue
		}
		filename := direntry.Name()
		if !strings.HasPrefix(filename, "app_") && !strings.HasSuffix(filename, ".csv") {
			//log.Printf("skipping: %s\n", filename)
			//err = errors.Join(err, fmt.Errorf("skipping: %s\n", filename)) //not fatal
			continue
		}

		csvfile, err := os.Open(filepath.Join(path, filename))
		if err != nil {
			//log.Fatalf("failed to read file: %s reason: %w\n", filename, err)
			errfin = errors.Join(errfin, err)
			continue
		}

		app, err := ApplicationFromCSV(csvfile)
		if err != nil {
			//log.Fatalf("%s failed: %s\n", filename, err)
			errfin = errors.Join(errfin, fmt.Errorf("%s failed:", filename), err)
			continue
		}

		csvfile.Close()

		apps.Apps = append(apps.Apps, app)
	}

	apps.NumApps = uint(len(apps.Apps))

	return apps, errfin
}
