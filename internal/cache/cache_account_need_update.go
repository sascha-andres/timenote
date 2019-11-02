package cache

import (
	"go.etcd.io/bbolt"
)

const accountBucketName = "account"
const accountValueKeyName = "value"

// NeedUpdate returns true if cache needs a refresh
func (c *Cache) AccountNeedUpdate() (needUpdate bool) {
	_ = c.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(accountBucketName))
		if err != nil {
			return err
		}
		v := bucket.Get([]byte(metaKeyName))
		needUpdate = len(v) == 0 || checkIfUpdateRequired(v)
		return nil
	})
	return
}
