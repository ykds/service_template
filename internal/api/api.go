package api

import (
	"github.com/gin-gonic/gin"
	"service_template/internal/service"
)

type Api struct {
	ExampleApi *ExampleApi
}

func InitApi(srv *service.Service) *Api {
	return &Api{
		ExampleApi: NewExampleApi(srv),
	}
}

func NewExampleApi(srv service.ExampleService) *ExampleApi {
	return &ExampleApi{
		srv: srv,
	}
}

type ExampleApi struct {
	srv service.ExampleService
}

type ExampleReq struct{}

type ExampleResp struct{}

// Example
// @ID Example
// @Summary Example
// @Description
// @Tags ex
// @Param field1 query int true "comment"
// @Param field2 body ExampleReq false "comment"
// @Param field3 header string false "comment"
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=ExampleResp} "comment"
// @Failure 500 {string} string "comment"
// @Failure 400 "comment"
// @Router /example [get]
func (s *ExampleApi) Example(c *gin.Context) {
	c.JSON(200, "ok")
}
