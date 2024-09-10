package application

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type INewErrorMsg interface {
	string | interface{}
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(Code int, msg INewErrorMsg) Error {
	message := fmt.Sprintf("%v", msg)
	return Error{
		Code:    Code,
		Message: message,
	}
}

func HandleError(c *gin.Context, error Error) {
	err := error
	c.JSON(err.Code, err)
	c.Abort()

}
