package testing

import (
	"fmt"
	"github.com/kpassapk/golang-labs/testing/test"
	"net/http"
	"time"
)

var wayback = time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

type exampleClock struct {
	call  int
	delta time.Duration
}

// Now returns a time that is wayback + call * delta for testing purposes
func (e *exampleClock) Now() time.Time {
	e.call = e.call + 1
	return wayback.Add(time.Duration(e.call) * e.delta)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o assert net/http.RoundTripper

func ExampleCustomTransport_Duration() {

	tr := &test.FakeRoundTripper{}
	// To make a real request
	// tr := http.DefaultTransport

	clock := &exampleClock{delta: 1 * time.Second}
	// To use the real system clock
	// clock := realClock{}

	ttr := newTimedTransport(tr, clock)

	c := http.Client{
		Transport: ttr,
	}

	c.Get("http://google.com")
	fmt.Println(ttr.Duration())
	// Output: 1s
}
