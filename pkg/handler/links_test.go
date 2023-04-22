package handler

import (
	"bytes"
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
	// Init Test Table
	type mockBehavior func(r *mockservice.MockRepository, input shortenLink.Link)

	tests := []struct {
		name                 string
		inputBody            string
		input                shortenLink.Link
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"original_url": "https://habr.com/ru/articles/555920/", "date": "2020-03-03"}`,
			input: shortenLink.Link{
				OriginalUrl: "https://habr.com/ru/articles/555920/",
				ShortUrl:    "",
				Date:        "2020-03-03",
			},
			mockBehavior: func(r *mockservice.MockRepository, input shortenLink.Link) {
				r.EXPECT().CreateShortUrl(input.OriginalUrl, input.ShortUrl, input.Date).Return("9_gfyrnTY5", nil).Times(1)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"ShortUrl":9_gfyrnTY5}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockservice.NewMockRepository(c)
			test.mockBehavior(repo, test.input)

			//services := &service.Service{Repos: repo}
			services := service.New(repo)
			//handler := Handler{services}
			handler := New(services)

			// Init Endpoint
			//gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/tokens/short", handler.CreateShortUrl)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/tokens/short", bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
