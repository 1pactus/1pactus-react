package handler

import (
	"net/http"

	"github.com/frimin/1pactus-react/app/onepacd/service/webapi/model"
	"github.com/frimin/1pactus-react/app/onepacd/store"
	"github.com/frimin/1pactus-react/log"
	"github.com/frimin/1pactus-react/proto/gen/go/api"
	"github.com/gin-gonic/gin"
)

func SetupNetworkStatus(group *gin.RouterGroup) {
	//networkStatusRoute := group.Group("/network_status")

	group.GET("/network_status", func(c *gin.Context) {
		httpResp := &api.GetNetworkHealthResponse{}
		datatype := "json"

		defer func() {
			if httpResp.Msg == "" {
				httpResp.Msg = model.ErrorFromCode(httpResp.Code).Error()
			}

			switch datatype {
			case "json":
				c.JSON(http.StatusOK, httpResp)
			case "pb":
				c.ProtoBuf(http.StatusOK, httpResp)
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid datatype"})
			}
		}()

		var req api.GetNetworkHealthRequest

		if err := c.ShouldBindQuery(&req); err != nil {
			log.Error("failed to bind query params: ", err)
			return
		}

		if req.Datatype != "" {
			datatype = req.Datatype
		}

		if req.Days == 0 {
			req.Days = 30 // default to 30 days
		}

		stats, err := store.Mongo.FetchNetworkGlobalStats(int64(req.Days))

		if err != nil {
			httpResp.Code = model.Code_InternalError
		}

		httpResp.Lines = make([]*api.NetworkStatusData, 0, len(stats))

		for _, s := range stats {
			httpResp.Lines = append(httpResp.Lines, s.ToProto())
		}

		if err != nil {
			httpResp.Code = model.Code_InternalError
		} else {
			httpResp.Code = model.Code_Success
		}

	})
}
