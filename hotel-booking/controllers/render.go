package controllers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func render(w http.ResponseWriter, page string, data any) {
	base := filepath.Join("views", "layout.html")
	if strings.HasPrefix(page, "admin/") {
		base = filepath.Join("views", "admin", "layout.html")
	}
	if strings.HasPrefix(page, "staff/") {
		base = filepath.Join("views", "staff", "layout.html")
	}

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
