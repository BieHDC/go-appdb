package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"

	. "biehdc.webapp.applister/application"
)

// Valid Get Parameters - All as string
//  appname
//  useability - if not a number, assume NONE
//  rosversion
//  tags

var urlParseFail = fmt.Errorf("Failed to parse url")
var emptySearch = fmt.Errorf("No search constraints present")

func parseSearchGetRequest(r *http.Request) (Constraints, error) {
	var constraints Constraints
	cliurl, err := url.Parse(r.URL.String())
	if err != nil {
		return Constraints{}, err
	}
	params := cliurl.Query()

	constraints.Name = params.Get("appname")
	useability_parsed, err := strconv.ParseUint(params.Get("useability"), 10, 0)
	if err != nil {
		useability_parsed = 0 //Is "None"
	}
	constraints.Useability = uint(useability_parsed)
	constraints.RosVersion = params.Get("rosversion")
	tags := params.Get("tags")
	if tags != "" { //If s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.
		for _, tag := range strings.Split(tags, ",") {
			constraints.Tags = append(constraints.Tags, strings.TrimSpace(tag))
		}
	}

	// If the user searched for nothing
	if constraints.Name == "" && constraints.Useability == 0 &&
		constraints.RosVersion == "" && len(constraints.Tags) == 0 {
		return Constraints{}, emptySearch
	}

	return constraints, nil
}

func (fulllist *Applications) GetMatches(cons Constraints) *Applications {
	searched := NewApplications(fulllist.Title)
	searched.Cons = cons                        // Needs to be reported back to the ui
	searched.Apps = make([]*Application, 0, 10) // 10 Results should be about accurate

	benched := time.Now()
	for _, app := range fulllist.Apps {
		if cons.Name != "" && !fuzzy.MatchFold(cons.Name, app.Name) {
			continue
		}

		if app.Useability < cons.Useability {
			continue
		}

		if cons.RosVersion != "" && !fuzzy.MatchFold(cons.RosVersion, app.RosVersion) {
			continue
		}

		if len(cons.Tags) > 0 {
			found := false // we succeed if one tag matches
			for _, tagcon := range cons.Tags {
				for _, tagog := range app.Tags {
					if tagcon == tagog {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			//if we found none, we fail
			if !found {
				continue
			}
		}

		searched.Apps = append(searched.Apps, app)
	}
	searched.NumApps = uint(len(searched.Apps))

	searched.Searchtime = time.Now().Sub(benched)
	return searched
}

func searchedpage(apps *Applications, tpl string) http.HandlerFunc {
	indextemplate := template.Must(template.ParseFiles(tpl))

	return func(w http.ResponseWriter, r *http.Request) {
		constraints, err := parseSearchGetRequest(r)
		if err != nil {
			if err == emptySearch {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			log.Println(err)
			http.Error(w, "Error while searching", http.StatusInternalServerError)
			return
		}

		buf := &bytes.Buffer{}
		err = indextemplate.Execute(buf, apps.GetMatches(constraints))
		if err != nil {
			log.Println(err)
			http.Error(w, "Error while searching", http.StatusInternalServerError)
			return
		}

		buf.WriteTo(w)
	}
}

func searchhandler(apps *Applications) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		constraints, err := parseSearchGetRequest(r)
		if err != nil {
			log.Println(err)
			//http.Error(w, "Error in search api 1", http.StatusInternalServerError)
			res, _ := json.Marshal("Error in search api 1")
			w.Write(res)
			return
		}

		matches := apps.GetMatches(constraints)
		res, err := json.Marshal(matches.Apps) //fixme maybe smaller response
		if err != nil {
			log.Println(err)
			//http.Error(w, "Error in search api 2", http.StatusInternalServerError)
			res, _ := json.Marshal("Error in search api 2")
			w.Write(res)
			return
		}
		_, err = w.Write(res)
		if err != nil {
			log.Println(err)
			//http.Error(w, "Error in search api 3", http.StatusInternalServerError)
			res, _ := json.Marshal("Error in search api 3")
			w.Write(res)
		}
	}
}
