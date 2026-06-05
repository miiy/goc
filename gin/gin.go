package gin

import "github.com/gin-gonic/gin"

type (
	Engine      = gin.Engine
	Context     = gin.Context
	H           = gin.H
	HandlerFunc = gin.HandlerFunc
	IRouter     = gin.IRouter
	IRoutes     = gin.IRoutes
)

func New() *Engine {
	return gin.New()
}

func Default() *Engine {
	return gin.Default()
}
