package alpaca_eaf

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/nmarsollier/nrm_ha_go/internal/alpaca_mount/server"
	"github.com/nmarsollier/nrm_ha_go/internal/eaf"
)

const switchDevice = "/api/v1/switch/0"

// Switch channels: 0 = fan power (the 12 V relay, boolean), 1 = fan speed
// (PWM 0..100 %). Both are writable.
const (
	switchFanPower = 0
	switchFanSpeed = 1
	maxSwitch      = 2
)

// initSwitch registers the Switch V3 endpoints for device 0.
func initSwitch(engine *gin.Engine) {
	g := engine.Group(switchDevice)

	// Common device endpoints.
	g.GET("/connected", getConnected)
	g.PUT("/connected", putConnected)
	g.PUT("/connect", putConnect)
	g.PUT("/disconnect", putDisconnect)
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

// getInterfaceVersionSwitch returns the ASCOM interface version (3).
func getInterfaceVersionSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 3, cid, stx)
}

// getMaxSwitch returns the number of switch channels (2).
func getMaxSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, maxSwitch, cid, stx)
}

// getCanWrite reports whether a channel is writable (both are).
func getCanWrite(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	server.WriteBool(c, true, cid, stx)
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
	server.WriteInt(c, v, cid, stx)
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

// getMinSwitchValue returns the minimum value (0 for both channels).
func getMinSwitchValue(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 0, cid, stx)
}

// getMaxSwitchValue returns the maximum value (1 for power, 100 for speed).
func getMaxSwitchValue(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	if !ok || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id", cid, stx)
		return
	}
	if id == switchFanPower {
		server.WriteInt(c, 1, cid, stx)
		return
	}
	server.WriteInt(c, 100, cid, stx)
}

// getSwitchStep returns the value step (1 for both channels).
func getSwitchStep(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	server.WriteInt(c, 1, cid, stx)
}

// setSwitch sets a channel's boolean state. For power it switches the relay;
// for speed it selects the minimum (0) or maximum (100).
func setSwitch(c *gin.Context) {
	cid, stx := server.ClientTransactionID(c), server.NextServerTx()
	id, ok := server.ParamInt(c, "Id")
	state, okState := server.ParamBool(c, "State")
	if !ok || !okState || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id or State", cid, stx)
		return
	}

	var err error
	if id == switchFanPower {
		err = eaf.Get().SetFanPower(state)
	} else {
		speed := 0
		if state {
			speed = 100
		}
		err = eaf.Get().SetFanSpeed(speed)
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
	value, okValue := server.ParamInt(c, "Value")
	if !ok || !okValue || id < 0 || id >= maxSwitch {
		server.WriteError(c, server.ErrInvalidValue, "Invalid Id or Value", cid, stx)
		return
	}

	var err error
	if id == switchFanPower {
		err = eaf.Get().SetFanPower(value != 0)
	} else {
		if value < 0 || value > 100 {
			server.WriteError(c, server.ErrInvalidValue, "Value out of range (0-100)", cid, stx)
			return
		}
		err = eaf.Get().SetFanSpeed(value)
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

	if v, err := switchValue(switchFanPower); err == nil {
		entries = append(entries, stateEntry{"GetSwitch0", server.RawBool(v != 0)})
		entries = append(entries, stateEntry{"GetSwitchValue0", server.RawInt(v)})
	}
	if v, err := switchValue(switchFanSpeed); err == nil {
		entries = append(entries, stateEntry{"GetSwitch1", server.RawBool(v != 0)})
		entries = append(entries, stateEntry{"GetSwitchValue1", server.RawInt(v)})
	}

	value, _ := json.Marshal(entries)
	server.WriteValue(c, value, cid, stx)
}

// switchValue returns the numeric value of a switch channel.
func switchValue(id int) (int, error) {
	if id == switchFanPower {
		on, err := eaf.Get().FanPower()
		if err != nil {
			return 0, err
		}
		if on {
			return 1, nil
		}
		return 0, nil
	}
	return eaf.Get().FanSpeed()
}

// switchName returns the channel name.
func switchName(id int) string {
	if id == switchFanPower {
		return "Alimentación del ventilador"
	}
	return "Velocidad del ventilador"
}

// switchDescription returns the channel description.
func switchDescription(id int) string {
	if id == switchFanPower {
		return "Enciende o apaga el relé de 12 V del ventilador"
	}
	return "Ciclo de trabajo PWM aplicado al ventilador, en % (0–100)"
}
