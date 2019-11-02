package cache

import (
	"github.com/sascha-andres/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

// AccountSet can be used to update the value of the account in the cache
func (c *Cache) AccountGet(account *toggl.Account) (acc *toggl.Account, err error) {
	_ = c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(accountBucketName))
		v := b.Get([]byte(accountValueKeyName))
		var a toggl.Account
		err = yaml.Unmarshal(v, &a)
		acc = &a
		return nil
	})
	return
}
