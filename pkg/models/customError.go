package models

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/sirupsen/logrus"
)

type CustomError struct {
	Message string
	Custom  bool
}

func (c CustomError) Error() string {
	return c.Message
}

func KIDLErrHook(c *gin.Context, e error) (int, interface{}) {

	errCode, errDesc, errMsg := 500, e.Error(), "Internal Server Error"
	if _, ok := e.(tonic.BindError); ok {
		errCode, errDesc = 400, e.Error()
	} else {
		switch {
		case errors.IsBadRequest(e) || errors.IsNotValid(e) || errors.IsAlreadyExists(e) || errors.IsNotSupported(e) || errors.IsNotAssigned(e) || errors.IsNotProvisioned(e):
			errCode, errDesc, errMsg = 400, e.Error(), "Error"
		case errors.IsForbidden(e):
			errCode, errDesc, errMsg = 403, e.Error(), "Error"
		case errors.IsMethodNotAllowed(e):
			errCode, errDesc, errMsg = 405, e.Error(), "Error"
		case errors.IsNotFound(e) || errors.IsUserNotFound(e):
			errCode, errDesc, errMsg = 404, e.Error(), "Error"
		case errors.IsUnauthorized(e):
			errCode, errDesc, errMsg = 401, e.Error(), "Error"
		case errors.IsNotImplemented(e):
			errCode, errDesc, errMsg = 501, e.Error(), "Error"
		}
	}

	// Write the error log
	logrus.Error(e)

	return errCode, ResponseError{
		Meta: MessageCodeError{
			Code:        errCode,
			Message:     errMsg,
			Description: errDesc,
		},
	}
}
