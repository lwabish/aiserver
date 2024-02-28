package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lwabish/cloudnative-ai-server/utils"
)

// RegisterRoutes register gin routes
func RegisterRoutes(router *gin.Engine) {
	for group, routes := range utils.Routes {
		if group == utils.RootGroup {
			for _, route := range routes {
				router.Handle(route.Method, route.Path, append(route.Middleware, route.Handler)...)
			}
			continue
		}
		g := router.Group(string(group))
		if m, ok := utils.GroupMiddlewares[group]; ok {
			g.Use(m...)
		}
		for _, route := range routes {
			g.Handle(route.Method, route.Path, append(route.Middleware, route.Handler)...)
		}
	}
}
