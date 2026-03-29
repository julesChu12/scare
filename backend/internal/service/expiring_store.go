package service

import (
	"strconv"
	"sync"
	"time"
)

type expiringValue struct {
	value     string
	expiresAt time.Time
}

type expiringStore struct {
	mu    sync.Mutex
	items map[string]expiringValue
}

func newExpiringStore() *expiringStore {
	return &expiringStore{
		items: make(map[string]expiringValue),
	}
}

func (s *expiringStore) get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if !ok {
		return "", false
	}
	if expired(item.expiresAt) {
		delete(s.items, key)
		return "", false
	}
	return item.value, true
}

func (s *expiringStore) set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = expiringValue{
		value:     value,
		expiresAt: expiresAt(ttl),
	}
}

func (s *expiringStore) del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

func (s *expiringStore) exists(key string) bool {
	_, ok := s.get(key)
	return ok
}

func (s *expiringStore) incr(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, ok := s.items[key]
	if ok && expired(item.expiresAt) {
		delete(s.items, key)
		ok = false
	}

	var current int64
	if ok {
		parsed, err := strconv.ParseInt(item.value, 10, 64)
		if err != nil {
			return 0, err
		}
		current = parsed
	}

	current++
	s.items[key] = expiringValue{
		value:     strconv.FormatInt(current, 10),
		expiresAt: item.expiresAt,
	}
	return current, nil
}

func expired(t time.Time) bool {
	return !t.IsZero() && time.Now().After(t)
}

func expiresAt(ttl time.Duration) time.Time {
	if ttl <= 0 {
		return time.Time{}
	}
	return time.Now().Add(ttl)
}
