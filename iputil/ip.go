// Package iputil defines type that will marshal/unmarshal an IP address as
// human readable IPv4 string rather than an unreadable byte stream.
package iputil

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
)

// IP is a special wrapped type of net.IP
type IP net.IP

// Gte will return if the IP is greater than or equal to the given IP
func (ip IP) Gte(start net.IP) bool {
	return bytes.Compare(ip, start) >= 0
}

// Lte will return if the IP is less than or equal to the given IP
func (ip IP) Lte(end net.IP) bool {
	return bytes.Compare(ip, end) <= 0
}

// InRange will return if the IP is within the given range of addresses
func (ip IP) InRange(start, end net.IP) bool {
	return ip.Gte(start) && ip.Lte(end)
}

// String returns the IP as a string
func (ip IP) String() string {
	return net.IP(ip).String()
}

// Value will get the value for storing in SQL
func (ip IP) Value() (driver.Value, error) {
	return ip.String(), nil
}

// Scan will return the IP from the database
func (ip *IP) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*ip = IP(net.ParseIP(string(v)))
	case string:
		*ip = IP(net.ParseIP(v))
	default:
		return fmt.Errorf("iputil: cannot scan type %T into ip.IP: %v", value, value)
	}
	return nil
}

// MarshalJSON specifies how the IP should be returned in the JSON. For this we
// can just utilise the existing custom marshaling on the net.IP object as it
// already does what we need.
func (ip *IP) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.ParseIP(ip.String()))
}

// UnmarshalJSON specifies how the IP should be parsed from the JSON. For this we
// can just utilise the existing custom marshaling on the net.IP object as it
// already does what we need.
func (ip *IP) UnmarshalJSON(data []byte) error {
	i := net.IP{}
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}
	*ip = IP(i)
	return nil
}
