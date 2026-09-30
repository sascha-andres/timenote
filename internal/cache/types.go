package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jason0x43/go-toggl"
)

type (
	Cache struct {
		maxAge int64
		dir    string
	}

	MetaData struct {
		Updated    time.Time
		NextUpdate time.Time
	}

	accountRecord struct {
		Meta    MetaData
		Account toggl.Account
	}

	projectsRecord struct {
		Meta     MetaData
		Projects []toggl.Project
	}

	clientsRecord struct {
		Meta    MetaData
		Clients []toggl.Client
	}
)

func accountPath(dir string) string {
	return filepath.Join(dir, "account.json")
}

func projectsPath(dir string, workspace int) string {
	return filepath.Join(dir, fmt.Sprintf("%d-projects.json", workspace))
}

func clientsPath(dir string, workspace int) string {
	return filepath.Join(dir, fmt.Sprintf("%d-clients.json", workspace))
}

func newMetaData(maxAge int64) MetaData {
	return MetaData{
		Updated:    time.Now(),
		NextUpdate: time.Now().Add(time.Duration(maxAge) * time.Minute),
	}
}

// NewCache creates a cache layer instance backed by JSON files in dir.
func NewCache(maxAge int, dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &Cache{
		maxAge: int64(maxAge),
		dir:    dir,
	}, nil
}
