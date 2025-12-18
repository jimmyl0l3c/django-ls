package safemap

import "sync"

type SafeMap[T any] struct {
	data map[string]T
	lock sync.RWMutex
}

func New[T any]() *SafeMap[T] {
	return &SafeMap[T]{data: make(map[string]T)}
}

func (sm *SafeMap[T]) Store(key string, item T) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	sm.data[key] = item
}

func (sm *SafeMap[T]) Load(key string) (T, bool) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	v, ok := sm.data[key]

	return v, ok
}

func (sm *SafeMap[T]) Delete(key string) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	delete(sm.data, key)
}
