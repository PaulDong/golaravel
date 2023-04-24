package golaravel

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

func (g *Golaravel) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single json value")
	}
	return nil
}

func (g *Golaravel) WriteJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (g *Golaravel) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (g *Golaravel) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Type", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

func (g *Golaravel) Error404(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusNotFound)
}

func (g *Golaravel) Error500(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusInternalServerError)
}

// ErrorUnauthorized sends an unauthorized status (client is not known)
func (g *Golaravel) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusUnauthorized)
}

// ErrorForbidden returns a forbidden status message (client is known)
func (g *Golaravel) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
	g.ErrorStatus(w, http.StatusForbidden)
}

func (g *Golaravel) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
