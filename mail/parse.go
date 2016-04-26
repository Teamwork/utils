package mail

import (
	"errors"
	gomail "net/mail"
	"regexp"
	"strings"

	"github.com/teamwork/mime"
)

var (
	addressFinderReg = regexp.MustCompile(`(.*?)<(.*?)>`)
)

func ParseAddress(address string) (addr *gomail.Address, err error) {
	// Let's try to parse just the address part via manual extraction
	dec := new(mime.WordDecoder)
	header, err := dec.DecodeHeader(address)
	if err != nil {
		return addr, err
	}

	addr, err = gomail.ParseAddress(header)
	if err != nil {
		// See the bug here for non ascii characters - https://github.com/golang/go/issues/12492
		// We're going to try our own parsing if mail.ParseAddress returns an error
		results := addressFinderReg.FindAllStringSubmatch(header, -1)
		if len(results) == 0 || len(results[0]) == 0 {
			return nil, errors.New("Invalid email address format")
		}

		addr, err = ParseAddress(results[0][2])
		if err == nil {
			addr.Name = strings.TrimSpace(results[0][1])
		}
	}
	return addr, err
}
