package api

import (
	"context"
	"kidl-api/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (api *Api) transactionCategoriesSubrouter(r *gin.RouterGroup) {
	transactionCategoriesAPI := r.Group("/transaction_categories")

	transactionCategoriesAPI.GET("", api.getTransactionCategories)
	transactionCategoriesAPI.POST("/:category", api.createTransactionCategory)

	transactionCategoriesMappingAPI := transactionCategoriesAPI.Group("/mapping")

	transactionCategoriesMappingAPI.POST("", api.createTransactionCategoryMapping)
}

func (api *Api) getTransactionCategories(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	categories, err := dbtx.GetAllCategories(ctx)
	if err != nil {
		log.Error("Error fetching categories: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", categories)
}

func (api *Api) createTransactionCategory(c *gin.Context) {
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	category := c.Param("category")

	// Check if category allready exist
	exist, err := dbtx.AnyCategoryByName(ctx, category)
	if err != nil {
		log.Error("Error check if category exits: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}
	if exist {
		log.Info("User try to create a category which already exists: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Category already exists")
		return
	}

	err = dbtx.CreateCategory(ctx, category)
	if err != nil {
		log.Error("Error creating category: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}

func (api *Api) createTransactionCategoryMapping(c *gin.Context) {
	var parameters db.CreateCategoryMappingParams
	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

	err := dbtx.CreateCategoryMapping(ctx, parameters)
	if err != nil {
		log.Error("Error creating category mapping: ", err)
		ReturnResponseError(c, http.StatusInternalServerError, "Error", "Internal server error")
	}

	ReturnResponseData(c, http.StatusOK, "OK", nil)
}
