package pkg

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"test-lbc/pkg/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFizzBuzzService_Run(t *testing.T) {
	testCases := []struct {
		name     string
		params   models.FizzBuzzParams
		expected []string
	}{
		{
			name:     "Standard FizzBuzz",
			params:   models.FizzBuzzParams{Int1: 3, Int2: 5, Limit: 15, Str1: "fizz", Str2: "buzz"},
			expected: []string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz"},
		},
		{
			name:     "Both int1 and int2 are zero",
			params:   models.FizzBuzzParams{Int1: 0, Int2: 0, Limit: 3, Str1: "fizz", Str2: "buzz"},
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "Only int1 is zero",
			params:   models.FizzBuzzParams{Int1: 0, Int2: 3, Limit: 5, Str1: "fizz", Str2: "buzz"},
			expected: []string{"1", "2", "buzz", "4", "5"},
		},
		{
			name:     "Only int2 is zero",
			params:   models.FizzBuzzParams{Int1: 3, Int2: 0, Limit: 5, Str1: "fizz", Str2: "buzz"},
			expected: []string{"1", "2", "fizz", "4", "5"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			// We expect the stats query to be executed for all cases except when both ints are 0
			if tc.params.Int1 != 0 || tc.params.Int2 != 0 {
				mock.ExpectExec("INSERT INTO `stats`").
					WithArgs(tc.params.Int1, tc.params.Int2, tc.params.Limit, tc.params.Str1, tc.params.Str2, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			service := NewFizzBuzzService(db)
			result, err := service.Run(tc.params)
			if err != nil {
				t.Errorf("error while running: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}

	t.Run("DB Error on incStats", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		params := models.FizzBuzzParams{Int1: 3, Int2: 5, Limit: 3, Str1: "fizz", Str2: "buzz"}
		expectedResult := []string{"1", "2", "fizz"}

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `stats`")).
			WithArgs(params.Int1, params.Int2, params.Limit, params.Str1, params.Str2, 1).
			WillReturnError(errors.New("db error"))

		service := NewFizzBuzzService(db)
		result, err := service.Run(params)
		if err != nil {
			t.Logf("error expected: %v", err)
		}

		if !reflect.DeepEqual(result, expectedResult) {
			t.Errorf("expected result %v even with db error, got %v", expectedResult, result)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestFizzBuzzService_GetMostRequested(t *testing.T) {
	query := regexp.QuoteMeta("SELECT `int1`,`int2`,`limit`,`str1`,`str2`,`hits` FROM `stats` ORDER BY `hits` desc LIMIT 1")

	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		service := NewFizzBuzzService(db)

		expectedStats := &models.FizzBuzzStats{
			Int1: 3, Int2: 5, Limit: 100, Str1: "fizz", Str2: "buzz", Hits: 20,
		}

		rows := sqlmock.NewRows([]string{"int1", "int2", "limit", "str1", "str2", "hits"}).
			AddRow(expectedStats.Int1, expectedStats.Int2, expectedStats.Limit, expectedStats.Str1, expectedStats.Str2, expectedStats.Hits)

		mock.ExpectQuery(query).WillReturnRows(rows)

		stats, err := service.GetMostRequested()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(stats, expectedStats) {
			t.Errorf("expected stats %v, got %v", expectedStats, stats)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("No rows found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		service := NewFizzBuzzService(db)

		rows := sqlmock.NewRows([]string{"int1", "int2", "limit", "str1", "str2", "hits"})
		mock.ExpectQuery(query).WillReturnRows(rows)

		stats, err := service.GetMostRequested()

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if stats != nil {
			t.Errorf("expected nil stats, got %v", stats)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		service := NewFizzBuzzService(db)

		dbErr := errors.New("query failed")
		mock.ExpectQuery(query).WillReturnError(dbErr)

		stats, err := service.GetMostRequested()

		if stats != nil {
			t.Errorf("expected nil stats on error, got %v", stats)
		}
		if err == nil {
			t.Errorf("expected an error, but got nil")
		} else if err.Error() != fmt.Sprintf("failed to query most requested: %s", dbErr.Error()) {
			t.Errorf("unexpected error message: %s", err.Error())
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()
		service := NewFizzBuzzService(db)

		rows := sqlmock.NewRows([]string{"int1", "int2", "limit", "str1", "str2", "hits"}).
			AddRow(3, 5, 100, "fizz", "buzz", "not-an-integer") // Invalid type for hits
		mock.ExpectQuery(query).WillReturnRows(rows)

		stats, err := service.GetMostRequested()

		if stats != nil {
			t.Errorf("expected nil stats on scan error, got %v", stats)
		}
		if err == nil {
			t.Errorf("expected a scan error, but got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
