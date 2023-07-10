package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Respond responds with converted data to the client.
func Respond(ctx context.Context, rw http.ResponseWriter, data any, code int) error {
	if code == http.StatusNoContent {
		rw.WriteHeader(code)
		return nil
	}

	var buf bytes.Buffer
	if err := EncodeBody(&buf, data); err != nil {
		return fmt.Errorf("encoding to buffer: %w", err)
	}

	rw.WriteHeader(code)

	if _, err := buf.WriteTo(rw); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}

// DecodeBody converts data from the client.
// If the value implements validation, it is executed.
func DecodeBody(body io.Reader, data any) error {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(data); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	if val, ok := data.(interface{ Validate() error }); ok {
		return fmt.Errorf("validation: %w", val.Validate())
	}

	return nil
}

// EncodeBody writes data to a writer after converting it to JSON.
func EncodeBody(rw io.Writer, data any) error {
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	return nil
}

// ProcessResponse processes the HTTP response by decoding its body into data.
// It returns an error if the status code indicates an error or if the body cannot be decoded.
func ProcessResponse(resp *http.Response, data any) error {
	if err := ErrFromStatusCode(resp.StatusCode); err != nil {
		return err
	}

	if err := DecodeBody(resp.Body, data); err != nil {
		return err
	}

	return nil
}
