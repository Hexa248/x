package controllers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func render(w http.ResponseWriter, page string, data any) {
	base := filepath.Join("views", "layout.html")
	view := filepath.Join("views", page)
	tmpl, err := template.ParseFiles(base, view)
	if err != nil {
		http.Error(w, "template parse error", http.StatusInternalServerError)
		log.Printf("template parse error (%s): %v", page, err)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		http.Error(w, "template render error", http.StatusInternalServerError)
		log.Printf("template render error (%s): %v", page, err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}
