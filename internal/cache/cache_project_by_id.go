package cache

import (
	"fmt"
	"github.com/sascha-andres/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

func (c *Cache) ProjectByID(projectID, workspace int) (project *toggl.Project, err error) {
	err = c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf("%6d-projects", workspace)))
		v := b.Get([]byte(fmt.Sprintf("%10d", projectID)))
		var proj toggl.Project
		err = yaml.Unmarshal(v, &proj)
		project = &proj
		return nil
	})
	return
}
