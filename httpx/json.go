package httpx

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

// JSONContentType represents MIME type for JSON content.
const JSONContentType = "application/json"

// JSON encodes json content to the ResponseWriter.
func JSON(w http.ResponseWriter, code int, v interface{}) error {
	raw, err := jsoniter.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(code)

	if _, err = w.Write(raw); err != nil {
		return err
	}

	return nil
}
