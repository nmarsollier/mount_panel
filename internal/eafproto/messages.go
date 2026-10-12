package eafproto

import "encoding/json"

// Error is a rejected request: the accessory answered ok=false with an error
// object {code, message, field}.
type Error struct {
	Code    string  `json:"code"`
	Message string  `json:"message"`
	Field   *string `json:"field"`
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Code != "" {
		return e.Code
	}
	return "accessory protocol error"
}

// envelope is the common fields of every accessory response. Error is kept raw
// because the transport injects a string form on a lost link while the
// accessory replies with the object form.
type envelope struct {
	Type  string          `json:"type"`
	ID    uint32          `json:"id"`
	OK    bool            `json:"ok"`
	Error json.RawMessage `json:"error"`
}

// parseError turns a raw `error` field into an error, accepting both the
// object form ({code,message,field}) and the transport's string form.
func parseError(raw json.RawMessage) error {
	if len(raw) == 0 || string(raw) == "null" {
		return &Error{}
	}
	var e Error
	if err := json.Unmarshal(raw, &e); err == nil {
		return &e
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return &Error{Message: s}
	}
	return &Error{Message: string(raw)}
}
