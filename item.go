package store

type Item struct{ Key, Data []byte }

func NewItem(key, data []byte) Item                 { return Item{Key: key, Data: data} }
func NewStringKeyItem(key string, data []byte) Item { return NewItem([]byte(key), data) }
