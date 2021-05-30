package cache

import (
	"github.com/jason0x43/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
	"time"
)

// AccountSet can be used to update the value of the account in the cache
func (c *Cache) AccountSet(account *toggl.Account) error {
	return c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(accountBucketName))
		v, err := yaml.Marshal(account)
		if err != nil {
			return err
		}
		err = b.Put([]byte(accountValueKeyName), v)
		if err != nil {
			return err
		}
		m := MetaData{
			Updated:    time.Now(),
			NextUpdate: time.Now().Add(time.Duration(c.maxAge) * time.Minute),
		}
		meta, err := yaml.Marshal(m)
		if err != nil {
			return err
		}
		err = b.Put([]byte(metaKeyName), meta)
		return err
	})
}
