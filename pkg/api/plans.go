package api

import (
	"context"
	"kidl-api/pkg/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (api *Api) plansSubrouter(r *gin.RouterGroup) {
	planAPI := r.Group("/plan")

	planAPI.GET("/:year", api.getPlans)
	planAPI.POST("", api.createPlan)
	planAPI.PATCH("", api.updatePlan)

}

func (api *Api) getPlans(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	yearstr := c.Param("year")

	// parse year to int
	year, err := strconv.Atoi(yearstr)
	if err != nil {
		log.Error("Cant convert year to string")
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing year")
	}

	plans, err := dbtx.GetPlansByYear(ctx, int32(year))
	if err != nil {
		log.Error("Error fetching plans by year: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", plans)
}

func (api *Api) createPlan(c *gin.Context) {
	var parameters db.InsertPlanParams

	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	err := dbtx.InsertPlan(ctx, parameters)
	if err != nil {
		log.Error("Error inserting plan: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}

func (api *Api) updatePlan(c *gin.Context) {
	var parameters db.UpdatePlanByIdParams
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	err := dbtx.UpdatePlanById(ctx, parameters)
	if err != nil {
		log.Error("Error update plan: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}
