package cache

import (
	"fmt"
	"github.com/sascha-andres/go-toggl"
	"go.etcd.io/bbolt"
	"gopkg.in/yaml.v2"
	"time"
)

func (c *Cache) SetClients(workspace int, clients []toggl.Client) error {
	return c.db.Update(func(tx *bbolt.Tx) error {
		allData, err := yaml.Marshal(clients)
		if err != nil {
			return err
		}
		b := tx.Bucket([]byte(fmt.Sprintf("%6d-clients", workspace)))
		err = b.Put([]byte("_all"), allData)
		if err != nil {
			return err
		}
		for _, client := range clients {
			singleData, err := yaml.Marshal(client)
			if err != nil {
				return err
			}
			err = b.Put([]byte(fmt.Sprintf("%10d", client.ID)), singleData)
			if err != nil {
				return err
			}
		}

		m := MetaData{
			Updated:    time.Now(),
			NextUpdate: time.Now().Add(time.Duration(c.maxAge) * time.Minute),
		}
		meta, err := yaml.Marshal(m)
		if err != nil {
			return err
		}
		err = b.Put([]byte("_meta"), meta)
		if err != nil {
			return err
		}
		return nil
	})
}
