package merrors

import (
	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/utils"
)

/* -------------------------------------------------------------------------- */
/*                              DOWNSTREAM ERROR                              */
/* -------------------------------------------------------------------------- */
func Downstream(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := 550
	smerror.Code = errorCode
	smerror.Type = errorType.Downstream
	smerror.Message = err
	res.Error = &smerror
	ctx.JSON(errorCode, res)
	ctx.Abort()
}
