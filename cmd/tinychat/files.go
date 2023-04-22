package main

import (
	"embed"
	"io"
	"net/http"
	"path"
	"strings"
)

//go:embed index.html spa.js style.css
var embeddedFiles embed.FS

func fileHandler(w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimPrefix(r.URL.Path, "/")
	if filePath == "" {
		filePath = "index.html" // Default to index.html if the path is empty
	}

	file, err := embeddedFiles.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	contentType := "text/plain"
	switch path.Ext(filePath) {
	case ".html":
		contentType = "text/html"
	case ".js":
		contentType = "text/javascript"
	case ".css":
		contentType = "text/css"
	}

	w.Header().Set("Content-Type", contentType)
	// Add cache-control headers for .html, .js, and .css files
	if path.Ext(filePath) == ".html" || path.Ext(filePath) == ".js" || path.Ext(filePath) == ".css" {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}
	io.Copy(w, file)
}
