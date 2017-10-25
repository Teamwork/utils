package aesutil

import (
	"fmt"
	"testing"

	"github.com/teamwork/test/diff"
)

const (
	testKeyString  = "abcd1234abcd1234"
	testPlaintext  = "HELLO WORLD"
	testCiphertext = "eTtvSIEnXOL6rhMSznY6HgfntkuWHZA16Z_s"
)

func TestInvalidKeysAndData(t *testing.T) {
	// Encrypt with empty key
	data := []byte(testPlaintext)
	_, err := Encrypt("", data)
	if err == nil {
		t.Errorf("Encrypt succeeded with an empty key")
	}

	// Decrypt with empty key
	_, err = Decrypt("", testCiphertext)
	if err == nil {
		t.Errorf("Decrypt succeeded with an empty key")
	}

	// Decrypt with an incorrect key
	res, err := Decrypt("aaaaffff12345678", testCiphertext)
	if err != nil {
		t.Errorf("Decrypt failed: %v", err)
	}
	if string(res) == testPlaintext {
		t.Errorf("Decrypt for '%s' succeeded with an incorrect key", testPlaintext)
	}

	// Decrypt with the correct key
	res, err = Decrypt(testKeyString, testCiphertext)
	if err != nil {
		t.Errorf("Decrypt failed: %v", err)
	}
	if string(res) != testPlaintext {
		t.Errorf("Decrypt for '%s' failed with an correct key", testPlaintext)
	}

	// Decrypt an short string (i.e. smaller than block size)
	_, err = Decrypt(testKeyString, "aaaabbbbcccc")
	if err == nil {
		t.Errorf("Decrypt succeeded with an invalid key size")
	}

	// Decrypt a non-base64 string
	_, err = Decrypt(testKeyString, fmt.Sprintf("%s#@?`", testCiphertext))
	if err == nil {
		t.Errorf("Decrypt succeeded with an invalid base64 string")
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	tests := []string{
		testPlaintext,
		"I love jam",
		testPlaintext,
		"pls work already",
		"",
		"haaaalp",
	}

	for i, testData := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			data := []byte(testData)
			cipher, err := Encrypt(testKeyString, data)
			if err != nil {
				t.Fatal(err)
			}
			if string(cipher) == "" {
				t.Fatalf("Encrypt failed, cipher result is empty string")
			}

			plain, err := Decrypt(testKeyString, cipher)
			if err != nil {
				t.Fatal(err)
			}
			if plain == nil {
				t.Fatal("Decrypt failed, plaintext result is nil")
			}
			if string(plain) != testData {
				t.Fatalf(diff.Cmp(testData, string(plain)))
			}
		})
	}
}
