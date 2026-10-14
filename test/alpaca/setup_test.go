//go:build integration

// Package alpaca holds the integration tests for the mount_panel Alpaca bridge.
//
// These tests drive the running bridge over HTTP (like the original Python
// suite in the firmware repo) and require a live mount + `go run .`. Point
// them at a different host with the ALPACA_BASE env var (default
// http://localhost:11111). They are gated behind the `integration` build tag
// so a plain `go test ./...` does not attempt them.
//
// They are a faithful port of ../bft_firmware/test/alpaca (pytest), minus the cases
// that depend on the firmware's debug REST endpoints (/api/status, /api/limits,
// /api/slew-to-coordinates), UDP discovery, or the NRM_TEST_MODE snapshot
// header. Those are t.Skip'd with a reason; the Alpaca-level assertions are
// otherwise unchanged.
package alpaca

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

// ── profile (mirrors the firmware's profile.json) ────────────────

const (
	defaultBase       = "http://localhost:11111"
	deviceNumber      = 0
	interfaceVersion  = 4
	alignmentMode     = 2
	equatorialSystem  = 1
	reportedPosDeg    = 0.1 // budget: reported position tolerance (deg)
	pollInterval      = 250 * time.Millisecond
	operationDeadline = 90 * time.Second
)

// requiredCapabilities are the ones a GEM mount must advertise.
var requiredCapabilities = map[string]string{
	"slew":        "canslew",
	"slew_async":  "canslewasync",
	"tracking":    "cansettracking",
	"pulse_guide": "canpulseguide",
	"park":        "canpark",
	"unpark":      "canunpark",
	"find_home":   "canfindhome",
}

func baseURL() string {
	if v := strings.TrimRight(os.Getenv("ALPACA_BASE"), "/"); v != "" {
		return v
	}
	return defaultBase
}

// ── client (mirrors the firmware's client.py) ────────────────────

// DeviceError is a non-zero ErrorNumber answered by the device.
type DeviceError struct {
	Number  int
	Message string
	Raw     string
}

func (e *DeviceError) Error() string {
	return fmt.Sprintf("Alpaca ErrorNumber=%d: %s", e.Number, e.Message)
}

// Response carries the standard Alpaca fields plus the raw body.
type Response struct {
	Raw                 string
	Value               any
	ErrorNumber         int
	ErrorMessage        string
	ClientTransactionID uint32
	ServerTransactionID uint32
	HTTPStatus          int
}

// Client is a minimal ASCOM Alpaca Telescope client.
type Client struct {
	base     string
	clientID int
	txn      uint32
	http     *http.Client
}

func newClient() *Client {
	return &Client{
		base:     baseURL(),
		clientID: 1,
		http:     &http.Client{Timeout: 8 * time.Second},
	}
}

// newTestClient returns a connected client and registers cleanup so motion and
// tracking do not leak between tests (mirrors the ensure_connected + _quiet
// fixtures). It also ensures a clean starting state, because the previous
// test's teardown is asynchronous (the firmware queues tracking commands).
func newTestClient(t *testing.T) *Client {
	t.Helper()
	c := newClient()
	c.putOK(t, "connected", url.Values{"Connected": {"true"}}, nil)

	c.putOK(t, "tracking", url.Values{"Tracking": {"false"}}, nil)
	waitTracking(t, c, false)

	t.Cleanup(func() {
		c.put("abortslew", nil, nil)
		c.put("tracking", url.Values{"Tracking": {"false"}}, nil)
		c.put("unpark", nil, nil)
	})
	return c
}

func (c *Client) nextTxn() uint32 {
	c.txn++
	return c.txn
}

// std fills ClientID/ClientTransactionID only if not already provided.
func (c *Client) std(params url.Values) url.Values {
	p := url.Values{}
	for k, v := range params {
		p[k] = v
	}
	if p.Get("ClientID") == "" {
		p.Set("ClientID", strconv.Itoa(c.clientID))
	}
	if p.Get("ClientTransactionID") == "" {
		p.Set("ClientTransactionID", strconv.FormatUint(uint64(c.nextTxn()), 10))
	}
	return p
}

func (c *Client) telescopeURL(name string) string {
	return fmt.Sprintf("%s/api/v1/telescope/%d/%s", c.base, deviceNumber, name)
}

// get issues a GET and returns the parsed response.
func (c *Client) get(name string, params url.Values) (*Response, error) {
	u := c.telescopeURL(name)
	if q := c.std(params).Encode(); q != "" {
		u += "?" + q
	}
	return c.do("GET", u, nil, true)
}

// put issues a PUT with a URL-encoded form body.
func (c *Client) put(name string, form, params url.Values) (*Response, error) {
	u := c.telescopeURL(name)
	if q := c.std(params).Encode(); q != "" {
		u += "?" + q
	}
	var body io.Reader
	if len(form) > 0 {
		body = strings.NewReader(form.Encode())
	}
	return c.do("PUT", u, body, true)
}

// management issues a management GET (no ErrorNumber contract).
func (c *Client) management(name string) (*Response, error) {
	u := fmt.Sprintf("%s/management/%s", c.base, name)
	return c.do("GET", u, nil, false)
}

func (c *Client) do(method, u string, body io.Reader, checkAlpaca bool) (*Response, error) {
	req, err := http.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 500 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, truncate(raw, 200))
	}

	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, fmt.Errorf("non-JSON response: %s", truncate(raw, 200))
	}
	if checkAlpaca {
		if _, ok := obj["ErrorNumber"]; !ok {
			return nil, fmt.Errorf("missing ErrorNumber: %s", truncate(raw, 200))
		}
	}

	r := &Response{Raw: string(raw), HTTPStatus: resp.StatusCode}
	r.Value = obj["Value"]
	r.ErrorNumber = asInt(obj["ErrorNumber"])
	r.ErrorMessage, _ = obj["ErrorMessage"].(string)
	r.ClientTransactionID = asUint32(obj["ClientTransactionID"])
	r.ServerTransactionID = asUint32(obj["ServerTransactionID"])
	return r, nil
}

// getValue returns the Value, erroring on a non-zero ErrorNumber.
func (c *Client) getValue(name string, params url.Values) (any, error) {
	r, err := c.get(name, params)
	if err != nil {
		return nil, err
	}
	if r.ErrorNumber != 0 {
		return nil, &DeviceError{r.ErrorNumber, r.ErrorMessage, r.Raw}
	}
	return r.Value, nil
}

// putOK is a PUT that expects success (ErrorNumber == 0).
func (c *Client) putOK(t *testing.T, name string, form, params url.Values) *Response {
	t.Helper()
	r, err := c.put(name, form, params)
	if err != nil {
		t.Fatalf("PUT %s: %v", name, err)
	}
	if r.ErrorNumber != 0 {
		t.Fatalf("PUT %s rejected: ErrorNumber=%d %s", name, r.ErrorNumber, r.ErrorMessage)
	}
	return r
}

// ── typed value getters ──────────────────────────────────────────

func (c *Client) boolValue(name string, params url.Values) (bool, error) {
	v, err := c.getValue(name, params)
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("%s is not bool: %v", name, v)
	}
	return b, nil
}

func (c *Client) floatValue(name string, params url.Values) (float64, error) {
	v, err := c.getValue(name, params)
	if err != nil {
		return 0, err
	}
	switch n := v.(type) {
	case float64:
		return n, nil
	case int:
		return float64(n), nil
	default:
		return 0, fmt.Errorf("%s is not numeric: %v", name, v)
	}
}

func (c *Client) intValue(name string, params url.Values) (int, error) {
	v, err := c.getValue(name, params)
	if err != nil {
		return 0, err
	}
	switch n := v.(type) {
	case float64:
		return int(n), nil
	case int:
		return n, nil
	default:
		return 0, fmt.Errorf("%s is not int: %v", name, v)
	}
}

func (c *Client) stringValue(name string, params url.Values) (string, error) {
	v, err := c.getValue(name, params)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s is not string: %v", name, v)
	}
	return s, nil
}

// must* wrappers fatal on error, matching the pytest client.get_value() usage.
func mustBool(t *testing.T, c *Client, name string) bool {
	t.Helper()
	v, err := c.boolValue(name, nil)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return v
}

func mustFloat(t *testing.T, c *Client, name string) float64 {
	t.Helper()
	v, err := c.floatValue(name, nil)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return v
}

func mustInt(t *testing.T, c *Client, name string) int {
	t.Helper()
	v, err := c.intValue(name, nil)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return v
}

func mustString(t *testing.T, c *Client, name string) string {
	t.Helper()
	v, err := c.stringValue(name, nil)
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	return v
}

// ── bounded polling helpers (mirrors conftest.wait_until) ────────

func waitUntil(t *testing.T, desc string, timeout time.Duration, getter func() (bool, error)) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		done, err := getter()
		if err == nil && done {
			return
		}
		lastErr = err
		time.Sleep(pollInterval)
	}
	t.Fatalf("timeout waiting for %s (last error: %v)", desc, lastErr)
}

// radec is a snapshot of the reported celestial position.
type radec struct {
	ra, dec    float64
	slewing    bool
	tracking   bool
	sideOfPier int
	atPark     bool
	atHome     bool
}

func readRADec(t *testing.T, c *Client) radec {
	t.Helper()
	return radec{
		ra:      mustFloat(t, c, "rightascension"),
		dec:     mustFloat(t, c, "declination"),
		slewing: mustBool(t, c, "slewing"),
	}
}

// currentRadec returns a full snapshot (mirrors conftest.current_radec).
func currentRadec(t *testing.T, c *Client) radec {
	t.Helper()
	return radec{
		ra:         mustFloat(t, c, "rightascension"),
		dec:        mustFloat(t, c, "declination"),
		slewing:    mustBool(t, c, "slewing"),
		tracking:   mustBool(t, c, "tracking"),
		sideOfPier: mustInt(t, c, "sideofpier"),
		atPark:     mustBool(t, c, "atpark"),
		atHome:     mustBool(t, c, "athome"),
	}
}

func waitNotSlewing(t *testing.T, c *Client, timeout ...time.Duration) radec {
	t.Helper()
	to := operationDeadline
	if len(timeout) > 0 {
		to = timeout[0]
	}
	var last radec
	waitUntil(t, "slewing=false", to, func() (bool, error) {
		ra, err := c.floatValue("rightascension", nil)
		if err != nil {
			return false, err
		}
		dec, err := c.floatValue("declination", nil)
		if err != nil {
			return false, err
		}
		slewing, err := c.boolValue("slewing", nil)
		if err != nil {
			return false, err
		}
		last = radec{ra: ra, dec: dec, slewing: slewing}
		return !slewing, nil
	})
	return last
}

func waitSlewing(t *testing.T, c *Client, timeout ...time.Duration) {
	t.Helper()
	to := 10 * time.Second
	if len(timeout) > 0 {
		to = timeout[0]
	}
	waitUntil(t, "slewing=true", to, func() (bool, error) {
		return c.boolValue("slewing", nil)
	})
}

func waitTracking(t *testing.T, c *Client, value bool, timeout ...time.Duration) {
	t.Helper()
	to := 10 * time.Second
	if len(timeout) > 0 {
		to = timeout[0]
	}
	waitUntil(t, fmt.Sprintf("tracking=%v", value), to, func() (bool, error) {
		v, err := c.boolValue("tracking", nil)
		if err != nil {
			return false, err
		}
		return v == value, nil
	})
}

// ── site/time setup (GOTO requirement) ───────────────────────────

func setupSiteTime(t *testing.T, c *Client) {
	t.Helper()
	c.putOK(t, "utcdate", url.Values{"UTCDate": {"2026-09-28T02:00:00"}}, nil)
	c.putOK(t, "sitelatitude", url.Values{"SiteLatitude": {"-32.89"}}, nil)
	c.putOK(t, "sitelongitude", url.Values{"SiteLongitude": {"-68.83"}}, nil)
}

// ── celestial separation assertion (mirrors oracle.py) ───────────

func angularSeparationDeg(ra1h, dec1, ra2h, dec2 float64) float64 {
	ra1 := math.Mod(ra1h, 24.0) * 15.0 * math.Pi / 180.0
	ra2 := math.Mod(ra2h, 24.0) * 15.0 * math.Pi / 180.0
	d1 := dec1 * math.Pi / 180.0
	d2 := dec2 * math.Pi / 180.0

	x1 := math.Cos(d1) * math.Cos(ra1)
	y1 := math.Cos(d1) * math.Sin(ra1)
	z1 := math.Sin(d1)
	x2 := math.Cos(d2) * math.Cos(ra2)
	y2 := math.Cos(d2) * math.Sin(ra2)
	z2 := math.Sin(d2)

	dot := x1*x2 + y1*y2 + z1*z2
	dot = math.Max(-1, math.Min(1, dot))
	return math.Acos(dot) * 180.0 / math.Pi
}

func assertAngularClose(t *testing.T, ra1, dec1, ra2, dec2, tol float64, ctx string) {
	t.Helper()
	sep := angularSeparationDeg(ra1, dec1, ra2, dec2)
	if sep > tol {
		t.Errorf("%s separation %.4f° > tolerance %.4f°", ctx, sep, tol)
	}
}

// gotoArrival issues a goto and verifies arrival on the reported coordinates.
// (Replaces assert_goto_duration: the raw axis distance the timing model needs
// is not observable over Alpaca, so only arrival + slewing state are checked.)
func gotoArrival(t *testing.T, c *Client, ra, dec float64) radec {
	t.Helper()
	c.putOK(t, "slewtocoordinatesasync",
		url.Values{"RightAscension": {fs(ra)}, "Declination": {fs(dec)}}, nil)
	final := waitNotSlewing(t, c)
	assertAngularClose(t, final.ra, final.dec, ra, dec, reportedPosDeg, "goto arrival")
	return final
}

// ── small helpers ─────────────────────────────────────────────────

func fs(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

func uitoa(v uint32) string { return strconv.FormatUint(uint64(v), 10) }

// mathMod is Python-style modulo: the result is always non-negative.
func mathMod(x, m float64) float64 {
	r := math.Mod(x, m)
	if r < 0 {
		r += m
	}
	return r
}

// discover performs ASCOM Alpaca Discovery (UDP 32227) against a host and
// returns the parsed JSON response.
func discover(t *testing.T, host string) map[string]any {
	t.Helper()
	conn, err := net.Dial("udp", net.JoinHostPort(host, "32227"))
	if err != nil {
		t.Fatalf("discover dial: %v", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("alpacadiscovery1")); err != nil {
		t.Fatalf("discover write: %v", err)
	}
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatalf("discover read: %v", err)
	}
	var resp map[string]any
	if err := json.Unmarshal(buf[:n], &resp); err != nil {
		t.Fatalf("discover parse: %v", err)
	}
	return resp
}

func truncate(b []byte, n int) string {
	s := string(b)
	if len(s) > n {
		return s[:n] + "…"
	}
	return s
}

func asInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	}
	return 0
}

func asUint32(v any) uint32 {
	switch n := v.(type) {
	case float64:
		return uint32(n)
	case int:
		return uint32(n)
	}
	return 0
}
