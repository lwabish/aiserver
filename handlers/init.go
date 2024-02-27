package handlers

import (
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
)

func init() {
	utils.RegisterRoutes(utils.RootGroup, []utils.RouteInfo{
		{Path: "status",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    BaseHdl.GetTaskStatus},
		{Path: "download",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    BaseHdl.DownloadResult},
	})
	utils.RegisterGroupMiddleware(utils.RootGroup, MidHdl.Authenticate)

}
