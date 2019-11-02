package cache

import (
	"fmt"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
)

// ProjectMetaData returns meta data about the projects cache
func (c *Cache) ProjectMetaData(workspace int) (m *MetaData, err error) {
	_ = c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(fmt.Sprintf(projectBucketNameTemplate, workspace)))
		v := b.Get([]byte("_meta"))
		var md MetaData
		err = yaml.Unmarshal(v, &md)
		m = &md
		return nil
	})
	return
}
