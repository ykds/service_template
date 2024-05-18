package response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"service_template/internal/errors"
	"service_template/pkg/logger"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty" swaggertype:"object"`
}

func HandleResponse(ctx *gin.Context, err error, req interface{}, data interface{}) {
	if err != nil {
		reqByte, _ := json.Marshal(req)
		logger.Errorf("req params: %+v, err: %+v", string(reqByte), err)
		var e = new(errors.Error)
		if errors.As(err, e) {
			ctx.JSON(200, Response{Code: e.Code(), Message: e.Message()})
		} else {
			ctx.JSON(200, Response{Code: errors.InternalError.Code(), Message: errors.InternalError.Message()})
		}
		return
	}
	ctx.JSON(200, Response{Code: errors.Success.Code(), Message: errors.Success.Message(), Data: data})
}
