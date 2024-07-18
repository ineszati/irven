package main

import (
	"fmt"
	"sync"
)

type MyStruct struct {
	// Define the fields of MyStruct here
}

type PoolManager struct {
	pool     *sync.Pool
	maxItems int
	current  int
	mu       sync.Mutex
}

func NewPoolManager(maxItems int) *PoolManager {
	manager := &PoolManager{
		maxItems: maxItems,
	}

	manager.pool = &sync.Pool{
		New: func() interface{} {
			manager.mu.Lock()
			defer manager.mu.Unlock()
			if manager.current < manager.maxItems {
				manager.current++
				return &MyStruct{}
			}
			return nil // Return nil if maxItems limit is reached
		},
	}

	return manager
}

func (pm *PoolManager) Get() *MyStruct {
	instance := pm.pool.Get()
	if instance == nil {
		return nil
	}
	return instance.(*MyStruct)
}

func (pm *PoolManager) Put(item *MyStruct) {
	pm.pool.Put(item)
}

func (pm *PoolManager) Resize(newMax int) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.maxItems = newMax

	// You may also want to adjust current count or manage existing items
}

func main() {
	manager := NewPoolManager(10)

	instance := manager.Get()
	if instance != nil {
		fmt.Println("Got an instance:", instance)
	} else {
		fmt.Println("No instance available")
	}

	manager.Put(instance)

	// Resize the pool
	manager.Resize(20)

	reusedInstance := manager.Get()
	if reusedInstance != nil {
		fmt.Println("Got a reused instance:", reusedInstance)
	} else {
		fmt.Println("No reused instance available")
	}

	fmt.Println(manager.maxItems) // Output: 20
}
