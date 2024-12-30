package utils

import (
	"encoding/json"
	"io"
	"log/slog"
)

// decodes reader with JSON-formatted content r into the STRUCT pointed at by st
func DecodeJson[T any](r io.Reader, st *T) error {
	decoder := json.NewDecoder(r)

	err := decoder.Decode(st)
	if err != nil {
		slog.Error("error decoding json", "error", err)
		return err
	}
	return nil
}
