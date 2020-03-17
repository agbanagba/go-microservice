package handlers

import (
	"net/http"
	"path/filepath"

	"github.com/agbanagba/go-microservice/images-api/files"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// Files handles the reading and writing to files
type Files struct {
	log     hclog.Logger
	storage files.Storage
}

// NewFiles creates a new files handler
func NewFiles(s files.Storage, log hclog.Logger) *Files {
	return &Files{log, s}
}

func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, rw, r)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid Path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: [id]/[filepath]", http.StatusBadRequest)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Save file for Product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.storage.Save(fp, r.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
