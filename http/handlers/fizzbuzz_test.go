package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	fModels "test-lbc/pkg/models"

	"github.com/gin-gonic/gin"
)

// MockService implements FizzBuzzService for testing purposes
type MockService struct {
	RunFunc              func(params fModels.FizzBuzzParams) ([]string, error)
	GetMostRequestedFunc func() (*fModels.FizzBuzzStats, error)
}

func (m *MockService) Run(params fModels.FizzBuzzParams) ([]string, error) {
	if m.RunFunc != nil {
		return m.RunFunc(params)
	}
	return nil, nil
}

func (m *MockService) GetMostRequested() (*fModels.FizzBuzzStats, error) {
	if m.GetMostRequestedFunc != nil {
		return m.GetMostRequestedFunc()
	}
	return nil, nil
}

func TestFizzBuzzRun(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Save original factory and restore after test
	origFactory := serviceFactory
	defer func() { serviceFactory = origFactory }()

	t.Run("Missing Parameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/fizzbuzz/run", nil)

		FizzBuzzRun(c, nil)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Invalid Parameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// int1 is invalid
		c.Request, _ = http.NewRequest("POST", "/fizzbuzz/run?int1=abc&int2=5&limit=100&str1=f&str2=b", nil)

		FizzBuzzRun(c, nil)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Success", func(t *testing.T) {
		expectedResp := []string{"1", "2", "fizz"}
		mockSvc := &MockService{
			RunFunc: func(params fModels.FizzBuzzParams) ([]string, error) {
				return expectedResp, nil
			},
		}
		serviceFactory = func(db *sql.DB) FizzBuzzService {
			return mockSvc
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/fizzbuzz/run?int1=3&int2=5&limit=3&str1=fizz&str2=buzz", nil)

		FizzBuzzRun(c, nil)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var resp []string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if len(resp) != 3 {
			t.Errorf("Expected 3 items, got %d", len(resp))
		}
	})
}

func TestFizzBuzzStats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	origFactory := serviceFactory
	defer func() { serviceFactory = origFactory }()

	t.Run("Success", func(t *testing.T) {
		mockSvc := &MockService{
			GetMostRequestedFunc: func() (*fModels.FizzBuzzStats, error) {
				return &fModels.FizzBuzzStats{Int1: 3, Int2: 5, Limit: 100, Str1: "f", Str2: "b", Hits: 10}, nil
			},
		}
		serviceFactory = func(db *sql.DB) FizzBuzzService {
			return mockSvc
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/fizzbuzz/stats/most-requested", nil)

		FizzBuzzStats(c, nil)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := &MockService{
			GetMostRequestedFunc: func() (*fModels.FizzBuzzStats, error) {
				return nil, errors.New("database error")
			},
		}
		serviceFactory = func(db *sql.DB) FizzBuzzService {
			return mockSvc
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/fizzbuzz/stats/most-requested", nil)

		FizzBuzzStats(c, nil)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
	})
}
