package cache

import (
	"fmt"
	"github.com/sascha-andres/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

func (c *Cache) Clients(workspace int) (clients []toggl.Client, err error) {
	err = c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf(clientBucketNameTemplate, workspace)))
		v := b.Get([]byte("_all"))
		err = yaml.Unmarshal(v, &clients)
		return nil
	})
	return
}
