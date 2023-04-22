package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"shortenLink"
	"shortenLink/pkg/service"
	mockservice "shortenLink/pkg/service/mocks"
	"testing"
)

func TestHandler_CreateShortUrl(t *testing.T) {

	type mockBehavior func(r *mockservice.MockRepository, input shortenLink.Link)

	testTable := []struct {
		name                 string
		inputBody            string
		input                shortenLink.Link
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"original_url": "https://habr.com/ru/articles/555920/","date": "2020-03-03"}`,
			input: shortenLink.Link{
				OriginalUrl: "https://habr.com/ru/articles/555920/",
				Date:        "2020-03-03",
			},
			mockBehavior: func(r *mockservice.MockRepository, input shortenLink.Link) {
				r.EXPECT().CreateShortUrl(input.OriginalUrl, gomock.Any(), input.Date).Return("9_gfyrnTY5", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"shortURL":"9_gfyrnTY5"}`,
		},
		{
			name:      "Some data is missing",
			inputBody: `{"date": "2020-03-03"}`,
			input: shortenLink.Link{
				Date: "2020-03-03",
			},
			mockBehavior:         func(r *mockservice.MockRepository, input shortenLink.Link) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"Key: 'LinkInput.OriginalUrl' Error:Field validation for 'OriginalUrl' failed on the 'required' tag"}`,
		},
		{
			name:      "Invalid Original URL",
			inputBody: `{"original_url": "555920", "date": "2020-03-03"}`,
			input: shortenLink.Link{
				OriginalUrl: "555920",
				Date:        "2020-03-03",
			},
			mockBehavior:         func(r *mockservice.MockRepository, input shortenLink.Link) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"555920 is a invalid original url"}`,
		},
		{
			name:      "Invalid Date",
			inputBody: `{"original_url":"https://habr.com/ru/articles/555920/", "date": "2020-33-33"}`,
			input: shortenLink.Link{
				OriginalUrl: "https://habr.com/ru/articles/555920/",
				Date:        "2020-33-33",
			},
			mockBehavior:         func(r *mockservice.MockRepository, input shortenLink.Link) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"2020-33-33 is a invalid date"}`,
		},
		{
			name:      "Problems on the service",
			inputBody: `{"original_url":"https://habr.com/ru/articles/555920/", "date": "2020-04-03"}`,
			input: shortenLink.Link{
				OriginalUrl: "https://habr.com/ru/articles/555920/",
				Date:        "2020-04-03",
			},
			mockBehavior: func(r *mockservice.MockRepository, input shortenLink.Link) {
				r.EXPECT().CreateShortUrl(input.OriginalUrl, gomock.Any(), input.Date).Return("", errors.New("problems on the service"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"problems on the service"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockservice.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.input)

			services := service.New(repo)
			handler := New(services)

			r := gin.New()
			r.POST("/tokens/short", handler.CreateShortUrl)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/tokens/short", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
func TestHandler_GetOriginalUrl(t *testing.T) {
	type mockBehavior func(r *mockservice.MockRepository, url string)

	testTable := []struct {
		name                 string
		inputBody            string
		url                  string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			url:  "9_gfyrnTY5",
			mockBehavior: func(r *mockservice.MockRepository, url string) {
				r.EXPECT().GetShortUrl(url).Return("https://habr.com/ru/articles/555920/", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"originalURL":"https://habr.com/ru/articles/555920/"}`,
		},
		{
			name:                 "Invalid URL",
			url:                  "9JDNVJ",
			mockBehavior:         func(r *mockservice.MockRepository, url string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"9JDNVJ is a invalid URL"}`,
		},
		{
			name: "Problems on the service",
			url:  "9_gfyrnTY5",
			mockBehavior: func(r *mockservice.MockRepository, url string) {
				r.EXPECT().GetShortUrl(url).Return("", errors.New("problems on the service"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"problems on the service"}`,
		},
		{
			name: "Invalid URL",
			url:  "9_gfyrnTY5",
			mockBehavior: func(r *mockservice.MockRepository, url string) {
				r.EXPECT().GetShortUrl(url).Return("", errors.New("not such ShortURL"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"not such ShortURL"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockservice.NewMockRepository(c)
			testCase.mockBehavior(repo, testCase.url)

			services := service.New(repo)
			handler := New(services)

			r := gin.New()
			r.GET("/tokens/:short", handler.GetOriginalUrl)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/tokens/"+testCase.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
