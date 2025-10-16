package webapi

import (
	"net/http"

	"github.com/frimin/1pactus-react/app/onepacd/service/webapi/handler"
	"github.com/gin-gonic/gin"
)

func (s *WebApiService) setupRoute(r *gin.Engine) error {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"data": nil,
		})
	})

	groupApi := r.Group("/api")
	{
		handler.SetupNetworkStatus(groupApi)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"data": nil,
		})
	})

	return nil
}
