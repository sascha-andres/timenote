package cache

import (
	jsonv2 "encoding/json/v2"
	"os"
	"path/filepath"
)

// readJSON decodes the JSON file at path into v. A missing file is treated
// as "no data yet" and leaves v at its zero value.
func readJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer func() { _ = f.Close() }()
	return jsonv2.UnmarshalRead(f, v)
}

// writeJSON atomically writes v as JSON to path.
func writeJSON(path string, v any) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp.Name()) }()

	if err := jsonv2.MarshalWrite(tmp, v); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
