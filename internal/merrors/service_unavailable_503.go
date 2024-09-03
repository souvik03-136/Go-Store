package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/utils"
)

func ServiceUnavailable(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusServiceUnavailable
	smerror.Code = errorCode
	smerror.Type = errorType.ServiceUnavailable
	smerror.Message = err

	res.Error = &smerror

	ctx.JSON(errorCode, res)
	ctx.Abort()
}
