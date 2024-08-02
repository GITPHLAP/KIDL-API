package api

import (
	"context"
	"kidl-api/pkg/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (api *Api) forecastsSubrouter(r *gin.RouterGroup) {
	forecastAPI := r.Group("/forecast")

	forecastAPI.GET("/:year", api.getForecasts)
	forecastAPI.POST("", api.createForecast)
	forecastAPI.PATCH("", api.updateForecast)

}

func (api *Api) getForecasts(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	yearstr := c.Param("year")

	// parse year to int
	year, err := strconv.Atoi(yearstr)
	if err != nil {
		log.Error("Cant convert year to string")
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing year")
	}

	forecasts, err := dbtx.GetForecastsByYear(ctx, int32(year))
	if err != nil {
		log.Error("Error fetching forecasts by year: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", forecasts)
}

func (api *Api) createForecast(c *gin.Context) {
	var parameters db.InsertForecastParams

	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	err := dbtx.InsertForecast(ctx, parameters)
	if err != nil {
		log.Error("Error inserting forecast: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}

func (api *Api) updateForecast(c *gin.Context) {
	var parameters db.UpdateForecastByIdParams
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	err := dbtx.UpdateForecastById(ctx, parameters)
	if err != nil {
		log.Error("Error update forecast: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}
