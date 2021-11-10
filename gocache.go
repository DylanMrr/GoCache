package gocache

import (
	"sync"
	"time"
)

type StorageType int

const (
	Specific StorageType = iota
	Default
	Permanent
)

type item struct {
	data    interface{}
	expires int64
}

type Cache struct {
	items          sync.Map
	defaultExpires time.Duration
	close          chan struct{}
}

type Expires struct {
	ExpiresDuration time.Duration
}

func New(cleanTime time.Duration, defaultExpires time.Duration) *Cache {
	cache := &Cache{
		close:          make(chan struct{}),
		defaultExpires: defaultExpires,
	}
	go func() {
		ticker := time.NewTicker(cleanTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				now := time.Now().UnixNano()

				cache.items.Range(func(key, value interface{}) bool {
					item := value.(item)
					if item.expires > 0 && now > int64(item.expires) {
						cache.items.Delete(key)
					}
					return true
				})
			case <-cache.close:
				return
			}
		}
	}()

	return cache
}

func (this *Cache) Get(key interface{}) (interface{}, bool) {
	obj, exists := this.items.Load(key)

	if !exists {
		return nil, false
	}
	return obj.(item).data, true
}

func (this *Cache) Add(key interface{}, value interface{}, storageType StorageType, expires Expires) {
	var expiresDuration int64
	switch storageType {
	case Specific:
		expiresDuration = time.Now().Add(expires.ExpiresDuration).UnixNano()
	case Default:
		expiresDuration = time.Now().Add(this.defaultExpires).UnixNano()
	case Permanent:
		expiresDuration = 0
	}
	this.items.Store(key, item{
		data:    value,
		expires: expiresDuration,
	})
}

func (this *Cache) Range(f func(key, value interface{}) bool) {
	now := time.Now().UnixNano()

	fn := func(key, value interface{}) bool {
		item := value.(item)

		if item.expires > 0 && now > item.expires {
			return true
		}

		return f(key, item.data)
	}

	this.items.Range(fn)
}

func (this *Cache) Count() int {
	length := 0
	this.items.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

func (this *Cache) Delete(key interface{}) {
	this.items.Delete(key)
}

func (this *Cache) Close() {
	this.close <- struct{}{}
	this.items = sync.Map{}
}
