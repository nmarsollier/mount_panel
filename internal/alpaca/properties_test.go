package alpaca

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/mount"
)

// propertyHandlers are the read-only handlers that fetch from the device. When
// the serial link is down, each must answer Not Connected (1027) and never a
// silent success with a false/zero value (fix #4).
var propertyHandlers = []struct {
	name string
	h    gin.HandlerFunc
}{
	{"slewing", getSlewing},
	{"tracking", getTracking},
	{"atpark", getAtPark},
	{"athome", getAtHome},
	{"ispulseguiding", getIsPulseGuiding},
	{"rightascension", getRightAscension},
	{"declination", getDeclination},
	{"sideofpier", getSideOfPier},
	{"siderealtime", getSiderealTime},
	{"altitude", getAltitude},
	{"azimuth", getAzimuth},
	{"sitelatitude", getSiteLatitude},
	{"sitelongitude", getSiteLongitude},
	{"siteelevation", getSiteElevation},
	{"guideraterightascension", getGuideRateRightAscension},
	{"guideratedeclination", getGuideRateDeclination},
	{"utcdate", getUTCDate},
	{"axisrates", getAxisRates},
	{"devicestate", getDeviceState},
}

// TestPropertyHandlersReportNotConnected drives every device read with the
// mount's serial link down and asserts the handler surfaces ErrorNumber=1027
// instead of a zero value. It is fully automatic: the singleton mount is wired
// to a serial port that can never open, so every read fails immediately.
func TestPropertyHandlersReportNotConnected(t *testing.T) {
	// No device can exist at this path, so conn stays nil and every read fails.
	mount.Init("/dev/nonexistent-serial-port")

	for _, ph := range propertyHandlers {
		t.Run(ph.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			ph.h(c)

			var r response
			if err := json.Unmarshal(w.Body.Bytes(), &r); err != nil {
				t.Fatalf("handler wrote non-JSON: %v", err)
			}
			if r.ErrorNumber != errNotConnected {
				t.Errorf("ErrorNumber=%d, want %d (Not Connected)", r.ErrorNumber, errNotConnected)
			}
		})
	}
}
