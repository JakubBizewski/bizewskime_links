package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func makeJSONRequestBody(targetURL string) []byte {
	json, _ := json.Marshal(map[string]string{
		"targetURL": targetURL,
	})

	return json
}

func makeJSONResponseBody(shortCode string) string {
	json, _ := json.Marshal(map[string]string{
		"shortCode": shortCode,
	})

	return string(json)
}

func hasUserIDCookie(recorder *httptest.ResponseRecorder) bool {
	setCookieValue := recorder.Header().Get("Set-Cookie")

	return strings.Contains(setCookieValue, "user-id=")
}

//nolint:gocognit // This is a test suite, so it's ok to have a lot of code here
func TestWebAppShort(t *testing.T) {
	encryptionService := &mocks.MockEncryptionService{}
	shortLinkRepository := &mocks.MockShortLinkRepository{}
	shortLinkService := driver.CreateShortLinkService(shortLinkRepository)

	webApp := web.CreateWebApp(shortLinkService, encryptionService)

	t.Run("RedirectForExistigShortLink", func(t *testing.T) {
		url := "/testShortCode"
		expectedShortCode := "testShortCode"
		expectedTargetURL := "https://example.com"

		shortLinkRepository.FindByShortCodeFunc = func(shortCode string) (model.ShortLink, error) {
			if shortCode != expectedShortCode {
				t.Errorf("Expected short code %s, but got %s", expectedShortCode, shortCode)
			}

			return model.ShortLink{
				ShortCode: shortCode,
				TargetURL: expectedTargetURL,
			}, nil
		}

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 302 {
			t.Errorf("Expected status code 302, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Header().Get("Location") != expectedTargetURL {
			t.Errorf("Expected redirect to %s, but got %s", expectedTargetURL, httpRecorder.Header().Get("Location"))
		}

		if httpRecorder.Header().Get("Set-Cookie") != "" {
			t.Errorf("Expected no cookie to be set, but got %s", httpRecorder.Header().Get("Set-Cookie"))
		}
	})

	t.Run("ShouldRedirectHomeOnShortLinkNotFound", func(t *testing.T) {
		shortLinkRepository.FindByShortCodeFunc = func(shortCode string) (model.ShortLink, error) {
			return model.ShortLink{}, nil
		}

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/notFoundShortCode", nil)
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 302 {
			t.Errorf("Expected status code 302, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Header().Get("Location") != "/" {
			t.Errorf("Expected redirect to %s, but got %s", "/", httpRecorder.Header().Get("Location"))
		}

		if httpRecorder.Header().Get("Set-Cookie") != "" {
			t.Errorf("Expected no cookie to be set, but got %s", httpRecorder.Header().Get("Set-Cookie"))
		}
	})

	t.Run("ShouldReturnAmbiguousTextOnError", func(t *testing.T) {
		expectedErrorMessage := "Something went wrong"

		shortLinkRepository.FindByShortCodeFunc = func(shortCode string) (model.ShortLink, error) {
			return model.ShortLink{}, errors.New("secret error")
		}

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/dummyErrorShortCode", nil)
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 500 {
			t.Errorf("Expected status code 500, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Body.String() != expectedErrorMessage {
			t.Errorf("Expected body %s, but got %s", expectedErrorMessage, httpRecorder.Body.String())
		}

		if httpRecorder.Header().Get("Set-Cookie") != "" {
			t.Errorf("Expected no cookie to be set, but got %s", httpRecorder.Header().Get("Set-Cookie"))
		}
	})

	t.Run("ShouldGenerateShortLink", func(t *testing.T) {
		expectedTargetURL := "https://example.com"
		var generatedShortCode string

		shortLinkRepository.StoreFunc = func(shortLink model.ShortLink) error {
			if shortLink.TargetURL != expectedTargetURL {
				t.Errorf("Expected target url %s, but got %s", expectedTargetURL, shortLink.TargetURL)
			}

			generatedShortCode = shortLink.ShortCode

			return nil
		}

		jsonBody := makeJSONRequestBody(expectedTargetURL)

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/new", bytes.NewBuffer(jsonBody))
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 200 {
			t.Errorf("Expected status code 200, but got %d", httpRecorder.Code)
		}

		expectedResponse := makeJSONResponseBody(generatedShortCode)
		if httpRecorder.Body.String() != expectedResponse {
			t.Errorf("Expected body %s, but got %s", expectedResponse, httpRecorder.Body.String())
		}

		if !hasUserIDCookie(httpRecorder) {
			t.Errorf("Expected to set cookie, but got %s", httpRecorder.Header().Get("Set-Cookie"))
		}
	})

	t.Run("ShouldReturnErrorOnInvalidURL", func(t *testing.T) {
		expectedErrorBody := makeExpectedErrorBody(
			"Key: 'targetURLPayload.TargetURL' Error:Field validation for 'TargetURL' failed on the 'url' tag",
		)

		jsonBody := makeJSONRequestBody("invalidURL")

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/new", bytes.NewBuffer(jsonBody))
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
		jsonBody := makeJSONRequestBody("https://example.com")

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/new", bytes.NewBuffer(jsonBody))
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 500 {
			t.Errorf("Expected status code 500, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Body.String() != expectedErrorBody {
			t.Errorf("Expected body %s, but got %s", expectedErrorBody, httpRecorder.Body.String())
		}
	})

	t.Run("ShouldReturnAmbiguousOnStoreError", func(t *testing.T) {
		expectedErrorBody := makeExpectedErrorBody("Something went wrong")

		shortLinkRepository.StoreFunc = func(shortLink model.ShortLink) error {
			return errors.New("Some error")
		}
		jsonBody := makeJSONRequestBody("https://example.com")

		httpRecorder := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/new", bytes.NewBuffer(jsonBody))
		webApp.Router.ServeHTTP(httpRecorder, req)

		if httpRecorder.Code != 500 {
			t.Errorf("Expected status code 500, but got %d", httpRecorder.Code)
		}

		if httpRecorder.Body.String() != expectedErrorBody {
			t.Errorf("Expected body %s, but got %s", expectedErrorBody, httpRecorder.Body.String())
		}
	})
}
