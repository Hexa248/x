package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func render(w http.ResponseWriter, page string, data any) {
	base := filepath.Join("views", "layout.html")
	view := filepath.Join("views", page)
	tmpl := template.Must(template.ParseFiles(base, view))
	_ = tmpl.ExecuteTemplate(w, "layout", data)
}
