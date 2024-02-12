package utils

import "github.com/gin-gonic/gin"

var (
	//Routes "group": path,method,handler ...
	Routes           = map[RouteGroup][]RouteInfo{}
	GroupMiddlewares = map[RouteGroup][]gin.HandlerFunc{}
)

const (
	RootGroup = "root"
)

type RouteGroup string

type RouteInfo struct {
	Path       string
	Method     string
	Middleware []gin.HandlerFunc
	Handler    gin.HandlerFunc
}

func RegisterGroupMiddleware(group RouteGroup, middlewares []gin.HandlerFunc) {
	GroupMiddlewares[group] = middlewares
}

func RegisterRoutes(group RouteGroup, routes []RouteInfo) {
	Routes[group] = append(Routes[group], routes...)
}
