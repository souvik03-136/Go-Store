package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/utils"
)

/* -------------------------------------------------------------------------- */
/*                                Conflict 409                                */
/* -------------------------------------------------------------------------- */

func Conflict(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusConflict
	smerror.Code = errorCode
	smerror.Type = errorType.conflict
	smerror.Message = err
	res.Error = &smerror
	ctx.JSON(errorCode, res)
	ctx.Abort()
}
