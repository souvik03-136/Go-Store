package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/utils"
)

/* -------------------------------------------------------------------------- */
/*                            BAD REQUEST ERROR 400                           */
/* -------------------------------------------------------------------------- */
func BadRequest(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusBadRequest

	smerror.Code = errorCode
	smerror.Type = errorType.Unauthorized
	smerror.Message = err

	res.Error = &smerror

	ctx.JSON(errorCode, res)
	ctx.Abort()
}
