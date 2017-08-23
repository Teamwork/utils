package htmlutil // import "github.com/teamwork/utils/htmlutil"

import "regexp"

var (
	reBase64Image = regexp.MustCompile(`data-src=["']data:.*?["']`)
)

// StripBase64Images removes any base64 encoded image from the article fixing a
// previous issue where the encoded data was not removed after the src was set
// to a valid URL.  This has caused massive documents to be created, transferred
// over the wire, and frequently locking up browsers. The src of the image will
// remain in tact and the image will be displayed correctly.
func StripBase64Images(data string) string {
	return reBase64Image.ReplaceAllString(data, "")
}
