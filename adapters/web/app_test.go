package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JakubBizewski/jakubme_links/adapters/web"
	"github.com/JakubBizewski/jakubme_links/domain/model"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
	"github.com/JakubBizewski/jakubme_links/mocks"
)

func makeExpectedErrorBody(errorMessage string) string {
	json, _ := json.Marshal(map[string]string{
		"error": errorMessage,
	})

	return string(json)
}

func makeJsonRequestBody(targetUrl string) []byte {
	json, _ := json.Marshal(map[string]string{
		"targetUrl": targetUrl,
	})

	return json
}

func makeJsonResponseBody(shortCode string) string {
	json, _ := json.Marshal(map[string]string{
		"shortCode": shortCode,
	})

	return string(json)
}

func TestWebAppShort(t *testing.T) {
	shortLinkRepository := &mocks.MockShortLinkRepository{}
	shortLinkService := driver.CreateShortLinkService(shortLinkRepository)
	webApp := web.CreateWebApp(shortLinkService)

	t.Run("RedirectForExistigShortLink", func(t *testing.T) {
		url := "/testShortCode"
		expectedShortCode := "testShortCode"
		expectedTargetUrl := "https://example.com"

		shortLinkRepository.FindByShortCodeFunc = func(shortCode string) (model.ShortLink, error) {
			if shortCode != expectedShortCode {
				t.Errorf("Expected short code %s, but got %s", expectedShortCode, shortCode)
			}

			return model.ShortLink{
				ShortCode: shortCode,
				TargetUrl: expectedTargetUrl,
			}, nil
		}

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", url, nil)
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 302 {
			t.Errorf("Expected status code 302, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Header().Get("Location") != expectedTargetUrl {
			t.Errorf("Expected redirect to %s, but got %s", expectedTargetUrl, httpRecorder.Header().Get("Location"))
		}
	})

	t.Run("ShouldRedirectHomeOnShortLinkNotFound", func(t *testing.T) {
		shortLinkRepository.FindByShortCodeFunc = func(shortCode string) (model.ShortLink, error) {
			return model.ShortLink{}, nil
		}

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/notFoundShortCode", nil)
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 302 {
			t.Errorf("Expected status code 302, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Header().Get("Location") != "/" {
			t.Errorf("Expected redirect to %s, but got %s", "/", httpRecorder.Header().Get("Location"))
		}
	})

	t.Run("ShouldReturnAmbiguousTextOnError", func(t *testing.T) {
		expectedErrorMessage := "Something went wrong"

		shortLinkRepository.FindByShortCodeFunc = func(shortCode string) (model.ShortLink, error) {
			return model.ShortLink{}, errors.New("secret error")
		}

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/dummyErrorShortCode", nil)
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 500 {
			t.Errorf("Expected status code 500, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Body.String() != expectedErrorMessage {
			t.Errorf("Expected body %s, but got %s", expectedErrorMessage, httpRecorder.Body.String())
		}
	})

	t.Run("ShouldGenerateShortLink", func(t *testing.T) {
		expectedTargetUrl := "https://example.com"
		var generatedShortCode string

		shortLinkRepository.StoreFunc = func(shortLink model.ShortLink) error {
			if shortLink.TargetUrl != expectedTargetUrl {
				t.Errorf("Expected target url %s, but got %s", expectedTargetUrl, shortLink.TargetUrl)
			}

			generatedShortCode = shortLink.ShortCode

			return nil
		}

		body := map[string]string{
			"targetUrl": expectedTargetUrl,
		}

		jsonBody, _ := json.Marshal(body)

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/new", bytes.NewBuffer(jsonBody))
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 200 {
			t.Errorf("Expected status code 200, but got %d", httpRecorder.Code)
		}

		expectedResponse := makeJsonResponseBody(generatedShortCode)
		if httpRecorder.Body.String() != expectedResponse {
			t.Errorf("Expected body %s, but got %s", expectedResponse, httpRecorder.Body.String())
		}
	})

	t.Run("ShouldReturnErrorOnInvalidUrl", func(t *testing.T) {
		expectedErrorBody := makeExpectedErrorBody("Key: 'targetUrlPayload.TargetUrl' Error:Field validation for 'TargetUrl' failed on the 'url' tag")

		body := map[string]string{
			"targetUrl": "invalidUrl",
		}

		jsonBody, _ := json.Marshal(body)

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/new", bytes.NewBuffer(jsonBody))
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 400 {
			t.Errorf("Expected status code 400, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Body.String() != expectedErrorBody {
			t.Errorf("Expected body %s, but got %s", expectedErrorBody, httpRecorder.Body.String())
		}
	})

	t.Run("ShouldReturnErrorIfFailedToGenerateUniqueShortCode", func(t *testing.T) {
		expectedErrorBody := makeExpectedErrorBody("Failed to generate unique short code")

		shortLinkRepository.StoreFunc = func(shortLink model.ShortLink) error {
			return driven.ErrShortCodeAlreadyExists
		}

		body := map[string]string{
			"targetUrl": "https://example.com",
		}

		jsonBody, _ := json.Marshal(body)

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/new", bytes.NewBuffer(jsonBody))
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 500 {
			t.Errorf("Expected status code 500, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Body.String() != expectedErrorBody {
			t.Errorf("Expected body %s, but got %s", expectedErrorBody, httpRecorder.Body.String())
		}
	})
}
