package valid

import (
	"fmt"
	"net/url"
	"regexp"
)

func ValidateOriginalURL(link string) error {
	if link == "" {
		return fmt.Errorf("empty query")
	}
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("%v is a invalid original url", link)
	}

	return nil
}

func ValidateShortURL(link string) error {
	pattern := `^[a-zA-Z0-9_]{10}$`
	if valid, _ := regexp.Match(pattern, []byte(link)); !valid {
		return fmt.Errorf("%v is a invalid URL", link)
	}

	return nil
}

func ValidateDate(date string) error {
	if date == "" {
		return fmt.Errorf("empty query")
	}

	pattern := `^[0-9]{4}\-(0[1-9]|1[012])\-(0[1-9]|[12][0-9]|3[01])$`
	if valid, _ := regexp.Match(pattern, []byte(date)); !valid {
		return fmt.Errorf("%v is a invalid date", date)
	}
	return nil
}
