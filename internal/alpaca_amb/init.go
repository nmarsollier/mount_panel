package alpaca_amb

import "github.com/gin-gonic/gin"

// Init registers the ambient controller's Alpaca device: ObservingConditions V2.
func Init(engine *gin.Engine) {
	initObservingConditions(engine)
}
