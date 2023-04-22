package postgres

import (
	"errors"
	"log"
	"shortenLink"
)

func (r Repository) CreateShortUrl(originalUrl, shortUrl, date string) (string, error) {
	log.Println("ps")
	_, err := r.db.Exec("INSERT INTO links (original_url, short_url, expiration_date) VALUES ($1,$2,$3)", originalUrl, shortUrl, date)
	if err != nil {
		if err == errors.New("duplicate key violates unique constraint") {
			return "", errors.New("there is already such a key")
		}
		return "", err
	}
	return shortUrl, nil
}

func (r Repository) GetShortUrl(url string) (string, error) {
	var link []shortenLink.OriginalLink
	err := r.db.Select(&link, "SELECT original_url FROM links  WHERE short_url=$1", url)
	if err != nil {
		return "", err
	}
	var origUrl string
	for _, v := range link {
		origUrl = v.Url
	}
	return origUrl, nil
}
func (r Repository) Delete(date string) error {
	_, err := r.db.Exec("DELETE FROM links WHERE expiration_date = $1", date)
	return err
}
