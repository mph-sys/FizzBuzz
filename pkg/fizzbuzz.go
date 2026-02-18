package pkg

import (
	"database/sql"
	"fmt"
	"strconv"
	"test-lbc/pkg/models"
)

type FizzBuzzService struct {
	db *sql.DB
}

func NewFizzBuzzService(db *sql.DB) FizzBuzzService {
	return FizzBuzzService{
		db: db,
	}
}

func (s FizzBuzzService) Run(params models.FizzBuzzParams) ([]string, error) {
	var (
		result       = make([]string, params.Limit)
		currentValue = 1
	)

	if params.Int1 == 0 && params.Int2 == 0 {
		for i := range result {
			result[i] = strconv.Itoa(currentValue)
			currentValue++
		}
		return result, nil
	}

	int1TestValue, int2TestValue := params.Int1, params.Int2
	if int1TestValue == 0 {
		int1TestValue = params.Limit + 1
	}
	if int2TestValue == 0 {
		int2TestValue = params.Limit + 1
	}

	for i := range result {
		// switch {
		// case currentValue%(int1TestValue*int2TestValue) == 0:
		// 	result[i] = params.Str1 + params.Str2
		// case currentValue%int1TestValue == 0:
		// 	result[i] = params.Str1
		// case currentValue%int2TestValue == 0:
		// 	result[i] = params.Str2
		// default:
		// 	result[i] = strconv.Itoa(currentValue)
		// }
		if currentValue%int1TestValue == 0 {
			result[i] = params.Str1
		}
		if currentValue%int2TestValue == 0 {
			result[i] += params.Str2
		}
		if result[i] == "" {
			result[i] = strconv.Itoa(currentValue)
		}
		currentValue++
	}

	return result, s.incStats(params)
}

func (s *FizzBuzzService) incStats(params models.FizzBuzzParams) error {
	_, err := s.db.Exec("INSERT INTO `stats` (`int1`,`int2`,`limit`,`str1`,`str2`,`hits`) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE `hits` = `hits`+1", params.Int1, params.Int2, params.Limit, params.Str1, params.Str2, 1)
	if err != nil {
		return fmt.Errorf("failed to save request: %v", err)
	}

	return nil
}

func (s FizzBuzzService) GetMostRequested() (*models.FizzBuzzStats, error) {
	rows, err := s.db.Query("SELECT `int1`,`int2`,`limit`,`str1`,`str2`,`hits` FROM `stats` ORDER BY `hits` desc LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("failed to query most requested: %v", err)
	}
	defer rows.Close()

	var mostRequested *models.FizzBuzzStats
	for rows.Next() {
		var (
			int1, int2, limit, hits int
			str1, str2              string
		)
		if err := rows.Scan(&int1, &int2, &limit, &str1, &str2, &hits); err != nil {
			return nil, fmt.Errorf("failed to scan most requested: %v", err)
		}
		mostRequested = &models.FizzBuzzStats{
			Int1:  int1,
			Int2:  int2,
			Limit: limit,
			Str1:  str1,
			Str2:  str2,
			Hits:  hits,
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %v", err)
	}

	return mostRequested, nil
}
