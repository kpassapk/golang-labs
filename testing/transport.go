package testing

import (
	"net/http"
	"time"
)

// Problem 2: Test the Duration() method
// Problem 3: Test the RoundTrip() method
//
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
