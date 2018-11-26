package imageutil

import (
	"net/http"
	"testing"

	"github.com/teamwork/test/image"
)

func TestDetectContentType(t *testing.T) {
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
