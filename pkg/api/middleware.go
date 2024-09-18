package api

import (
	"errors"
	"kidl-api/pkg/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func APIStatistics(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Time the request starts
		t := time.Now()

		//Before request is processed

		c.Next()

		//After request is processed

		url := c.Request.URL.Path
		status := c.Writer.Status()
		method := c.Request.Method
		size := c.Writer.Size()
		ip := c.ClientIP()
		elapsedTime := time.Since(t)

		stats := db.InsertApiStatisticParams{
			Url:                url,
			StatusCode:         int32(status),
			HttpMethod:         method,
			SizeBytes:          int64(size),
			ClientIp:           ip,
			ElapsedTimeSeconds: pgtype.Float8{Float64: elapsedTime.Seconds(), Valid: true},
		}

		err := api.db.InsertApiStatistic(c, stats)
		if err != nil {
			log.Warning("Error inserting API statistics: ", err.Error())
		}

	}
}
