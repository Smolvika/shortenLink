package inMemory

import "sync"

type UrlAndDate struct {
	url  string
	date string
}

type inMemory struct {
	mu                     *sync.RWMutex
	shortToOriginalAndDate map[string]UrlAndDate
}

func New() *inMemory {
	shortToOriginalAndDate := make(map[string]UrlAndDate)
	return &inMemory{
		shortToOriginalAndDate: shortToOriginalAndDate,
		mu:                     &sync.RWMutex{},
	}
}

func (m *inMemory) CreateShortUrl(originalUrl, shortUrl, date string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.shortToOriginalAndDate[shortUrl]; ok {
		return "", nil
	}
	m.shortToOriginalAndDate[shortUrl] = UrlAndDate{
		url:  originalUrl,
		date: date,
	}
	return shortUrl, nil
}

func (m *inMemory) GetShortUrl(url string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if original, ok := m.shortToOriginalAndDate[url]; ok {
		return original.url, nil
	}

	return "", nil
}

func (m *inMemory) Delete(date string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for key, val := range m.shortToOriginalAndDate {
		if val.date <= date {
			delete(m.shortToOriginalAndDate, key)
		}
	}
	return nil
}
