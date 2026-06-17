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

// Mode helpers re-exported so callers using goc/gin don't need gin-gonic/gin.
var (
	SetMode     = gin.SetMode
	TestMode    = gin.TestMode
	DebugMode   = gin.DebugMode
	ReleaseMode = gin.ReleaseMode
)

func New() *Engine {
	return gin.New()
}

func Default() *Engine {
	return gin.Default()
}
