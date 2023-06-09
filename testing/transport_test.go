package testing

import (
	"fmt"
	"github.com/kpassapk/golang-labs/testing/test"
	"github.com/stretchr/testify/assert"
	"github.com/yalochat/go-components/tester"
	"net/http"
	"testing"
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

func aRoundTripperReturning(response *http.Response, err error) *test.FakeRoundTripper {
	rt := &test.FakeRoundTripper{}
	rt.RoundTripReturns(response, err)
	return rt
}

// Test_timedTransport_RoundTrip was created using the Goland IDE
// and modified a little bit:
// - got rid of 'want' and 'wantErr'
// - added 'ret' and 'feature' type aliases
// - started the feature
//
// TODO: add success and error test cases
func Test_timedTransport_RoundTrip(t *testing.T) {
	type fields struct {
		rtp   http.RoundTripper
		clock Clock
	}
	type args struct {
		r *http.Request
	}
	type ret = *http.Response
	type feature = tester.Feature[fields, args, ret]

	// TODO uncomment and implement
	// roundTrip := func(f fields, a args) (ret, error) {
	// }

	tests := []struct {
		name     string
		args     args
		scenario func(*feature)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			s := tester.NewFeature[fields, args, ret](a, tt.args)
			tt.scenario(s)
		})
	}
}

// Test_timedTransport_Duration is a Table Test, and was created using the Goland IDE unmodified
func Test_timedTransport_Duration(t *testing.T) {
	type fields struct {
		rtp      http.RoundTripper
		clock    Clock
		reqStart time.Time
		reqEnd   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Duration
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &timedTransport{
				rtp:      tt.fields.rtp,
				clock:    tt.fields.clock,
				reqStart: tt.fields.reqStart,
				reqEnd:   tt.fields.reqEnd,
			}
			if got := tr.Duration(); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
