package factory

import "sync"

type Supplier func(string) any

type Singleton struct {
	cache    map[string]any
	locker   *sync.Mutex
	supplier Supplier
}

func NewSingleton(supplier Supplier) *Singleton {
	return &Singleton{
		cache:    make(map[string]any),
		locker:   &sync.Mutex{},
		supplier: supplier,
	}
}

func (s *Singleton) Get(key string) any {
	if value, hit := s.cache[key]; hit {
		return value
	}
	s.locker.Lock()
	defer s.locker.Unlock()
	if value, hit := s.cache[key]; hit {
		return value
	}
	s.cache[key] = s.supplier(key)
	return s.cache[key]
}
