package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/utils"
)

/* -------------------------------------------------------------------------- */
/*                                FORBIDDEN 403                               */
/* -------------------------------------------------------------------------- */
func Forbidden(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusForbidden
	smerror.Code = errorCode
	smerror.Type = errorType.Forbidden
	smerror.Message = err
	res.Error = &smerror
	ctx.JSON(errorCode, res)
	ctx.Abort()
}
