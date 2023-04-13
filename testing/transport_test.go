package testing

import (
	"fmt"
	"github.com/kpassapk/golang-labs/testing/test"
	"github.com/stretchr/testify/assert"
	"github.com/yalochat/go-components/tester"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"net"
	"net/http"
	"net/http/httptest"
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

type timedTransportFields = struct {
	rtp      *test.FakeRoundTripper
	dialer   *net.Dialer
	clock    Clock
	reqStart time.Time
	reqEnd   time.Time
}

func aSampleHTTPRequest() *http.Request {
	return httptest.NewRequest("", "/", nil)
}

func TestTimedTransport_RoundTrip(t *testing.T) {
	type fields = timedTransportFields
	type args struct {
		r *http.Request
	}
	// In this case, we're testing for a side effect, so the return type is the timed transport itself.
	type ret = *timedTransport
	type tst = tester.Feature[fields, args, ret]

	roundTrip := func(fields fields, args args) (ret, error) {
		rt := &timedTransport{
			clock: fields.clock,
			rtp:   fields.rtp,
		}
		_, err := rt.RoundTrip(args.r)
		return rt, err
	}

	assertDurationIs := func(ts *tst, d time.Duration) {
		got := ts.Response.reqEnd.Sub(ts.Response.reqStart)
		ts.Assert.Equal(d, got)
	}

	tests := []struct {
		name string
		args args
		test func(*tst)
	}{
		{
			name: "RoundTrip returns the response from the underlying RoundTripper",
			args: args{
				r: aSampleHTTPRequest(),
			},
			test: func(ts *tst) {
				f := fields{
					clock: &exampleClock{delta: 1 * time.Second},
				}
				ts.GivenFields(f)
				ts.When(roundTrip)
				ts.AssertNoError()
				assertDurationIs(ts, 1*time.Second)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			ts := tester.NewFeature[fields, args, ret](a, tt.args)
			tt.test(ts)
		})
	}
}

func Test_timedTransport_Duration(t *testing.T) {
	type args struct{}
	type fields struct {
		rtp      http.RoundTripper
		clock    Clock
		reqStart time.Time
		reqEnd   time.Time
	}

	getDuration := func(fields fields, args args) (time.Duration, error) {
		tr := &timedTransport{
			rtp:      fields.rtp,
			clock:    fields.clock,
			reqStart: fields.reqStart,
			reqEnd:   fields.reqEnd,
		}
		return tr.Duration(), nil
	}

	tests := []struct {
		name   string
		fields fields
		assert func(feature *tester.Feature[fields, args, time.Duration])
	}{
		{
			name: "Duration returns the difference between the start and end times",
			fields: fields{
				reqStart: wayback,
				reqEnd:   wayback.Add(1 * time.Second),
			},
			assert: func(ts *tester.Feature[fields, args, time.Duration]) {
				ts.When(getDuration)
				ts.AssertNoError()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			ts := tester.NewFeature[fields, args, time.Duration](a, args{})
			tt.assert(ts)
		})
	}
}

type fooFeature[A any, R any] struct {
	*tester.Feature[*timedTransportFields, A, R]
	*tester.Tracing
}

func (f *fooFeature[A, R]) AssertDeltaIs(t time.Duration) {
	got := f.Fields.reqEnd.Sub(f.Fields.reqStart)
	f.Feature.Assert.Equal(t, got)
}

func TestTimedTransport_RoundTrip2(t *testing.T) {
	// TODO set fields, args, res and feature types
	type fields = timedTransportFields
	type args struct {
		r *http.Request
	}
	type res = *http.Response
	type feature = fooFeature[args, res]

	execute := func(fields *fields, args args) (res, error) {
		rt := &timedTransport{
			clock: fields.clock,
			rtp:   fields.rtp,
		}

		res, err := rt.RoundTrip(args.r)

		fields.reqStart = rt.reqStart
		fields.reqEnd = rt.reqEnd
		return res, err
	}

	tests := []struct {
		name     string
		args     args
		scenario func(*feature)
	}{
		{
			name: "emits an empty error",
			args: args{},
			scenario: func(s *feature) {
				// TODO add your steps here
				s.Given(&fields{
					clock: &exampleClock{delta: time.Second},
					rtp:   &test.FakeRoundTripper{},
				})
				s.When(execute)
				s.AssertNoError()
				s.AssertSpanCountIs(0)
				s.AssertDeltaIs(time.Second)
				s.Then(func(is *assert.Assertions, res res) {
					is.Equal(res, nil)
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			mt := mocktracer.Start()
			defer mt.Stop()
			f := tester.NewFeature[*fields, args, res](a, tt.args)
			tf := tester.NewTracingFeature[*fields, args, res](f, mt)
			tt.scenario(&feature{Feature: f, Tracing: tf})
		})
	}
}
