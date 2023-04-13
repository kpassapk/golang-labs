package testing

import "testing"

func TestBisection(t *testing.T) {
	type args struct {
		f   func(float64) float64
		eps float64
		l   float64
		r   float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "f(x) = x^2 - 2",
			args: args{
				f:   func(x float64) float64 { return x*x - 2 },
				eps: 0.0001,
				l:   0,
				r:   2,
			},
			want: 1.4142,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bisection(tt.args.f, tt.args.eps, tt.args.l, tt.args.r); got-tt.want > tt.args.eps {
				t.Errorf("Bisection() = %v, want %v", got, tt.want)
			}
		})
	}

}
