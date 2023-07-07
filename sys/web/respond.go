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
		return err
	}

	rw.WriteHeader(code)

	if _, err := buf.WriteTo(rw); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}

	return nil
}

// Decode converts data from the client.
// If the value implements validation, it is executed.
func Decode(req *http.Request, data any) error {
	if err := DecodeBody(req.Body, data); err != nil {
		return err
	}

	if val, ok := data.(interface{ OK() error }); !ok {
		return fmt.Errorf("validation: %w", val.OK())
	}

	return nil
}

func DecodeBody(body io.ReadCloser, data any) error {
	defer body.Close()

	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("decoding body: %w", err)
	}

	return nil
}

func EncodeBody(rw io.Writer, data any) error {
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		return fmt.Errorf("encoding body: %w", err)
	}

	return nil
}

func ProcessResponse[T any](resp *http.Response) (*T, error) {
	if err := ErrFromStatusCode(resp.StatusCode); err != nil {
		return nil, err
	}

	var data T
	if err := DecodeBody(resp.Body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
