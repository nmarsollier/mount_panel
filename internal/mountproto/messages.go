package mountproto

// Error is a rejected request: the mount answered ok=false with a message.
type Error struct {
	Message string
}

func (e *Error) Error() string {
	if e.Message == "" {
		return "protocol error"
	}
	return e.Message
}

// envelope is the common fields of every response, used to route the frame to
// the pending request and to detect rejections.
type envelope struct {
	Type  string `json:"type"`
	ID    uint32 `json:"id"`
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}
