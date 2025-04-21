package api

import (
	"kidl-api/pkg/db"
	"kidl-api/pkg/models"
	"kidl-api/pkg/utils"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/loopfz/gadgeto/tonic"
	swagger "github.com/num30/gin-swagger-ui"
	log "github.com/sirupsen/logrus"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type Api struct {
	db   *db.Queries
	conn *pgxpool.Pool
}

func StartAPI() {
	time.Sleep(4 * time.Second)
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

	// Tonic error hook
	tonic.SetErrorHook(models.KIDLErrHook)
	r := gin.New()

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
	r.Use(gin.Logger())

	// Fizz initiate
	f := fizz.NewFromEngine(r)

	f.Generator().OverrideDataType(reflect.TypeOf(&pgtype.Text{}), "string", "string")
	f.Generator().OverrideDataType(reflect.TypeOf(&pgtype.Int4{}), "int", "int32")
	f.Generator().OverrideDataType(reflect.TypeOf(&pgtype.Int8{}), "int", "int64")
	f.Generator().OverrideDataType(reflect.TypeOf(&pgtype.Bool{}), "bool", "true")
	f.Generator().OverrideDataType(reflect.TypeOf(&pgtype.Float8{}), "number", "float64")
	f.Generator().OverrideDataType(reflect.TypeOf(&pgtype.Float4{}), "number", "float32")

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

	api.initializeRouters(f)

	// Add Open API description
	infos := &openapi.Info{
		Title:       "KIDL API",
		Description: "This is my KIDL API",
		Version:     utils.GetEnvOrDefault("VERSION", "develop"),
	}

	// Create an endpoint for openapi.json file
	f.GET("/openapi.json", nil, f.OpenAPI(infos, "json"))

	// Now add a UI handler
	swagger.AddOpenApiUIHandler(r, "swagger", "/openapi.json")

	err = f.Engine().Run(":8080")
	if err != nil {
		log.Fatalf("Error starting API, err: %s", err.Error())
	}
}

func (api *Api) healthCheckEndpoint(c *gin.Context, req *models.HealthCheckParams) (models.HealthCheckReturn, error) {

	version := utils.GetEnvOrDefault("VERSION", "develop")
	return models.HealthCheckReturn{
		Version: version,
	}, nil
	// ReturnResponseData(c, http.StatusOK, "OK", )
}

func (api *Api) initializeRouters(f *fizz.Fizz) {
	v1 := f.Group("/v1", "V1", "Version group")

	v1.GET("/health_check", HealthCheckInfo, tonic.Handler(api.healthCheckEndpoint, http.StatusOK))
	// api.transactionsSubrouter(v1)
	// api.transactionCategoriesSubrouter(v1)
	// api.plansSubrouter(v1)
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
