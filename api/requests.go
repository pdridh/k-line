package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

// Read a request body and parse it.
// The parsed json is loaded into v unless an error occurs
func ParseJSON(r *http.Request, v any) error {

	b, err := io.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}

// Given the query params in the url and a struct that has strongly typed fields,
// ParseQueryParams extracts the field from the query param if it exists and parses it into the type of the field it
// using the json tag in the struct. Params should be a pointer to a struct to get loaded into.
func ParseQueryParams(q url.Values, params any) {
	v := reflect.ValueOf(params).Elem()

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := v.Type().Field(i)
		paramValue := q.Get(fieldType.Tag.Get("json")) // Get from query

		if paramValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.Int64, reflect.Int32, reflect.Int:
			val, err := strconv.Atoi(paramValue)
			if err == nil {
				field.SetInt(int64(val))
			}
		case reflect.Float64:
			val, err := strconv.ParseFloat(paramValue, 64)
			if err == nil {
				field.SetFloat(val)
			}
		case reflect.Bool:
			val, err := strconv.ParseBool(paramValue)
			if err == nil {
				field.SetBool(val)
			}
		default:
			field.SetString(paramValue)
		}
	}
}
