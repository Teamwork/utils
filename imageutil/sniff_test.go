package imageutil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/teamwork/test/image"
)

func TestDetectImage(t *testing.T) {
	if ct := DetectImage(image.GIF); ct != "image/gif" {
		t.Error(ct)
	}
	if ct := DetectImage(image.JPEG); ct != "image/jpeg" {
		t.Error(ct)
	}
	if ct := DetectImage(image.PNG); ct != "image/png" {
		t.Error(ct)
	}
}

func TestDetectImageStream(t *testing.T) {
	fp := bytes.NewReader(image.JPEG)
	ct, err := DetectImageStream(fp)
	if err != nil {
		t.Fatal(err)
	}

	if ct != "image/jpeg" {
		t.Error(ct)
	}

	// No Tell() in Go? Hmm. Just compare data to see if Seek() worked.
	d, _ := ioutil.ReadAll(fp)
	if !reflect.DeepEqual(d, image.JPEG) {
		t.Errorf("read data wrong; first 20 bytes:\nwant: %#v\ngot:  %#v",
			image.JPEG[:19], d[:19])
	}
}

func BenchmarkDetect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DetectImage(image.PNG)
	}
}

func BenchmarkHTTPDetect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		http.DetectContentType(image.PNG)
	}
}
