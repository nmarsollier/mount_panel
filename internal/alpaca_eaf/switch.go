package alpaca_eaf

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/mount_panel/internal/alpaca_mount/server"
	"github.com/nmarsollier/mount_panel/internal/fan"
	"github.com/nmarsollier/mount_panel/internal/temp"
)

const switchDevice = "/api/v1/switch/0"

// Switch channels exposed by the accessory controller:
//
//	0 = Fan Power            (boolean, writable)
//	1 = Fan Speed            (numeric 0..100 %, writable)
//	2 = Fan RPM              (numeric, read-only)
//	3 = Mirror Temperature   (numeric, read-only)
//	4 = Ambient Temperature  (numeric, read-only)
const (
	switchFanPower    = 0
	switchFanSpeed    = 1
	switchFanRPM      = 2
	switchMirrorTemp  = 3
	switchAmbientTemp = 4
	maxSwitch         = 5
)

var errInvalidSwitchID = errors.New("invalid switch id")

// initSwitch registers the Switch V3 endpoints for device 0.
func initSwitch(engine *gin.Engine) {
	g := engine.Group(switchDevice)

	// Common device endpoints.
	g.GET("/connected", getConnectedSwitch)
	g.PUT("/connected", putConnectedSwitch)
	g.PUT("/connect", putConnectSwitch)
	g.PUT("/disconnect", putDisconnectSwitch)
	g.GET("/connecting", getConnecting)
	g.GET("/description", getDescription)
	g.GET("/driverinfo", getDriverInfo)
	g.GET("/driverversion", getDriverVersion)
	g.GET("/interfaceversion", getInterfaceVersionSwitch)
	g.GET("/name", getName)
	g.GET("/supportedactions", getSupportedActions)
	g.GET("/devicestate", getDeviceStateSwitch)
	g.PUT("/commandblind", putCommandBlind)
	g.PUT("/commandbool", putCommandBool)
	g.PUT("/commandstring", putCommandString)
	g.PUT("/action", putAction)

	// Switch endpoints.
	g.GET("/maxswitch", getMaxSwitch)
	g.GET("/canwrite", getCanWrite)
	g.GET("/getswitch", getSwitch)
	g.GET("/getswitchdescription", getSwitchDescription)
	g.GET("/getswitchname", getSwitchName)
	g.PUT("/setswitchname", setSwitchName)
	g.GET("/getswitchvalue", getSwitchValue)
	g.GET("/minswitchvalue", getMinSwitchValue)
	g.GET("/maxswitchvalue", getMaxSwitchValue)
	g.GET("/switchstep", getSwitchStep)
	g.PUT("/setswitch", setSwitch)
	g.PUT("/setswitchvalue", setSwitchValue)
	g.GET("/canasync", getCanAsync)
	g.PUT("/setasync", setAsync)
	g.PUT("/setasyncvalue", setAsyncValue)
	g.GET("/statechangecomplete", getStateChangeComplete)
	g.PUT("/cancelasync", cancelAsync)
}

// getConnectedSwitch reports whether the Switch client is logically connected.
func getConnectedSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, fan.Get().Connected(), cid, stx)
}

// putConnectedSwitch tracks the Switch client's logical link.
func putConnectedSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	connected, ok := server.ParamBool(c, "Connected")
	if !ok {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Connected", cid, stx)
		return
	}
	fan.Get().SetConnected(connected)
	server.WriteOK(c, cid, stx)
}

// putConnectSwitch establishes the Switch logical client link.
func putConnectSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	fan.Get().SetConnected(true)
	server.WriteOK(c, cid, stx)
}

// putDisconnectSwitch drops the Switch logical client link.
func putDisconnectSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	fan.Get().SetConnected(false)
	server.WriteOK(c, cid, stx)
}

// getInterfaceVersionSwitch returns the ASCOM interface version (3).
func getInterfaceVersionSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 3, cid, stx)
}

// getMaxSwitch returns the number of switch channels (5).
func getMaxSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, maxSwitch, cid, stx)
}

// getCanWrite reports whether a channel is writable (only power and speed).
func getCanWrite(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteBool(c, id == switchFanPower || id == switchFanSpeed, cid, stx)
}

// getSwitch reports the boolean state of a channel (true = non-minimum value).
func getSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	v, err := switchValue(id)
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteBool(c, v != 0, cid, stx)
}

// getSwitchValue returns the numeric value of a channel.
func getSwitchValue(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	v, err := switchValue(id)
	if err != nil {
		server.WriteError(c, server.ErrNotConnected, err.Error(), cid, stx)
		return
	}
	server.WriteFloat(c, v, cid, stx)
}

// getSwitchName returns the channel name.
func getSwitchName(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteString(c, switchName(id), cid, stx)
}

// getSwitchDescription returns the channel description.
func getSwitchDescription(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteString(c, switchDescription(id), cid, stx)
}

// setSwitchName persists a channel alias in Go (no CDC command).
func setSwitchName(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	// Names are stored in Go and, for now, are not persisted — accepting the
	// write keeps the V3 member implemented.
	if _, ok := server.ParamRaw(c, "Name"); !ok {
		server.WriteError(c, server.ErrInvalidValue, "Missing Name", cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// getMinSwitchValue returns the minimum value of a channel.
func getMinSwitchValue(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteFloat(c, switchMin(id), cid, stx)
}

// getMaxSwitchValue returns the maximum value of a channel.
func getMaxSwitchValue(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteFloat(c, switchMax(id), cid, stx)
}

// getSwitchStep returns the value step of a channel.
func getSwitchStep(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteFloat(c, switchStep(id), cid, stx)
}

// setSwitch sets a channel's boolean state. For the fan power it selects the
// minimum (off) or maximum (manual); for speed it selects 1 or 100.
func setSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	state, okState := server.ParamBool(c, "State")
	if !ok || !okState || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id or State", cid, stx)
		return
	}
	if id != switchFanPower && id != switchFanSpeed {
		server.WriteError(c, server.ErrNotImplemented, "Channel is read-only", cid, stx)
		return
	}

	var err error
	if id == switchFanPower {
		mode := "off"
		if state {
			mode = "manual"
		}
		err = fan.Get().SetFanMode(mode)
	} else {
		speed := 1
		if state {
			speed = 100
		}
		err = fan.Get().SetFanSpeed(speed)
	}
	if err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// setSwitchValue sets a channel's numeric value.
func setSwitchValue(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	value, okValue := server.ParamFloat(c, "Value")
	if !ok || !okValue || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id or Value", cid, stx)
		return
	}
	if id != switchFanPower && id != switchFanSpeed {
		server.WriteError(c, server.ErrNotImplemented, "Channel is read-only", cid, stx)
		return
	}

	var err error
	if id == switchFanPower {
		mode := "off"
		if value != 0 {
			mode = "manual"
		}
		err = fan.Get().SetFanMode(mode)
	} else {
		speed := int(value)
		if speed < 1 {
			speed = 1
		} else if speed > 100 {
			speed = 100
		}
		err = fan.Get().SetFanSpeed(speed)
	}
	if err != nil {
		server.WriteError(c, server.ErrInvalidValue, err.Error(), cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// getCanAsync reports whether a channel supports async operations (false).
func getCanAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteBool(c, false, cid, stx)
}

// setAsync is not implemented (CanAsync=false).
func setAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Async not supported", cid, stx)
}

// setAsyncValue is not implemented (CanAsync=false).
func setAsyncValue(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Async not supported", cid, stx)
}

// getStateChangeComplete is not implemented (CanAsync=false).
func getStateChangeComplete(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteError(c, server.ErrNotImplemented, "Async not supported", cid, stx)
}

// cancelAsync validates the Id and does nothing (no async operation exists).
func cancelAsync(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteOK(c, cid, stx)
}

// getDeviceStateSwitch returns the Switch operational properties.
func getDeviceStateSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	entries := []stateEntry{}

	for id := 0; id < maxSwitch; id++ {
		v, err := switchValue(id)
		if err != nil {
			continue
		}
		entries = append(entries,
			stateEntry{switchStateName(id), server.RawBool(v != 0)},
			stateEntry{switchValueName(id), server.RawFloat(v)},
		)
	}

	value, _ := json.Marshal(entries)
	server.WriteValue(c, value, cid, stx)
}

// switchValue returns the numeric value of a switch channel.
func switchValue(id int) (float64, error) {
	switch id {
	case switchFanPower:
		on, err := fan.Get().FanPower()
		if err != nil {
			return 0, err
		}
		if on {
			return 1, nil
		}
		return 0, nil
	case switchFanSpeed:
		s, err := fan.Get().FanSpeed()
		return float64(s), err
	case switchFanRPM:
		r, err := fan.Get().FanRPM()
		return float64(r), err
	case switchMirrorTemp:
		return fan.Get().MirrorTemperature()
	case switchAmbientTemp:
		return temp.Get().AmbientTemperature()
	default:
		return 0, errInvalidSwitchID
	}
}

// switchName returns the channel name.
func switchName(id int) string {
	switch id {
	case switchFanPower:
		return "Fan Power"
	case switchFanSpeed:
		return "Fan Speed"
	case switchFanRPM:
		return "Fan RPM"
	case switchMirrorTemp:
		return "Mirror Temperature"
	case switchAmbientTemp:
		return "Ambient Temperature"
	default:
		return ""
	}
}

// switchDescription returns the channel description.
func switchDescription(id int) string {
	switch id {
	case switchFanPower:
		return "Enciende/apaga el relé de 12 V del ventilador"
	case switchFanSpeed:
		return "Velocidad PWM del ventilador, en % (0–100)"
	case switchFanRPM:
		return "RPM real del ventilador (tacómetro)"
	case switchMirrorTemp:
		return "Temperatura del espejo primario (DS18B20)"
	case switchAmbientTemp:
		return "Temperatura ambiente (SHT31)"
	default:
		return ""
	}
}

// switchMin / switchMax / switchStep return the value bounds of a channel.
func switchMin(id int) float64 {
	if id == switchMirrorTemp || id == switchAmbientTemp {
		return -50
	}
	return 0
}

func switchMax(id int) float64 {
	switch id {
	case switchFanPower:
		return 1
	case switchFanSpeed:
		return 100
	case switchFanRPM:
		return 100000
	default:
		return 100
	}
}

func switchStep(id int) float64 {
	if id == switchMirrorTemp || id == switchAmbientTemp {
		return 0.1
	}
	return 1
}

// switchStateName / switchValueName name the DeviceState entries for a channel.
func switchStateName(id int) string {
	return "GetSwitch" + strconv.Itoa(id)
}

func switchValueName(id int) string {
	return "GetSwitchValue" + strconv.Itoa(id)
}
