package inMemory

type UrlAndDate struct {
	url  string
	date string
}

type Repository struct {
	shortToOriginalAndDate map[string]UrlAndDate
}

func New(n int) *Repository {
	shortToOriginalAndDate := make(map[string]UrlAndDate, n)
	return &Repository{
		shortToOriginalAndDate: shortToOriginalAndDate,
	}
}
