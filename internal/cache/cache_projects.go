package cache

import (
	"fmt"
	"github.com/sascha-andres/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

func (c *Cache) Projects(workspace int) (projects []toggl.Project, err error) {
	err = c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf("%6d-projects", workspace)))
		v := b.Get([]byte("_all"))
		err = yaml.Unmarshal(v, &projects)
		return nil
	})
	return
}
