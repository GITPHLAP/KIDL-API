package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

func TransactionHandler(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx, err := api.conn.Begin(c)
		if err != nil {
			log.Error("[TransactionHandler] Error creating transaction: ", err)
			AbortWithError(c, http.StatusInternalServerError, "Internal Server Error", "Contact an Administrator")
			return
		}

		qtx := api.db.WithTx(tx)

		c.Set("dbtx", qtx)
		c.Set("tx", tx)

		c.Next()

		err = tx.Rollback(c)
		if err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				log.Error(err)
			}
		} else {
			log.Warn("Rolled back transaction in middleware", c.Request.URL.String(), " ", c.Request.Method)
		}
	}
}
