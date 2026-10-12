package server

import "github.com/gin-gonic/gin"

var engine *gin.Engine

// Router returns the singleton GIN engine. GIN's request logger and debug
// output are disabled; only Recovery (which logs panics) is kept.
func Router() *gin.Engine {
	if engine != nil {
		return engine
	}

	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
	engine.Use(gin.Recovery())

	return engine
}
