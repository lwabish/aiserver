package roop

import (
	"github.com/lwabish/cloudnative-ai-server/handlers"
	"github.com/lwabish/cloudnative-ai-server/utils"
	"net/http"
)

const (
	group    = "openvoice"
	TaskType = group
)

func init() {
	utils.RegisterRoutes(group, []utils.RouteInfo{
		{Path: "upload",
			Method:     http.MethodPost,
			Middleware: nil,
			Handler:    Handler.UploadFile},
	})
	utils.RegisterGroupMiddleware(group, handlers.MidHdl.Authenticate)
}
