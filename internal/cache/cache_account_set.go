package cache

import (
	"github.com/sascha-andres/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

// AccountSet can be used to update the value of the account in the cache
func (c *Cache) AccountSet(account *toggl.Account) error {
	return c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(accountBucketName))
		v, err := yaml.Marshal(account)
		if err != nil {
			return err
		}
		return b.Put([]byte(accountValueKeyName), v)
	})
}
