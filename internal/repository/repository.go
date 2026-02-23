package repository

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("short code not found")

type URL struct {
	Code        string
	OriginalURL string
	Clicks      int64
}

type URLRepository interface {
	Save(url URL) error
	FindByCode(code string) (URL, error)
	IncrementClicks(code string) error
	ListAll() []URL
}

type InMemoryRepository struct {
	mu   sync.RWMutex
	urls map[string]URL
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		urls: make(map[string]URL),
	}
}

func (r *InMemoryRepository) Save(url URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.urls[url.Code] = url
	return nil
}

func (r *InMemoryRepository) FindByCode(code string) (URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.urls[code]
	if !ok {
		return URL{}, ErrNotFound
	}
	return url, nil
}

func (r *InMemoryRepository) IncrementClicks(code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	url, ok := r.urls[code]
	if !ok {
		return ErrNotFound
	}
	url.Clicks++
	r.urls[code] = url
	return nil
}

func (r *InMemoryRepository) ListAll() []URL {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]URL, 0, len(r.urls))
	for _, url := range r.urls {
		list = append(list, url)
	}
	return list
}
