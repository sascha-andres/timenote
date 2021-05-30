package cache

import (
	"github.com/jason0x43/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

// AccountSet can be used to update the value of the account in the cache
func (c *Cache) GetAccount() (acc toggl.Account, err error) {
	_ = c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(accountBucketName))
		v := b.Get([]byte(accountValueKeyName))
		err = yaml.Unmarshal(v, &acc)
		return nil
	})
	return
}
