package iputil

import (
	"fmt"
	"net"
	"testing"
)

func TestIPRange(t *testing.T) {
	cases := []struct {
		ip     IP
		start  net.IP
		end    net.IP
		result bool
	}{
		{
			IP(net.ParseIP("127.0.0.1")),
			net.ParseIP("127.0.0.0"),
			net.ParseIP("127.0.0.2"),
			true,
		},
		{
			IP(net.ParseIP("127.0.0.2")),
			net.ParseIP("127.0.0.0"),
			net.ParseIP("127.0.0.2"),
			true,
		},
		{
			IP(net.ParseIP("127.0.0.3")),
			net.ParseIP("127.0.0.0"),
			net.ParseIP("127.0.0.2"),
			false,
		},
		{
			IP(net.ParseIP("127.1.0.1")),
			net.ParseIP("127.0.0.0"),
			net.ParseIP("127.0.0.2"),
			false,
		},
		{
			IP(net.ParseIP("127.1.0.1")),
			net.ParseIP("127.0.0.0"),
			net.ParseIP("127.1.0.0"),
			false,
		},
		{
			IP(net.ParseIP("127.0.1.1")),
			net.ParseIP("127.0.0.0"),
			net.ParseIP("127.1.0.0"),
			true,
		},
	}

	for idx, test := range cases {
		t.Run(fmt.Sprintf("TestIPRange/%d", idx), func(t *testing.T) {
			if result := test.ip.InRange(test.start, test.end); result != test.result {
				t.Errorf("expected %v got %v", test.result, result)
			}
		})
	}
}
