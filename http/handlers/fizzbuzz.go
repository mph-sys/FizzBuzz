package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"test-lbc/http/models"
	"test-lbc/pkg"
	fModels "test-lbc/pkg/models"

	"github.com/gin-gonic/gin"
)

func FizzBuzzRun(c *gin.Context, db *sql.DB) {
	params, errMes := getFizzBuzzParams(c)
	if len(errMes) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ResponseError{
			Errors: errMes,
		})
		return
	}
	c.JSON(http.StatusOK, models.ResponseSuccess{
		Data: pkg.NewFizzBuzzService(db).Run(*params),
	})
}

func FizzBuzzStats(c *gin.Context, db *sql.DB) {
	mostRequested, err := pkg.NewFizzBuzzService(db).GetMostRequested()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, models.ResponseError{
			Errors: []string{err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, models.ResponseSuccess{
		Data: mostRequested,
	})
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

	return &fModels.FizzBuzzParams{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}, errMes
}
