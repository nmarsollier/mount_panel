package alpaca_eaf

import "github.com/gin-gonic/gin"

// Init registers the focuser controller's Alpaca device: Focuser V4.
func Init(engine *gin.Engine) {
	initFocuser(engine)
}
