package cryptoutil

import (
	"testing"

	"github.com/teamwork/test/diff"
)

const (
	keyString         = "abcd1234abcd1234"
	testPlaintextData = "HELLO WORLD"
)

func TestEncryptAndDecrypt(t *testing.T) {
	tests := []string{
		"HELLO WORLD",
		"I love jam",
		"HELLO WORLD",
		"pls work already",
		"haaaalp",
	}

	for _, testData := range tests {
		data := []byte(testData)
		cipher, err := Encrypt(keyString, data)
		if err != nil {
			t.Error(err)
		}

		plain, err := Decrypt(keyString, cipher)
		if err != nil {
			t.Error(err)
		}
		if string(plain) != testData {
			t.Errorf(diff.Cmp(testData, string(plain)))
		}
	}
}
