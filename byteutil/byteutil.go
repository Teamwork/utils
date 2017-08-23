// Package byteutil provides a set if functions for working with bytes.
package byteutil // import "github.com/teamwork/utils/byteutil"

// ToUTF8 converts a string to UTF8.
//
// TODO: This assumes the input string is in ISO-8859-1, which it may not be
// (especially not as called from ExtractTNEF()).
func ToUTF8(in []byte, encoding string) string {
	buf := make([]rune, len(in))
	for i, b := range in {
		buf[i] = rune(b)
	}
	return string(buf)
}
