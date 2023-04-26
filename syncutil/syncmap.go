package syncmap

import "sync"

// Map a generic sync map.
type Map[K comparable, V any] struct {
	mx *sync.RWMutex
	m  map[K]V
}

// NewMap create a new threadsafe map.
func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		mx: &sync.RWMutex{},
		m:  make(map[K]V),
	}
}

// Get a value from this map.
func (sm *Map[K, V]) Get(k K) (V, bool) {
	sm.mx.RLock()
	defer sm.mx.RUnlock()
	v, ok := sm.m[k]
	return v, ok
}

// Store a value in this map.
func (sm *Map[K, V]) Store(k K, v V) {
	sm.mx.Lock()
	defer sm.mx.Unlock()
	sm.m[k] = v
}

// Range iterate the map.
func (sm *Map[K, V]) Range(fn func(K, V) bool) {
	cm := func() map[K]V {
		sm.mx.RLock()
		defer sm.mx.RUnlock()
		copyMap := make(map[K]V, len(sm.m))
		for k, v := range sm.m {
			copyMap[k] = v
		}

		return copyMap
	}()

	for k, v := range cm {
		if !fn(k, v) {
			break
		}
	}
}
