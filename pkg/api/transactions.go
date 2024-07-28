package api

import (
	"context"
	"kidl-api/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (api *Api) transactionsSubrouter(r *gin.RouterGroup) {
	transactionAPI := r.Group("/transaction")

	transactionAPI.POST("", api.uploadTransactions)

}

func (api *Api) uploadTransactions(c *gin.Context) {

	ctx := context.Background()
	dbtx := c.MustGet("dbtx").(*db.Queries)

	if err := c.BindJSON(&parameters); err != nil {
		log.Error("Error parsing request data: ", err)
		ReturnResponseError(c, http.StatusBadRequest, "Error", "Error parsing request")
		return
	}

}
