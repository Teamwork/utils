package header

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Constants for DispositionArgs.
const (
	TypeInline     = "inline"
	TypeAttachment = "attachment"
)

// DispositionArgs are arguments for SetContentDisposition().
type DispositionArgs struct {
	Type     string // disposition-type
	Filename string // filename-parm
	//CreationDate     time.Time // creation-date-parm
	//ModificationDate time.Time // modification-date-parm
	//ReadDate         time.Time // read-date-parm
	//Size             int       // size-parm
}

// SetContentDisposition sets the Content-Disposition header. Any previous value
// will be overwritten.
//
// https://tools.ietf.org/html/rfc2183
// https://tools.ietf.org/html/rfc6266
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
func SetContentDisposition(header http.Header, args DispositionArgs) error {
	if header == nil {
		return errors.New("header is nil map")
	}

	if args.Type == "" {
		return errors.New("the Type field is mandatory")
	}
	if args.Type != TypeInline && args.Type != TypeAttachment {
		return fmt.Errorf("the Type field must be %#v or %#v", TypeInline, TypeAttachment)
	}
	v := args.Type

	if args.Filename != "" {
		// Format filename= according to <quoted-string> as defined in RFC822.
		// We don't don't allow \, and % though. Replacing \ is a slightly lazy
		// way to prevent certain injections in case of user-provided strings
		// (ending the quoting and injecting their own values or even headers).
		// % because some user agents interpret percent-encodings, and others do
		// not (according to the RFC anyway). Finally escape " with \".
		r := strings.NewReplacer("\\", "", "%", "", `"`, `\"`)
		args.Filename = r.Replace(args.Filename)

		// Don't allow unicode.
		ascii, hasUni := hasUnicode(args.Filename)
		v += fmt.Sprintf(`; filename="%v"`, ascii)

		// Add filename* for unicode, encoded according to
		// https://tools.ietf.org/html/rfc5987
		if hasUni {
			v += fmt.Sprintf("; filename*=UTF-8''%v",
				url.QueryEscape(args.Filename))
		}
	}

	header.Set("Content-Disposition", v)
	return nil
}

func hasUnicode(s string) (string, bool) {
	deuni := make([]rune, len(s))
	has := false
	i := 0
	for _, c := range s {
		// TODO: maybe also disallow any escape chars?
		switch {
		case c > 255:
			has = true
		default:
			deuni[i] = c
			i++
		}
	}

	return strings.TrimRight(string(deuni), "\x00"), has
}
