package alpaca_eaf

import "encoding/json"

// stateEntry is a single {Name, Value} pair of a DeviceState response.
type stateEntry struct {
	Name  string          `json:"Name"`
	Value json.RawMessage `json:"Value"`
}
