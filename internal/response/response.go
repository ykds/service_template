package response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"service_template/internal/errors"
	"service_template/pkg/logger"
)

type Response struct {
	Code int         `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data,omitempty"`
}

func HandleResponse(ctx *gin.Context, err error, req interface{}, data interface{}) {
	if err != nil {
		reqByte, _ := json.Marshal(req)
		logger.Errorf("req params: %+v, err: %+v", string(reqByte), err)
		var e = new(errors.Error)
		if errors.As(err, e) {
			ctx.JSON(200, Response{Code: e.Code(), Info: e.Message()})
		} else {
			ctx.JSON(200, Response{Code: errors.InternalError.Code(), Info: errors.InternalError.Message()})
		}
		return
	}
	ctx.JSON(200, Response{Code: errors.Success.Code(), Info: errors.Success.Message(), Data: data})
}
