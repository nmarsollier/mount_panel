package alpaca_fan

import "github.com/gin-gonic/gin"

// Init registers the fan controller's Alpaca device: Switch V3.
func Init(engine *gin.Engine) {
	initSwitch(engine)
}
