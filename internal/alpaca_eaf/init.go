package alpaca_eaf

import "github.com/gin-gonic/gin"

// Init registers the accessory controller's Alpaca devices: Focuser V4,
// ObservingConditions V2 and Switch V3.
func Init(engine *gin.Engine) {
	initFocuser(engine)
	initObservingConditions(engine)
	initSwitch(engine)
}
