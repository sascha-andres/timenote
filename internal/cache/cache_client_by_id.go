package cache

import (
	"fmt"
	"github.com/jason0x43/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

func (c *Cache) ClientByID(clientID, workspace int) (client *toggl.Client, err error) {
	err = c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf(clientBucketNameTemplate, workspace)))
		v := b.Get([]byte(fmt.Sprintf("%10d", clientID)))
		var cl toggl.Client
		err = yaml.Unmarshal(v, &cl)
		client = &cl
		return nil
	})
	return
}
