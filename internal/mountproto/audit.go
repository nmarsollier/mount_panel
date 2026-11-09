package mountproto

import (
	"encoding/json"

	"github.com/nmarsollier/mount_panel/internal/events"
)

// commandDetails builds the typed mount.command detail from the request family
// and its (possibly nil) typed payload. A family with a bare request
// (capabilities, state) yields only the envelope columns.
func commandDetails(family string, id uint32, payload any) events.MountCommand {
	cmd := events.MountCommand{Family: family, MsgID: id}
	switch p := payload.(type) {
	case ConfigRequest:
		cmd.Lat = p.Lat
		cmd.Lon = p.Lon
		cmd.Elevation = p.Elevation
		cmd.GuideRateRa = p.GuideRateRa
		cmd.GuideRateDec = p.GuideRateDec
		cmd.Utc = p.Utc
	case ControlRequest:
		cmd.Tracking = p.Tracking
		cmd.ManualRaDps = p.ManualRaDps
		cmd.ManualDecDps = p.ManualDecDps
	case ActionRequest:
		cmd.Action = p.Action
		cmd.Scope = p.Scope
		cmd.Ra = p.Ra
		cmd.Dec = p.Dec
		cmd.Speed = p.Speed
		cmd.Direction = p.Direction
		cmd.DurationMs = p.DurationMs
		cmd.Axis = p.Axis
		cmd.Degrees = p.Degrees
		cmd.Param = p.Param
	}
	return cmd
}

// responseDetails builds the typed mount.response detail from the response
// envelope and the (already unmarshalled) typed response, or only the envelope
// columns when the family's response carries no extra fields (control, action).
func responseDetails(family string, id uint32, env envelope, out any) events.MountResponse {
	resp := events.MountResponse{Family: family, MsgID: id, OK: env.OK, Error: env.Error}
	switch o := out.(type) {
	case *CapabilitiesResponse:
		resp.BootID = strPtr(o.BootID)
		resp.Proto = intPtr(o.Proto)
		resp.DeviceType = strPtr(o.DeviceType)
		resp.DeviceID = strPtr(o.DeviceID)
		resp.Name = strPtr(o.Name)
		resp.FirmwareVersion = strPtr(o.FirmwareVersion)
		caps := o.Capabilities.Caps
		resp.CapGuide = boolPtr(caps.Guide)
		resp.CapHome = boolPtr(caps.Home)
		resp.CapPark = boolPtr(caps.Park)
		resp.CapUnpark = boolPtr(caps.Unpark)
		resp.CapMoveAxis = boolPtr(caps.MoveAxis)
		resp.CapSlew = boolPtr(caps.Slew)
		resp.CapTracking = boolPtr(caps.Tracking)
		resp.CapSync = boolPtr(caps.Sync)
		resp.CapMeridianFlip = caps.MeridianFlip
		resp.CapPierSide = boolPtr(caps.PierSide)
		resp.Capabilities = marshalArrays(o.Capabilities)
	case *StateResponse:
		resp.State = strPtr(o.State)
		resp.Tracking = strPtr(o.Tracking)
		resp.RaSteps = int64Ptr(o.RaSteps)
		resp.DecSteps = int64Ptr(o.DecSteps)
		resp.RaAxisDeg = floatPtr(o.RaAxisDeg)
		resp.DecAxisDeg = floatPtr(o.DecAxisDeg)
		resp.Ra = floatPtr(o.Ra)
		resp.Dec = floatPtr(o.Dec)
		resp.Lst = floatPtr(o.Lst)
		resp.PierSide = strPtr(o.PierSide)
		resp.RaSpeedDps = floatPtr(o.RaSpeedDps)
		resp.DecSpeedDps = floatPtr(o.DecSpeedDps)
		resp.Guiding = boolPtr(o.Guiding)
		resp.AtHome = boolPtr(o.AtHome)
		resp.AtPark = boolPtr(o.AtPark)
		resp.Power = boolPtr(o.Power)
		resp.TimeValid = boolPtr(o.TimeValid)
		resp.RaMin = floatPtr(o.Limits.RaMin)
		resp.RaMax = floatPtr(o.Limits.RaMax)
		resp.DecMin = floatPtr(o.Limits.DecMin)
		resp.DecMax = floatPtr(o.Limits.DecMax)
		resp.ConfigRev = uint32Ptr(o.ConfigRev)
		resp.UptimeMs = uint64Ptr(o.UptimeMs)
	case *ConfigResponse:
		resp.ConfigRev = uint32Ptr(o.ConfigRev)
		resp.Utc = strPtr(o.Utc)
		resp.Lat = floatPtr(o.Config.Lat)
		resp.Lon = floatPtr(o.Config.Lon)
		resp.Elevation = intPtr(o.Config.Elevation)
		resp.GuideRateRa = floatPtr(o.Config.GuideRateRa)
		resp.GuideRateDec = floatPtr(o.Config.GuideRateDec)
	}
	return resp
}

// marshalArrays keeps only the array-valued capabilities (the caps booleans are
// already columns) as raw JSON for the capabilities_json column.
func marshalArrays(c MountCaps) json.RawMessage {
	b, _ := json.Marshal(struct {
		Axes              []string  `json:"axes"`
		SlewSpeedsDps     []float64 `json:"slew_speeds_dps"`
		MoveAxisSpeedsDps []float64 `json:"move_axis_speeds_dps"`
		TrackingModes     []string  `json:"tracking_modes"`
	}{
		Axes:              c.Axes,
		SlewSpeedsDps:     c.SlewSpeedsDps,
		MoveAxisSpeedsDps: c.MoveAxisSpeedsDps,
		TrackingModes:     c.TrackingModes,
	})
	return b
}

func strPtr(v string) *string     { return &v }
func intPtr(v int) *int           { return &v }
func boolPtr(v bool) *bool        { return &v }
func floatPtr(v float64) *float64 { return &v }
func int64Ptr(v int64) *int64     { return &v }
func uint32Ptr(v uint32) *uint32  { return &v }
func uint64Ptr(v uint64) *uint64  { return &v }
