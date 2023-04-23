package inMemory

type UrlAndDate struct {
	url  string
	date string
}

type inMemory struct {
	shortToOriginalAndDate map[string]UrlAndDate
}

func New() *inMemory {
	shortToOriginalAndDate := make(map[string]UrlAndDate)
	return &inMemory{
		shortToOriginalAndDate: shortToOriginalAndDate,
	}
}

func (m *inMemory) CreateShortUrl(originalUrl, shortUrl, date string) (string, error) {
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
	if original, ok := m.shortToOriginalAndDate[url]; ok {
		return original.url, nil
	}
	return "", nil
}

func (m *inMemory) Delete(date string) error {
	for key, val := range m.shortToOriginalAndDate {
		if val.date == date {
			delete(m.shortToOriginalAndDate, key)
		}
	}
	return nil
}
