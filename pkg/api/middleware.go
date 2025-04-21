package api

import (
	"kidl-api/pkg/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

func TransactionHandler(api *Api) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx, err := api.conn.Begin(c)
		if err != nil {
			log.Error()
			c.AbortWithError(http.StatusInternalServerError, errors.Errorf("Error creating transaction: %s", err))
			return
		}

		qtx := api.db.WithTx(tx)

		c.Set("dbtx", qtx)
		c.Set("tx", tx)

		c.Next()

		if c.Writer.Status() == http.StatusOK {
			log.Debug("Committing transaction")
			if err := tx.Commit(c); err != nil {
				log.Error("Error on committing transaction: ", err)
			}
		} else {
			log.Debug("Rolling back transaction in middleware")
			tx.Rollback(c)
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
