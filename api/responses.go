package api

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
func WriteError(w http.ResponseWriter, r *http.Request, status int, apiErr APIError, errors any) {
	e := NewErrorResponse(status, apiErr.Code, apiErr.Message, errors)

	if err := WriteJSON(w, r, status, e); err != nil {
		log.Println("failed to write to request")
	}
}

// TODO make these generalized? tala ko

// Helper that calls WriteError() with args for an internal server error
func WriteInternalError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusInternalServerError, ErrHTTPInternal, nil)
}

// Helper that calls WriteError() with args for a bad request error
func WriteBadRequestError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusBadRequest, ErrHTTPBadRequest, nil)
}

// Helper that calls WriteError() with args for a resource not found error
func WriteNotFoundError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusNotFound, ErrHTTPNotFound, nil)
}

func WriteInvalidJWTError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusUnauthorized, ErrJWTInvalid, nil)
}

func WriteForbiddenError(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusForbidden, ErrHTTPForbidden, nil)
}

func WriteValidationError(w http.ResponseWriter, r *http.Request, err error) {
	v := FormatValidationErrors(err)
	WriteError(w, r, http.StatusBadRequest, ErrJSONValidation, v)
}
