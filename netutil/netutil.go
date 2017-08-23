package netutil // import "github.com/teamwork/utils/netutil"

import "net"

// RemovePort removes the "port" part of an hostname.
func RemovePort(host string) string {
	shost, _, err := net.SplitHostPort(host)
	// Probably doesn't have a port, which is an error.
	if err != nil {
		return host
	}
	return shost
}
