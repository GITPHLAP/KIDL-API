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
	"github.com/juju/errors"

	"github.com/loopfz/gadgeto/tonic"
	log "github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
)

func (api *Api) overviewSubrouter(f *fizz.RouterGroup) {
	overviewAPI := f.Group("/overview", "Overview", "Group of all overview endpoints")

	overviewAPI.GET("/:year", nil, tonic.Handler(api.getOverviews, http.StatusOK))

	// overviewAPI.GET("/template/:year/:type", api.getTemplateData)

}

func (api *Api) getOverviews(c *gin.Context, params *models.GetOverviewsParams) (result models.OverviewReturn, err error) {
	var totalMonthData []models.MonthData

	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	// parse year to int
	year := params.Year

	// Converts the query string array to array of months in type int
	months := params.Months

	// Validate the months
	err = monthValidator(months)
	if err != nil {
		return nil, errors.NotValidf("Bindings %s", err)
	}

	overviewDataByte, err := dbtx.GetOverviewByYearNMonth(ctx, db.GetOverviewByYearNMonthParams{
		Year:   int32(year),
		Months: months,
	})
	if err != nil {
		return nil, errors.Errorf("Cant get overview by year and month: %s", err)
	}

	err = json.Unmarshal(overviewDataByte, &result)
	if err != nil {
		return nil, errors.Errorf("Cant unmarshal overview data: %s", err)
	}

	// Fetch the "13th month" Total for the year
	totalOverviewDataByte, err := dbtx.Get13thMonthByYear(ctx, int32(year))
	if err != nil {
		return nil, errors.Errorf("Cant get total overview by year: %s", err)
	}

	// Unmarshal the data to MonthData struct
	err = json.Unmarshal(totalOverviewDataByte.Calculation, &totalMonthData)
	if err != nil {
		return nil, errors.Errorf("Cant unmarshal overview data: %s", err)
	}
	// Create a MonthWrapper to append to result slice of MonthWrapper
	totalMonthWrapper := models.MonthWrapper{}
	totalMonthWrapper["13"] = totalMonthData

	// Append the "13" to rest of the overview
	result[0].Months = append(result[0].Months, totalMonthWrapper)

	return result, nil
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

// getTemplateOverviews godoc
//
//	@Summary		Get data for the Nordea template
//	@Description	get data by year and type
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			year	path		int	true	"Year"
//	@Param			type	path		db.TransType	true	"Type"
//	@Success		200	{object}	api.ResponseData{data=db.YearlyOverviewByCategory}
//	@Failure		400	{object}	api.ResponseError
//	@Failure		500	{object}	api.ResponseError
//	@Router			/overview/template/{year}/{type} [get]
func (api *Api) getTemplateData(c *gin.Context) {
	// var totalMonthData []MonthData

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

	typeStr := c.Param("type")

	overviewData, err := dbtx.GetYearOverviewByYearAndType(ctx, db.GetYearOverviewByYearAndTypeParams{
		Year: int32(year),
		Type: db.TransType(typeStr),
	})
	if err != nil {
		log.Error("Cant get year-overview by year and type: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", overviewData)
}
