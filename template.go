package main

import "html/template"

type BasePageData struct {
	StoreName string
}

func parsePageTemplate(filename string) *template.Template {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		filename,
	))
	return tmpl
}
