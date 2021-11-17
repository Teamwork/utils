package ptrutil

import (
	"reflect"
	"testing"
	"time"

	"github.com/teamwork/test/diff"
)

func TestBool(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{{
		name: "true",
		args: args{
			b: true,
		},
		want: func() *bool { b := true; return &b }(),
	}, {
		name: "false",
		args: args{
			b: false,
		},
		want: func() *bool { b := false; return &b }(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bool(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bool() = %s", diff.Cmp(tt.want, got))
			}
		})
	}
}

func TestString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{{
		name: "example",
		args: args{
			s: "example",
		},
		want: func() *string { s := "example"; return &s }(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("String() = %s", diff.Cmp(tt.want, got))
			}
		})
	}
}

func TestInt(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want *int
	}{{
		name: "positive number",
		args: args{
			i: 10,
		},
		want: func() *int { i := 10; return &i }(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int() = %s", diff.Cmp(tt.want, got))
			}
		})
	}
}

func TestInt64(t *testing.T) {
	type args struct {
		i int64
	}
	tests := []struct {
		name string
		args args
		want *int64
	}{{
		name: "positive number",
		args: args{
			i: 10,
		},
		want: func() *int64 { i := int64(10); return &i }(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Int64() = %s", diff.Cmp(tt.want, got))
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	type args struct {
		f float64
	}
	tests := []struct {
		name string
		args args
		want *float64
	}{{
		name: "positive number",
		args: args{
			f: 10,
		},
		want: func() *float64 { f := 10.0; return &f }(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Float64() = %s", diff.Cmp(tt.want, got))
			}
		})
	}
}

func TestTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want *time.Time
	}{{
		name: "past date",
		args: args{
			t: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		want: func() *time.Time { t := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC); return &t }(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time() = %s", diff.Cmp(tt.want, got))
			}
		})
	}
}

func TestBoolOr(t *testing.T) {
	type args struct {
		v *bool
		d bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "do not apply default",
		args: args{
			v: Bool(true),
			d: false,
		},
		want: true,
	}, {
		name: "apply default",
		args: args{
			d: false,
		},
		want: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BoolOr(tt.args.v, tt.args.d); got != tt.want {
				t.Errorf("BoolOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringOr(t *testing.T) {
	type args struct {
		v *string
		d string
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "do not apply default",
		args: args{
			v: String("example"),
			d: "default",
		},
		want: "example",
	}, {
		name: "apply default",
		args: args{
			d: "default",
		},
		want: "default",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringOr(tt.args.v, tt.args.d); got != tt.want {
				t.Errorf("StringOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntOr(t *testing.T) {
	type args struct {
		v *int
		d int
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "do not apply default",
		args: args{
			v: Int(10),
			d: 20,
		},
		want: 10,
	}, {
		name: "apply default",
		args: args{
			d: 20,
		},
		want: 20,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntOr(tt.args.v, tt.args.d); got != tt.want {
				t.Errorf("IntOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Or(t *testing.T) {
	type args struct {
		v *int64
		d int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{{
		name: "do not apply default",
		args: args{
			v: Int64(10),
			d: 20,
		},
		want: 10,
	}, {
		name: "apply default",
		args: args{
			d: 20,
		},
		want: 20,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int64Or(tt.args.v, tt.args.d); got != tt.want {
				t.Errorf("Int64Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64Or(t *testing.T) {
	type args struct {
		v *float64
		d float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{{
		name: "do not apply default",
		args: args{
			v: Float64(10.0),
			d: 20.0,
		},
		want: 10.0,
	}, {
		name: "apply default",
		args: args{
			d: 20.0,
		},
		want: 20.0,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Float64Or(tt.args.v, tt.args.d); got != tt.want {
				t.Errorf("Float64Or() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimeOr(t *testing.T) {
	type args struct {
		v *time.Time
		d time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{{
		name: "do not apply default",
		args: args{
			v: Time(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			d: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		want: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
	}, {
		name: "apply default",
		args: args{
			d: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		want: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TimeOr(tt.args.v, tt.args.d); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TimeOr() = %v, want %v", got, tt.want)
			}
		})
	}
}
