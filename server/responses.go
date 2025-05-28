package server

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// Write "v" (value that can be json marshaled) and send as response
// If the requestor accepts encoding then this encodes with gzip
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	var out io.Writer = w
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		out = gzw
		defer gzw.Close()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := out.Write(b); err != nil {
		return err
	}

	return nil
}

// Utility function that turns the given status, message and errors into an object ErrorResponse
// Which is used as a json response.
func WriteError(w http.ResponseWriter, r *http.Request, status int, message string, errors any) {
	e := ErrorResponse{
		Status:  status,
		Message: message,
		Errors:  errors,
	}

	if err := WriteJSON(w, r, status, e); err != nil {
		log.Println("failed to write to request")
	}
}

// Helper that calls WriteError() with args for an internal server error
func WriteInternalError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusInternalServerError, "Internal server error :(", nil)
}

// Helper that calls WriteError() with args for a bad request error
func WriteBadRequestError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusBadRequest, "Bad request", nil)
}
