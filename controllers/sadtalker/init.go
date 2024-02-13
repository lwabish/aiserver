package sadtalker

import (
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
)

const (
	group    = "sadTalker"
	TaskType = group
)

func init() {
	utils.RegisterRoutes(group, []utils.RouteInfo{
		{Path: "upload",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    StCtl.UploadFile},
		{Path: "status",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    StCtl.GetTaskStatus},
		{Path: "download",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    StCtl.DownloadResult},
	})
}
