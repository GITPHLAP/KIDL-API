package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type MessageCodeData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ResponseData struct {
	Meta MessageCodeData `json:"meta"`
	Data interface{}     `json:"data"`
}
type ResponseError struct {
	Meta MessageCodeError `json:"meta"`
}

type MessageCodeError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

/*
Takes in a response code, a message and the return data to frontend
*/
func ReturnResponseData(c *gin.Context, responseCode int, Message string, data interface{}) {

	tx := c.MustGet("tx").(pgx.Tx)
	err := tx.Commit(c)
	if err != nil {
		var response ResponseError
		response.Meta.Code = http.StatusInternalServerError
		response.Meta.Message = "Internal Server Error"
		response.Meta.Description = "Contact an Administrator"
		c.JSON(http.StatusInternalServerError, response)
	} else {
		var response ResponseData
		response.Meta.Code = responseCode
		response.Meta.Message = Message
		response.Data = data
		c.JSON(responseCode, response)
	}
}

/*
Takes in a response code, an error message and a description of said error message
*/
func ReturnResponseError(c *gin.Context, responseCode int, errorMessage string, description string, bodyResponseCode ...int) {
	var response ResponseError

	tx := c.MustGet("tx").(pgx.Tx)
	err := tx.Rollback(c)
	if err != nil && err.Error() != "tx is closed" {
		response.Meta.Code = http.StatusInternalServerError
		response.Meta.Message = "Internal Server Error"
		response.Meta.Description = "Contact an Administrator"
	} else {
		response.Meta.Code = responseCode
		response.Meta.Message = errorMessage
		response.Meta.Description = description
	}

	if len(bodyResponseCode) > 0 {
		response.Meta.Code = bodyResponseCode[0]
	} else {
		response.Meta.Code = responseCode
	}

	c.JSON(responseCode, response)
}

/*
Takes in a response code, an error message and a description of said error message
*/
func AbortWithError(c *gin.Context, responseCode int, errorMessage string, description string, bodyResponseCode ...int) {
	var response ResponseError

	qtx, exists := c.Get("tx")
	if !exists {
		response.Meta.Code = http.StatusInternalServerError
		response.Meta.Message = "Internal Server Error"
		response.Meta.Description = "Contact an Administrator"
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	tx := qtx.(pgx.Tx)
	err := tx.Rollback(c)
	if err != nil {
		response.Meta.Code = http.StatusInternalServerError
		response.Meta.Message = "Internal Server Error"
		response.Meta.Description = "Contact an Administrator"
	} else {
		response.Meta.Code = responseCode
		response.Meta.Message = errorMessage
		response.Meta.Description = description
	}
	if len(bodyResponseCode) > 0 {
		response.Meta.Code = bodyResponseCode[0]
	} else {
		response.Meta.Code = responseCode
	}
	c.AbortWithStatusJSON(responseCode, response)
}
