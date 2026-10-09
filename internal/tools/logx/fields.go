package logx

// Field keys used to tag log lines with request context.
const (
	CORRELATION_ID = "correlationId"
	CONTROLLER     = "controller"
	HTTP_METHOD    = "method"
	HTTP_PATH      = "path"
	HTTP_STATUS    = "status"
)

// Fields is an ordered set of key=value tags attached to a log line.
type Fields struct {
	keys []string
	data map[string]string
}

// NewFields returns an empty Fields set.
func NewFields() *Fields {
	return &Fields{data: map[string]string{}}
}

// Add sets a field and returns the receiver for chaining. Insertion order is
// preserved so the log output is deterministic.
func (f *Fields) Add(key, value string) *Fields {
	if _, ok := f.data[key]; !ok {
		f.keys = append(f.keys, key)
	}
	f.data[key] = value
	return f
}

// Each calls fn for each field in insertion order.
func (f *Fields) Each(fn func(key, value string)) {
	for _, k := range f.keys {
		fn(k, f.data[k])
	}
}
