package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/utils"
)

/* -------------------------------------------------------------------------- */
/*                            INTERNAL SERVER ERROR 500                       */
/* -------------------------------------------------------------------------- */
func InternalServer(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusInternalServerError

	smerror.Code = errorCode
	smerror.Type = errorType.server
	smerror.Message = err

	res.Error = &smerror

	ctx.JSON(errorCode, res)
	ctx.Abort()
}
