package testing

import (
	"net"
	"net/http"
	"time"
)

// Problem 2: Testing an HTTP transport
// In Go, it is possible to change the behavior of the HTTP transport by implementing the http.RoundTripper interface.
// This is useful for testing, for example, to measure the time it takes to establish a connection.

type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

type timedTransport struct {
	rtp   http.RoundTripper
	clock Clock

	reqStart time.Time
	reqEnd   time.Time
}

func newTimedTransport(rtp http.RoundTripper, clock Clock) *timedTransport {
	return &timedTransport{
		clock: clock,
		rtp:   rtp,
	}
}

func (tr *timedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	tr.reqStart = tr.clock.Now()
	resp, err := tr.rtp.RoundTrip(r)
	tr.reqEnd = tr.clock.Now()
	return resp, err
}

func (tr *timedTransport) Duration() time.Duration {
	return tr.reqEnd.Sub(tr.reqStart)
}

func newTimedDialer(dialer *net.Dialer, clock Clock) *timedDialer {
	return &timedDialer{
		dialer: dialer,
		clock:  clock,
	}
}

type timedDialer struct {
	dialer *net.Dialer
	clock  Clock

	connStart time.Time
	connEnd   time.Time
}

func (tr *timedDialer) Duration() time.Duration {
	return tr.connEnd.Sub(tr.connStart)
}

func (tr *timedDialer) dial(network, addr string) (net.Conn, error) {
	tr.connStart = tr.clock.Now()
	cn, err := tr.dialer.Dial(network, addr)
	tr.connEnd = tr.clock.Now()
	return cn, err
}

type durationReporter interface {
	Duration() time.Duration
}

func newTimedClient(rtp durationReporter, conn durationReporter) *timedClient {
	return &timedClient{
		rtp:  rtp,
		conn: conn,
	}
}

type timedClient struct {
	rtp  durationReporter
	conn durationReporter
}

func (tc *timedClient) Duration() time.Duration {
	return tc.rtp.Duration() - tc.conn.Duration()
}
