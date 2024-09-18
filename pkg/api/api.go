package api

import (
	"kidl-api/pkg/db"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Api struct {
	db   *db.Queries
	conn *pgxpool.Pool
}

func StartAPI() {
	log.Info("Starting API")
	// Setting up db connection
	db, conn, err := db.InitConn()
	if err != nil {
		log.Fatal("[NBI2] Failed to connect to database: ", err)
	}

	api := Api{
		db:   db,
		conn: conn,
	}

	r := gin.Default()

	// Cors options
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:  []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	r.Use(TransactionHandler(&api))
	r.Use(APIStatistics(&api))
	r.Use(gin.Recovery())

	// Validations of pgtypes
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		log.Println("Registering new types")
		v.RegisterCustomTypeFunc(ValidatePgtypeText, pgtype.Text{})
		v.RegisterCustomTypeFunc(ValidatePgtypeInt4, pgtype.Int4{})
		v.RegisterCustomTypeFunc(ValidatePgtypeInt8, pgtype.Int8{})
		v.RegisterCustomTypeFunc(ValidatePgtypeBool, pgtype.Bool{})
		v.RegisterCustomTypeFunc(ValidatepgtypeFloat8, pgtype.Float8{})
		v.RegisterCustomTypeFunc(ValidatepgtypeFloat4, pgtype.Float4{})
	} else {
		log.Println("This failed for some reason")
	}

	api.initializeRouters(r)
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting API, err: %s", err.Error())
	}
}

func (api *Api) healthCheckEndpoint(c *gin.Context) {
	type HealthCheckReturn struct {
		Version string `json:"version"`
	}
	version := os.Getenv("VERSION")
	ReturnResponseData(c, http.StatusOK, "OK", HealthCheckReturn{
		Version: version,
	})
}

func (api *Api) initializeRouters(r *gin.Engine) {
	v1 := r.Group("/v1")

	v1.GET("/health_check", api.healthCheckEndpoint)
	api.transactionsSubrouter(v1)
	api.transactionCategoriesSubrouter(v1)
	api.plansSubrouter(v1)
	api.overviewSubrouter(v1)

}

func ValidatePgtypeText(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(pgtype.Text{}) {
		if !field.Interface().(pgtype.Text).Valid {
			return nil
		}
		return field.Interface().(pgtype.Text).String
	}

	return nil
}

func ValidatePgtypeInt4(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(pgtype.Int4{}) {
		if !field.Interface().(pgtype.Int4).Valid {
			return nil
		}
		return field.Interface().(pgtype.Int4).Int32
	}

	return nil
}

func ValidatePgtypeInt8(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(pgtype.Int8{}) {
		if !field.Interface().(pgtype.Int8).Valid {
			return nil
		}
		return field.Interface().(pgtype.Int8).Int64
	}

	return nil
}

func ValidatePgtypeBool(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(pgtype.Bool{}) {
		if !field.Interface().(pgtype.Bool).Valid {
			return nil
		}
		return field.Interface().(pgtype.Bool).Bool
	}

	return nil
}

func ValidatepgtypeFloat8(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(pgtype.Float8{}) {
		if !field.Interface().(pgtype.Float8).Valid {
			return nil
		}
		return field.Interface().(pgtype.Float8).Float64
	}

	return nil
}

func ValidatepgtypeFloat4(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(pgtype.Float4{}) {
		if !field.Interface().(pgtype.Float4).Valid {
			return nil
		}
		return field.Interface().(pgtype.Float4).Float32
	}

	return nil
}
