package driver_test

import (
	"errors"
	"testing"

	"github.com/JakubBizewski/jakubme_links/domain/model"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
	"github.com/JakubBizewski/jakubme_links/mocks"
)

type ShortLinkGeneratorTestSuite struct {
	mockRepository *mocks.MockShortLinkRepository
	service        *driver.ShortLinkService
}

func (suite *ShortLinkGeneratorTestSuite) SetupTest() {
	suite.mockRepository = &mocks.MockShortLinkRepository{}
	suite.service = driver.CreateShortLinkService(suite.mockRepository)
}

func TestShortLinkGeneratorTestSuite(t *testing.T) {
	suite := new(ShortLinkGeneratorTestSuite)
	suite.SetupTest()

	generateAttemptsTestCases := []struct {
		name     string
		attempts int
	}{
		{"StoreUniqueOnFirstTry", 1},
		{"StoreUniqueOnSecondTry", 2},
		{"StoreUniqueOnNinethTry", 9},
		{"StoreUniqueOnTenthTry", 10},
	}

	for _, testCase := range generateAttemptsTestCases {
		t.Run(testCase.name, func(t *testing.T) {
			attempts := 0
			suite.mockRepository.StoreFunc = func(shortLink model.ShortLink) error {
				attempts++
				if attempts < testCase.attempts {
					return driven.ErrShortCodeAlreadyExists
				}
				return nil
			}

			_, err := suite.service.GenerateShortLink("https://example.com")

			if err != nil {
				t.Errorf("Expected no error, but got %v", err)
			}

			if attempts != testCase.attempts {
				t.Errorf("Expected %d attempts, but got %d", testCase.attempts, attempts)
			}
		})
	}

	t.Run("ErrShortCodeGenerationFailedAfter10StoreFails", func(t *testing.T) {
		attempts := 0
		suite.mockRepository.StoreFunc = func(shortLink model.ShortLink) error {
			attempts++

			return driven.ErrShortCodeAlreadyExists
		}

		shortCode, err := suite.service.GenerateShortLink("https://example.com")

		if shortCode != "" {
			t.Errorf("Expected short code %s, but got %s", "", shortCode)
		}

		if attempts != 10 {
			t.Errorf("Expected 10 attempts, but got %d", attempts)
		}

		if err != driver.ErrShortCodeGenerationFailed {
			t.Errorf("Expected error %v, but got %v", driver.ErrShortCodeGenerationFailed, err)
		}
	})

	t.Run("ErrorWithoutRetriesOnNonShortCodeAlreadyExistsError", func(t *testing.T) {
		expectedError := errors.New("Some error")
		storeCalls := 0

		suite.mockRepository.StoreFunc = func(shortLink model.ShortLink) error {
			storeCalls++
			return expectedError
		}

		shortCode, err := suite.service.GenerateShortLink("https://example.com")

		if shortCode != "" {
			t.Errorf("Expected short code %s, but got %s", "", shortCode)
		}

		if err != expectedError {
			t.Errorf("Expected error %v, but got %v", expectedError, err)
		}

		if storeCalls != 1 {
			t.Errorf("Expected 1 call to store, but got %d", storeCalls)
		}
	})
}
