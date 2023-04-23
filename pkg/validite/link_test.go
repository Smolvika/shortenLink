package validite

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidateOriginalURL(t *testing.T) {
	testTable := []struct {
		name     string
		arr      string
		expected error
	}{
		{
			name:     "OK",
			arr:      "https://habr.com/ru/articles/555920/",
			expected: nil,
		},
		{
			name:     "Invalid original url",
			arr:      "empty/query",
			expected: errors.New("empty/query is a invalid original URL"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := OriginalURL(testCase.arr)
			assert.Equal(t, testCase.expected, err)

		})
	}
}

func Test_ValidateShortURL(t *testing.T) {
	testTable := []struct {
		name     string
		arr      string
		expected error
	}{
		{
			name:     "OK",
			arr:      "NdJgND_234",
			expected: nil,
		},
		{
			name:     "Invalid values in url",
			arr:      "gh0c+lfv,1",
			expected: errors.New("gh0c+lfv,1 is a invalid URL"),
		},
		{
			name:     "Invalid length url",
			arr:      "gh0c",
			expected: errors.New("gh0c is a invalid URL"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := ShortURL(testCase.arr)
			assert.Equal(t, testCase.expected, err)

		})
	}
}
func Test_ValidateDate(t *testing.T) {
	testTable := []struct {
		name     string
		arr      string
		expected error
	}{
		{
			name:     "OK",
			arr:      "2020-03-23",
			expected: nil,
		},
		{
			name:     "Invalid values in date",
			arr:      "2023-23-23",
			expected: errors.New("2023-23-23 is a invalid date"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := Date(testCase.arr)
			assert.Equal(t, testCase.expected, err)

		})
	}
}
