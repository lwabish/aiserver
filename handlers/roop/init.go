package roop

import (
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
)

const (
	group    = "roop"
	TaskType = group
)

func init() {
	utils.RegisterRoutes(group, []utils.RouteInfo{
		{Path: "upload",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    Handler.UploadFile},
		{Path: "status",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    Handler.GetTaskStatus},
		{Path: "download",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    Handler.DownloadResult},
	})
	utils.RegisterGroupMiddleware(group, handlers.MidHdl.Authenticate)
}
