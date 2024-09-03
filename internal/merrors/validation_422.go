package merrors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souvik03-136/Go-Store/internal/utils"
)

/* -------------------------------------------------------------------------- */
/*                            VALIDATION ERROR 422                            */
/* -------------------------------------------------------------------------- */
func Validation(ctx *gin.Context, err string) {
	var res utils.BaseResponse
	var smerror utils.Error
	errorCode := http.StatusUnprocessableEntity

	smerror.Code = errorCode
	smerror.Type = errorType.validation
	smerror.Message = err

	res.Error = &smerror

	ctx.JSON(errorCode, res)
	ctx.Abort()
}
