package store

import (
	"hash/fnv"
	"sync"
)

type Shard struct {
	items map[string]interface{}
	mu   sync.RWMutex
}

func store() {

}

type Cache struct {
	shards []*Shard
}

func (c *Cache) getShardIndex(key string) int {
	size := len(c.shards)
	hashValue := fnv.New32a()
	hashValue.Write([]byte(key))
	keyValue := int(hashValue.Sum32())
	return keyValue % size
}

func NewCache(numberOfShards int) *Cache {
	cache := &Cache{
		shards: make([]*Shard, numberOfShards),
	}

	for i := range numberOfShards {
		cache.shards[i] = &Shard{
			items: make(map[string]interface{}),
		}
	}

	return cache
}

func (c *Cache) Set(key string, value interface{}) {
	shardIndex := c.getShardIndex(key)
	
	c.shards[shardIndex].mu.Lock()
	defer c.shards[shardIndex].mu.Unlock()

	c.shards[shardIndex].items[key] = value
}

func (c *Cache) Get(key string) (interface{}, bool) {
	shardIndex := c.getShardIndex(key)
	
	c.shards[shardIndex].mu.RLock()
	defer c.shards[shardIndex].mu.RUnlock()

	val, ok := c.shards[shardIndex].items[key];
	if ok {
		return  val, true
	}
	return nil, false
}