package aesutil

import (
	"testing"

	"github.com/teamwork/test/diff"
)

const (
	testKeyString  = "abcd1234abcd1234"
	testPlaintext  = "HELLO WORLD"
	testCiphertext = "eTtvSIEnXOL6rhMSznY6HgfntkuWHZA16Z_s"
)

func TestInvalidKeys(t *testing.T) {
	// Encrypt with empty key
	data := []byte("HELLO WORLD")
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
}

func TestEncryptAndDecrypt(t *testing.T) {
	tests := []string{
		"HELLO WORLD",
		"I love jam",
		"HELLO WORLD",
		"pls work already",
		"",
		"haaaalp",
	}

	for _, testData := range tests {
		data := []byte(testData)
		cipher, err := Encrypt(testKeyString, data)
		if err != nil {
			t.Error(err)
		}
		if string(cipher) == "" {
			t.Errorf("Encrypt failed, cipher result is empty string")
		}

		plain, err := Decrypt(testKeyString, cipher)
		if err != nil {
			t.Error(err)
		}
		if plain == nil {
			t.Errorf("Decrypt failed, plaintext result is nil")
		}
		if string(plain) != testData {
			t.Errorf(diff.Cmp(testData, string(plain)))
		}
	}
}
