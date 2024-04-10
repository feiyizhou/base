package clients

import (
	"sync"
	"time"
)

// Memory mem cache struct
type Memory struct {
	sync.Mutex
	data map[string]*data
}

type data struct {
	Data    interface{}
	Expired *time.Time
}

// NewMemory create new mem cache
func NewMemory() *Memory {
	return &Memory{
		data: map[string]*data{},
	}
}

// Get return cached value
func (mem *Memory) Get(key string) interface{} {
	if ret, ok := mem.data[key]; ok {
		if ret.Expired != nil && ret.Expired.Before(time.Now()) {
			mem.deleteKey(key)
			return nil
		}
		return ret.Data
	}
	return nil
}

// IsExist check value exists in mem cache.
func (mem *Memory) IsExist(key string) bool {
	if ret, ok := mem.data[key]; ok {
		if ret.Expired != nil && ret.Expired.Before(time.Now()) {
			mem.deleteKey(key)
			return false
		}
		return true
	}
	return false
}

// Set cached value with key and expire time.
func (mem *Memory) Set(key string, val interface{}, timeout *time.Duration) {
	mem.Lock()
	defer mem.Unlock()
	value := &data{Data: val}
	if timeout != nil {
		expired := time.Now().Add(*timeout)
		value.Expired = &expired
	}
	mem.data[key] = value
	return
}

// Delete delete value in mem cache.
func (mem *Memory) Delete(key string) error {
	mem.deleteKey(key)
	return nil
}

// deleteKey
func (mem *Memory) deleteKey(key string) {
	mem.Lock()
	defer mem.Unlock()
	delete(mem.data, key)
}
