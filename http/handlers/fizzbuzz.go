package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"test-lbc/http/models"
	"test-lbc/pkg"
	fModels "test-lbc/pkg/models"
	"test-lbc/prometheus"

	"github.com/gin-gonic/gin"
)

type FizzBuzzService interface {
	Run(params fModels.FizzBuzzParams) ([]string, error)
	GetMostRequested() (*fModels.FizzBuzzStats, error)
}

var serviceFactory = func(db *sql.DB) FizzBuzzService {
	return pkg.NewFizzBuzzService(db)
}

func FizzBuzzRun(c *gin.Context, db *sql.DB) {
	prometheus.IncRequest("run")
	params, errMes := getFizzBuzzParams(c)
	if len(errMes) > 0 {
		prometheus.IncStats("run", "error")
		log.Println("failed to run fizzbuzz:")
		for _, err := range errMes {
			log.Println("\t", err)
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ResponseError{
			Errors: errMes,
		})
		return
	}

	result, err := serviceFactory(db).Run(*params)
	if err != nil {
		log.Printf("failed to save stats: %v", err)
		prometheus.IncStats("run", "error_on_stat_save")
	} else {
		prometheus.IncStats("run", "success")
	}

	c.JSON(http.StatusOK, result)
}

func FizzBuzzStats(c *gin.Context, db *sql.DB) {
	prometheus.IncRequest("stats")
	mostRequested, err := serviceFactory(db).GetMostRequested()
	if err != nil {
		prometheus.IncStats("stats", "error")
		log.Printf("failed to retrieve fizzbuzz stats: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ResponseError{
			Errors: []string{err.Error()},
		})
		return
	}

	prometheus.IncStats("stats", "success")
	c.JSON(http.StatusOK, mostRequested)
}

func getFizzBuzzParams(c *gin.Context) (*fModels.FizzBuzzParams, []string) {
	var (
		int1Str  = c.Query("int1")
		int2Str  = c.Query("int2")
		limitStr = c.Query("limit")
		str1     = c.Query("str1")
		str2     = c.Query("str2")

		errMes []string
	)

	if int1Str == "" || int2Str == "" || limitStr == "" || str1 == "" || str2 == "" {
		return nil, []string{"int1, int2, limit, str1 and str2 are all mandatory"}
	}
	int1, err := strconv.Atoi(int1Str)
	if err != nil {
		errMes = append(errMes, "int1 err: "+err.Error())
	}

	int2, err := strconv.Atoi(int2Str)
	if err != nil {
		errMes = append(errMes, "int2 err: "+err.Error())
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		errMes = append(errMes, "limit err: "+err.Error())
	}

	if limit <= 0 {
		errMes = append(errMes, "limit must be greater than 0")
	}

	return &fModels.FizzBuzzParams{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}, errMes
}
