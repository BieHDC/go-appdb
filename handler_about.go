package main

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

func aboutpage(apps *Applications, tpl string) http.HandlerFunc {
	abouttemplate := template.Must(template.ParseFiles(tpl))

	buf := &bytes.Buffer{}
	err := abouttemplate.Execute(buf, apps)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, bytes.NewBuffer(buf.Bytes()))
	}
}
