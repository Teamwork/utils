package header_test

import (
	"net/http"

	"github.com/teamwork/utils/httputilx/header"
)

func ExampleSetCSP() {
	static := "static.example.com"
	headers := make(http.Header)
	header.SetCSP(headers, header.CSPArgs{ // nolint
		header.CSPDefaultSrc: {header.CSPSourceNone},
		header.CSPScriptSrc:  {static},
		header.CSPStyleSrc:   {static, header.CSPSourceUnsafeInline},
		header.CSPFormAction: {header.CSPSourceSelf},
		header.CSPReportURI:  {"/csp"},
	})

	// Output:
}

func ExampleSetContentDisposition() {
	headers := make(http.Header)
	header.SetContentDisposition(headers, header.DispositionArgs{ // nolint
		Type:     "image/png",
		Filename: "foo.png",
	})

	// Output:
}
