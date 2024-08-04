package api

import (
	"context"
	"encoding/json"
	"fmt"
	"kidl-api/pkg/db"
	"kidl-api/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	log "github.com/sirupsen/logrus"
)

func (api *Api) overviewSubrouter(r *gin.RouterGroup) {
	forecastAPI := r.Group("/overview")

	forecastAPI.GET("/:year", api.getOverviews)

}

type OverviewReturn []struct {
	Year   int            `json:"year"`
	Months []MonthWrapper `json:"months"`
}

type MonthWrapper map[string][]MonthData

type TypeData struct {
	ActualSum   pgtype.Float8 `json:"actual_sum"`
	Category    pgtype.Text   `json:"category"`
	ForecastSum pgtype.Float8 `json:"forecast_sum"`
}
type TypeTotalData struct {
	TotalActual   pgtype.Float8 `json:"total_actual,omitempty"`
	TotalForecast pgtype.Float8 `json:"total_forecast,omitempty"`
}

type MonthData struct {
	Expenses      []TypeData     `json:"expenses,omitempty"`
	ExpensesTotal *TypeTotalData `json:"expenses_total,omitempty"`
	Income        []TypeData     `json:"income,omitempty"`
	IncomeTotal   *TypeTotalData `json:"income_total,omitempty"`
}

func (api *Api) getOverviews(c *gin.Context) {
	var result OverviewReturn
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	yearstr := c.Param("year")

	// parse year to int
	year, err := strconv.Atoi(yearstr)
	if err != nil {
		log.Error("Cant convert year to int: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing year")
		return
	}

	monthsStr, ok := c.GetQueryArray("months")
	if !ok {
		log.Error("No months is entered: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error finding months in query")
		return
	}

	// Converts the query string array to array of months in type int
	months, err := convertStringArrayToIntArr(monthsStr)
	if err != nil {
		log.Error("Cant convert months to int: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing months")
		return
	}

	// Validate the months
	err = monthValidator(months)
	if err != nil {
		log.Error(err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", err.Error())
		return
	}

	overviewDataByte, err := dbtx.GetOverviewByYearNMonth(ctx, db.GetOverviewByYearNMonthParams{
		Year:   int32(year),
		Months: months,
	})
	if err != nil {
		log.Error("Cant get overview by year and month: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	log.Infof("%+v\n", string(overviewDataByte))

	err = json.Unmarshal(overviewDataByte, &result)
	if err != nil {
		log.Error("Cant unmarshal overview data: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", result)
}

func convertStringArrayToIntArr(input []string) ([]int32, error) {
	var result []int32
	for _, s := range input {
		log.Infoln("s: ", s)
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Error("Cant convert string to int")
			return nil, err
		}
		result = append(result, int32(i))
	}
	return result, nil
}

func monthValidator(months []int32) error {
	for _, m := range months {
		// from 1 to 12 is valid months
		if m < 1 || m > 12 {
			msg := fmt.Sprintf("%d is not a valid month", m)
			return models.CustomError{Message: msg, Custom: true}
		}
	}
	return nil
}
