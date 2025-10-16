package handler

import (
	"net/http"

	"github.com/frimin/1pactus-react/app/onepacd/service/webapi/model"
	"github.com/frimin/1pactus-react/app/onepacd/store"
	"github.com/frimin/1pactus-react/proto/gen/go/api"
	"github.com/gin-gonic/gin"
)

func SetupNetworkStatus(group *gin.RouterGroup) {
	networkStatusRoute := group.Group("/network_status")

	networkStatusRoute.GET("/get", func(c *gin.Context) {
		/*var req api.GetNetworkHealthRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}*/

		httpResp := model.ApiResponse[api.GetNetworkHealthResponse]{
			Data: &api.GetNetworkHealthResponse{},
		}

		defer func() {
			if httpResp.Msg == "" {
				httpResp.Msg = model.ErrorFromCode(httpResp.Code).Error()
			}
			c.JSON(http.StatusOK, httpResp)
		}()

		stats, err := store.Mongo.FetchNetworkGlobalStats(-1)

		if err != nil {
			httpResp.Code = model.Code_InternalError
			httpResp.Data = nil
		}

		httpResp.Data.Lines = make([]*api.NetworkStatusData, 0, len(stats))

		for _, s := range stats {
			httpResp.Data.Lines = append(httpResp.Data.Lines, s.ToProto())
		}

		if err != nil {
			httpResp.Code = model.Code_InternalError
			httpResp.Data = nil
		} else {
			httpResp.Code = model.Code_Success
		}

	})
}
