package inMemory

import (
	"errors"
)

func (r Repository) CreateShortUrl(originalUrl, shortUrl, date string) (string, error) {
	if _, ok := r.shortToOriginalAndDate[shortUrl]; ok {
		return "", errors.New("there is already such a key")
	}
	r.shortToOriginalAndDate[shortUrl] = UrlAndDate{
		url:  originalUrl,
		date: date,
	}
	return shortUrl, nil
}

func (r Repository) GetShortUrl(url string) (string, error) {
	if original, ok := r.shortToOriginalAndDate[url]; ok {
		return original.url, nil
	}
	return "", nil
}

func (r Repository) Delete(date string) error {
	for key, val := range r.shortToOriginalAndDate {
		if val.date == date {
			delete(r.shortToOriginalAndDate, key)
		}
	}
	return nil
}
