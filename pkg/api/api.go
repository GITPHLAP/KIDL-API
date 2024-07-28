package api

import (
	"kidl-api/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Api struct {
	db   *db.Queries
	conn *pgxpool.Pool
}

func StartAPI() {

}

func (api *Api) initializeRouters(r *gin.Engine) {
	v1 := r.Group("/v1")

}
