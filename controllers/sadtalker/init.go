package sadtalker

import "github.com/lwabish/cloudnative-ai-server/utils"

const (
	group = "sadTalker"
)

func init() {
	utils.RegisterRoutes(group, []utils.RouteInfo{
		{Path: "upload",
			Method:     "POST",
			Middleware: nil,
			Handler:    stCtl.UploadFile},
	})
}
