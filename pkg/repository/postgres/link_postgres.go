package postgres

import (
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) CreateShortUrl(originalUrl, shortUrl, date string) (string, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return "", err
	}
	var url string
	row, err := p.db.Query("SELECT original_url FROM links  WHERE short_url=$1", shortUrl)
	if err != nil {
		return "", err
	}
	if err := row.Scan(&url); err != nil {
		tx.Rollback()
		return "", err
	}
	if url != "" {
		tx.Rollback()
		return "", nil
	}
	_, err = p.db.Exec("INSERT INTO links (original_url, short_url, expiration_date) VALUES ($1,$2,$3)", originalUrl, shortUrl, date)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	return shortUrl, tx.Commit()
}

func (p *Postgres) GetShortUrl(url string) (string, error) {
	var origUrl string
	row, err := p.db.Query("SELECT original_url FROM links  WHERE short_url=$1", url)
	if err != nil {
		return "", err
	}
	if err := row.Scan(&origUrl); err != nil {
		return "", err
	}
	return origUrl, nil
}
func (p *Postgres) Delete(date string) error {
	_, err := p.db.Exec("DELETE FROM links WHERE expiration_date = $1", date)
	return err
}
