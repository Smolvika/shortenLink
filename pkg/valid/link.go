package valid

import (
	"fmt"
	"net/url"
	"regexp"
	"time"
)

func ValidateOriginalURL(link string) error {
	if link == "" {
		return fmt.Errorf("empty query")
	}
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return fmt.Errorf("%v is a invalid base url", link)
	}

	return nil
}

func ValidateShortURL(link string) error {
	if link == "" {
		return fmt.Errorf("empty query")
	}
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
	layout := "2020-03-03"
	if _, err := time.Parse(layout, date); err != nil {
		return fmt.Errorf("%v is a invalid date", date)
	}
	return nil
}
