package main

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
)

//go:embed doc
var docFS embed.FS

func swaggerHandler(file string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := fs.ReadFile(docFS, path.Join("doc", file))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		switch path.Ext(file) {
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		}
		w.Write(data)
	}
}
