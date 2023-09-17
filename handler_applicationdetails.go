package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"

	"biehdc.webapp.applister/cache"
)

func applicationparser(apps *Applications, tpl string) http.HandlerFunc {
	apptemplate := template.Must(template.ParseFiles(tpl))

	templatecache := cache.NewTemplateCache()

	return func(w http.ResponseWriter, r *http.Request) {
		cliurl, err := url.Parse(r.URL.String())
		if err != nil {
			log.Println(err)
			http.Error(w, "Error in Applicationdetails 1", http.StatusInternalServerError)
			return
		}
		name := path.Base(cliurl.Path)

		cache, exists := templatecache.GetEntry(name)
		if exists {
			// We already got the thing, return it
			//println(name, " from cache")
			io.Copy(w, bytes.NewBuffer(cache))
			return
		}

		theapp := apps.FindEntry(name)
		if theapp == nil {
			// If we didnt find the app, dispatch a search
			http.Redirect(w, r, "/search?appname="+name, http.StatusSeeOther)
			return
		}

		buf := &bytes.Buffer{}
		err = apptemplate.Execute(buf, *theapp)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error in Applicationdetails 2", http.StatusInternalServerError)
			return
		}
		bufasbytes := buf.Bytes()
		templatecache.SetEntry(name, bufasbytes)
		io.Copy(w, bytes.NewBuffer(bufasbytes))
		//println(name, " from cold")
		return
	}
}
