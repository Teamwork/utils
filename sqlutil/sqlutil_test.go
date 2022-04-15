package sqlutil

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/teamwork/test"
)

func TestIntListValue(t *testing.T) {
	cases := []struct {
		in   IntList
		want string
	}{
		{IntList{}, ""},
		{IntList{}, ""},
		{IntList{4, 5}, "4, 5"},
		{IntList{1, 1}, "1, 1"},
		{IntList{1}, "1"},
		{IntList{1, 0, 2}, "1, 0, 2"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.Value()
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestIntListScan(t *testing.T) {
	cases := []struct {
		in      string
		want    IntList
		wantErr string
	}{
		{"", IntList{}, ""},
		{"1", IntList{1}, ""},
		{"4, 5", IntList{4, 5}, ""},
		{"4,   5", IntList{4, 5}, ""},
		{"1, 1", IntList{1, 1}, ""},
		{"1, 0, 2", IntList{1, 0, 2}, ""},
		{"1,0,2", IntList{1, 0, 2}, ""},
		{"1,    0,    2    ", IntList{1, 0, 2}, ""},
		{"1,", IntList{1}, ""},
		{"1,,,,", IntList{1}, ""},
		{",,1,,", IntList{1}, ""},
		{"1,NaN", IntList{}, "strconv.ParseInt"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := IntList{}
			err := out.Scan(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestStringListValue(t *testing.T) {
	cases := []struct {
		in   StringList
		want string
	}{
		{StringList{}, ""},
		{StringList{}, ""},
		{StringList{"4", "5"}, "4,5"},
		{StringList{"1", "1"}, "1,1"},
		{StringList{"€"}, "€"},
		{StringList{"1", "", "1"}, "1,1"},
		{StringList{"لوحة المفاتيح العربية", "xx"}, "لوحة المفاتيح العربية,xx"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.Value()
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestStringListScan(t *testing.T) {
	cases := []struct {
		in      string
		want    StringList
		wantErr string
	}{
		{"", StringList{}, ""},
		{"1", StringList{"1"}, ""},
		{"4, 5", StringList{"4", "5"}, ""},
		{"1, 1", StringList{"1", "1"}, ""},
		{"1,", StringList{"1"}, ""},
		{"1,,,,", StringList{"1"}, ""},
		{",,1,,", StringList{"1"}, ""},
		{"€", StringList{"€"}, ""},
		{"لوحة المفاتيح العربية, xx", StringList{"لوحة المفاتيح العربية", "xx"}, ""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out := StringList{}
			err := out.Scan(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolScan(t *testing.T) {
	cases := []struct {
		in      interface{}
		want    Bool
		wantErr string
	}{
		{[]byte("true"), true, ""},
		{float64(1.0), true, ""},
		{[]byte{0x1}, true, ""},
		{int64(1), true, ""},
		{"true", true, ""},
		{true, true, ""},
		{"1", true, ""},
		{[]byte("false"), false, ""},
		{float64(0.0), false, ""},
		{[]byte{0x0}, false, ""},
		{int64(0), false, ""},
		{"false", false, ""},
		{false, false, ""},
		{"0", false, ""},
		{nil, false, ""},
		{"not a valid bool", false, "invalid value 'not a valid bool'"},
		{time.Time{}, false, "unsupported format time.Time"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			var out Bool
			err := out.Scan(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolScanValue(t *testing.T) {
	cases := []struct {
		in   Bool
		want driver.Value
	}{
		{false, false},
		{true, true},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.Value()
			if err != nil {
				t.Fatal(err)
			}
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolMarshalText(t *testing.T) {
	cases := []struct {
		in      Bool
		want    []byte
		wantErr string
	}{
		{false, []byte("false"), ""},
		{true, []byte("true"), ""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			out, err := tc.in.MarshalText()
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestBoolUnmarshalText(t *testing.T) {
	cases := []struct {
		in      []byte
		want    Bool
		wantErr string
	}{
		{[]byte("  true  "), true, ""},
		{[]byte(` "true"`), true, ""},
		{[]byte(`  1 `), true, ""},
		{[]byte("false  "), false, ""},
		{[]byte(`"false" `), false, ""},
		{[]byte(` 0 `), false, ""},
		{[]byte(`not a valid bool`), false, "invalid value 'not a valid bool'"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.in), func(t *testing.T) {
			var out Bool
			err := out.UnmarshalText(tc.in)
			if !test.ErrorContains(err, tc.wantErr) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", err, tc.wantErr)
			}
			if !reflect.DeepEqual(out, tc.want) {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}

func TestInterpolate(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		params []interface{}
		want   string
	}{
		{
			name: "it should ignore queries without parameters",
			in:   "SELECT 1 FROM table",
			want: "SELECT 1 FROM table",
		},
		{
			name:   "it should ignore queries without named parameters",
			in:     "SELECT 1 FROM table WHERE condition = ?",
			params: []interface{}{10},
			want:   "SELECT 1 FROM table WHERE condition = ?",
		},
		{
			name: "it should replace the longest strings first",
			in:   "SELECT 1 FROM (SELECT 1 FROM table LIMIT :Page, 50) t LIMIT :PageThreshold, :PageSize",
			params: []interface{}{
				map[string]interface{}{
					"Page":          1,
					"PageThreshold": 10,
					"PageSize":      50,
				},
			},
			want: "SELECT 1 FROM (SELECT 1 FROM table LIMIT 1, 50) t LIMIT 10, 50",
		},
		{
			name: "it should replace slice of strings",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": []string{"a", "b", "c"},
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN ('a', 'b', 'c')",
		},
		{
			name: "it should replace a string",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": "a",
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN ('a')",
		},
		{
			name: "it should replace slice of booleans",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": []bool{true, false},
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (true, false)",
		},
		{
			name: "it should replace slice of int",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": []int{1, 2, 3},
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (1, 2, 3)",
		},
		{
			name: "it should replace slice of int64",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": []int64{1, 2, 3},
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (1, 2, 3)",
		},
		{
			name: "it should replace slice of float64",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": []float64{1, 2, 3},
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (1, 2, 3)",
		},
		{
			name: "it should replace slice of time",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": []time.Time{
						time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC),
						time.Date(2022, 4, 16, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN ('2022-04-15 00:00:00', '2022-04-16 00:00:00')",
		},
		{
			name: "it should replace a time",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN ('2022-04-15 00:00:00')",
		},
		{
			name: "it should replace an int",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": int(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an int8",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": int8(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an int16",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": int16(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an int32",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": int32(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an int64",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": int64(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an uint",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": uint(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an uint16",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": uint16(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an uint32",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": uint32(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace an uint64",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": uint64(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace a float32",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": float32(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace a float64",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: []interface{}{
				map[string]interface{}{
					"Value": float64(10),
				},
			},
			want: "SELECT 1 FROM table WHERE condition IN (10)",
		},
		{
			name: "it should replace a custom type",
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: func() []interface{} {
				type example string
				e := example("test")

				return []interface{}{
					map[string]interface{}{
						"Value": e,
					},
				}
			}(),
			want: "SELECT 1 FROM table WHERE condition IN ('test')",
		},
		{
			name: "it should replace a strange custom type", // best effort
			in:   "SELECT 1 FROM table WHERE condition IN (:Value)",
			params: func() []interface{} {
				type example struct {
					Value string
				}
				e := example{
					Value: "test",
				}

				return []interface{}{
					map[string]interface{}{
						"Value": e,
					},
				}
			}(),
			want: "SELECT 1 FROM table WHERE condition IN ({test})",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := Interpolate(tc.in, tc.params...)
			if out != tc.want {
				t.Errorf("\nout:  %#v\nwant: %#v\n", out, tc.want)
			}
		})
	}
}
