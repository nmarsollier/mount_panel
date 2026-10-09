package mount

import "github.com/nmarsollier/nrm_ha_go/internal/proto"

// Limits parameters accepted by the ESP's action:limits command.
const (
	LimitsSetHome     = "set_home"
	LimitsSetRALeft   = "set_ra_left"
	LimitsSetRARight  = "set_ra_right"
	LimitsSetDecLeft  = "set_dec_left"
	LimitsSetDecRight = "set_dec_right"
	LimitsClear       = "clear_limits"
)

// SetLimits fixes home or an axis limit from the current physical position, or
// clears the limits (LimitsClear).
func (m *Mount) SetLimits(param string) error {
	return m.action(proto.ActionRequest{Action: "limits", Param: param})
}
