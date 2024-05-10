package middleware

import (
	"github.com/gin-gonic/gin"
	"service_template/internal/common"
	"service_template/internal/errors"
	"service_template/internal/response"
)

func Authentication(auth func(token string) (int, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(common.TokenHeader)
		if token == "" {
			response.HandleResponse(c, errors.Unauthorized, nil, nil)
			c.Abort()
			return
		}
		userId, err := auth(token)
		if err != nil {
			response.HandleResponse(c, err, nil, nil)
			c.Abort()
			return
		}
		c.Set(common.UserIdKey, userId)
		c.Next()
	}
}
