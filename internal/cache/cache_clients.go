package cache

import (
	"fmt"
	"github.com/jason0x43/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

func (c *Cache) Clients(workspace int) (clients []toggl.Client, err error) {
	err = c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf(clientBucketNameTemplate, workspace)))
		v := b.Get([]byte(allKeyName))
		err = yaml.Unmarshal(v, &clients)
		return nil
	})
	return
}
