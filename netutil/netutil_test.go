package netutil

import "testing"

func TestRemovePort(t *testing.T) {
	cases := []struct {
		in, expected string
	}{
		{"127.0.0.1:2345", "127.0.0.1"},
		{"127.0.0.1", "127.0.0.1"},
		{"127.0.0.1:", "127.0.0.1"},
		{"::1", "::1"},
		{"[::1]:80", "::1"},
		{"arp242.net:", "arp242.net"},
		{"arp242.net:", "arp242.net"},
		{"arp242.net:8080", "arp242.net"},
	}

	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			out := RemovePort(tc.in)
			if out != tc.expected {
				t.Errorf("\nout:      %#v\nexpected: %#v\n", out, tc.expected)
			}
		})
	}
}
