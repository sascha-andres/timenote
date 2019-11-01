package cache

type Cache struct{}

func NewCache() (*Cache, error) {
	return &Cache{}, nil
}
