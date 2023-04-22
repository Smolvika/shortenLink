package shortenLink

type OriginalLink struct {
	Url string `db:"original_url"`
}

type Link struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ShortUrl    string
	Date        string `json:"date" binding:"required"`
}
