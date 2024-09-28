package api

import (
	"context"
	"database/sql"
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
	planAPI.POST("/copy_plan/:fromYear/:toYear", api.copyPlan)
	planAPI.GET("", api.getPlanYears)

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
		return
	}

	plans, err := dbtx.GetPlansByYear(ctx, int32(year))
	if err != nil {
		log.Error("Error fetching plans by year: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
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
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}

func (api *Api) copyPlan(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	fromYearStr := c.Param("fromYear")
	toYearStr := c.Param("toYear")

	// parse years to int
	fromYear, err := strconv.Atoi(fromYearStr)
	if err != nil {
		log.Error("Cant convert fromYear to string")
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing fromYear")
		return
	}
	toYear, err := strconv.Atoi(toYearStr)
	if err != nil {
		log.Error("Cant convert toYear to string")
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing toYear")
		return
	}

	fromPlans, err := dbtx.GetPlansByYear(ctx, int32(fromYear))
	if err != nil {
		log.Error("Error fetching plans by fromYear: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
		return
	}

	// Check that plans for toYear not already exists
	_, err = dbtx.GetPlansByYear(ctx, int32(toYear))
	if err != sql.ErrNoRows {
		log.Error("toYear plans already exists: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "toYear plans already exists")
		return
	}

	var newPlans []db.InsertCopiedPlansParams
	// Map the from plans to newPlans struct
	for _, p := range fromPlans {
		newPlans = append(newPlans, db.InsertCopiedPlansParams{
			Year:            int32(toYear),
			CategoryID:      p.CategoryID,
			JanuaryAmount:   p.JanuaryAmount,
			FebruaryAmount:  p.FebruaryAmount,
			MarchAmount:     p.MarchAmount,
			AprilAmount:     p.AprilAmount,
			MayAmount:       p.MayAmount,
			JuneAmount:      p.JuneAmount,
			JulyAmount:      p.JulyAmount,
			AugustAmount:    p.AugustAmount,
			SeptemberAmount: p.SeptemberAmount,
			OctoberAmount:   p.OctoberAmount,
			NovemberAmount:  p.NovemberAmount,
			DecemberAmount:  p.DecemberAmount,
			Type:            p.Type,
		})
	}

	addedPlans, err := dbtx.InsertCopiedPlans(ctx, newPlans)

	ReturnResponseData(c, http.StatusOK, "OK", addedPlans)
}

func (api *Api) getPlanYears(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	// Check that plans for toYear not already exists
	planYears, err := dbtx.GetPlanYearsOrderedAsc(ctx)
	if err != nil {
		log.Error("Error perform GetPlanYearsOrderedAsc: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Internal Server Error")
		return
	}

	ReturnResponseData(c, http.StatusOK, "OK", planYears)
}
