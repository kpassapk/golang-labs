package testing

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

var wayback = time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

type exampleClock struct {
	delta time.Duration
}

func (e exampleClock) Now() time.Time { return wayback + e.delta }

func ExampleCustomTransport_ConnDuration() {
	d := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	td := newTimedDialer(d, exampleClock{})

	tr := http.DefaultTransport

	// To use the real system clock
	// clock := realClock{}
	tr := newTimedTransport(tr, exampleClock{})

	c := http.Client{
		Transport: tr,
	}
	ct := newTimedClient(tr, td)


	c.Get("http://google.com")
	fmt.Println(tr.ConnDuration())
	// Output: 0s
}

type timedTransportFields = struct {
	rtp    http.RoundTripper
	dialer *net.Dialer
	clock  Clock
}

func TestTimedTransport_RoundTrip(t *testing.T) {
	type fields = timedTransportFields
	type args struct {
		r *http.Request
	}
	type ret = *http.Response

	roundTrip := func(fields fields, args args) (ret, error) {
		rt := &timedTransport{
			clock: fields.clock,
			rtp: fields.rtp,
			dialer: fields.dialer,
		}
		return rt.RoundTrip(args.r)
	}

	tests := []struct {
		name string
		args args
		test func([args, ret])
	}{
		{
			name: "test1",
			args: args{
				r:
			},
			test: func(ts Tester[fields, args, ret]) {
				fields := fields{
					clock: exampleClock{delta: 1 * time.Second},
				}
				ts.GivenFields(fields)
				ts.When(roundTrip)
				ts.AssertNoError()
			},

		},
	}
}
