package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

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

// UploadREST post a file to a file storage location
func (f *Files) UploadREST(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	filename := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", filename)

	f.saveFile(id, filename, rw, r.Body)
}

// UploadMultiPart ...
func (f *Files) UploadMultiPart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 * 1024)
	if err != nil {
		f.log.Error("No multipart form detected", "error", err)
		http.Error(rw, "No multipart form data provided", http.StatusBadRequest)
		return
	}

	id, idErr := strconv.Atoi(r.FormValue("id"))
	f.log.Info("Form ", "id", id)
	if idErr != nil {
		f.log.Error("id in mulipart not an parsable as an integer", "error", err)
		http.Error(rw, "Expected integer as id", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		f.log.Error("No file with multipart request", "error", err)
		http.Error(rw, "Expected file with the multipart request", http.StatusBadRequest)
		return
	}

	f.saveFile(string(id), header.Filename, rw, file)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid Path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: [id]/[filepath]", http.StatusBadRequest)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for Product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.storage.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
