package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

func indexpage(apps *Applications, tpl string, entriesperpage uint) http.HandlerFunc {
	indextemplate := template.Must(template.ParseFiles(tpl))

	allentries := &bytes.Buffer{}
	err := indextemplate.Execute(allentries, apps)
	if err != nil {
		panic(err)
	}

	allpages := generateAllPages(apps, entriesperpage, indextemplate)

	return func(w http.ResponseWriter, r *http.Request) {
		cliurl, err := url.Parse(r.URL.String())
		if err != nil {
			log.Println(err)
			http.Error(w, "Error in indexhandler 1", http.StatusInternalServerError)
			return
		}
		params := cliurl.Query()
		if cliurl.Path != "/" {
			//if the path is not empty, we start a search query
			http.Redirect(w, r, "/search?appname="+cliurl.Path[1:], http.StatusSeeOther)
			return
		}

		page := params.Get("page")
		if page == "" {
			page = "0"
		}
		if page == "all" {
			io.Copy(w, bytes.NewBuffer(allentries.Bytes()))
			return
		}

		pagei, err := strconv.Atoi(page)
		if err != nil || pagei < 0 {
			//if for whatever reason
			io.Copy(w, bytes.NewBuffer(allentries.Bytes()))
			return
		}

		cache, exists := allpages[uint(pagei)]
		if exists {
			io.Copy(w, bytes.NewBuffer(cache.Bytes()))
			return
		} else {
			//invalid page number -> all entries
			io.Copy(w, bytes.NewBuffer(allentries.Bytes()))
			return
		}
	}
}

func generateAllPages(apps *Applications, entriesperpage uint, indextemplate *template.Template) map[uint]*bytes.Buffer {
	allpages := make(map[uint]*bytes.Buffer)

	//static
	paged := NewApplications(apps.Title)
	paged.PG.EntriesPerPage = entriesperpage
	paged.PG.TotalPages = uint(math.Ceil(float64(len(apps.Apps)) / float64(paged.PG.EntriesPerPage)))

	//dynamic
	var nextPage uint
	var err error
	for i := uint(0); i < paged.PG.TotalPages; i++ {
		nextPage = i + 1
		paged.PG.NextPage = nextPage
		paged.GetSlice(apps)

		buf := &bytes.Buffer{}
		err = indextemplate.Execute(buf, paged)
		if err != nil {
			panic(err)
		}
		allpages[i] = buf
	}

	return allpages
}
