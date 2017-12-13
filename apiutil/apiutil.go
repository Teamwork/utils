// Package apiutil provides a set of functions for Bearer and other
// authentication across our APIs and services
package apiutil // import "github.com/teamwork/utils/apiutil"

import (
	"crypto/aes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"encoding/hex"

	"encoding/base64"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"github.com/jmcvetta/randutil"
	"github.com/pkg/errors"
)

// TokenData is used for generating and parsing Projects-style API tokens.
type TokenData struct {
	InstallationID int64 `json:"installationId"`
	Shard          int64 `json:"shard"`
	UserID         int64 `json:"userId"`
	ObjectID       int64 `json:"objectId"`
}

// GenerateAPIAuthToken is ported from Projects (precise function is
// cfcs/utility/security.cfc@generateSecureAPIToken). It generates a secure(ish)
// permanent token that can be used to authenticate with Projects API.
// These will be used as part of app login flow until a more robust OAuth2
// or API key system is implemented and supported across all products.
func GenerateAPIAuthToken(aesKey string, tokenData *TokenData) (string, error) {
	if tokenData == nil {
		return "", errors.New("nil pointer provided for tokenData")
	}

	seed, err := randutil.IntRange(1, 1000000)
	if err != nil {
		return "", err
	}

	tokenString := fmt.Sprintf("%d_%d_%d_%v_%d_%d",
		tokenData.InstallationID,
		tokenData.Shard,
		tokenData.UserID,
		time.Now().Format(time.RFC3339),
		seed,
		tokenData.ObjectID,
	)

	encryptedToken, err := encrypt([]byte(tokenString), aesKey)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encryptedToken), nil
}

// ValidateAPIAuthToken is ported from Projects (precise function is
// cfcs/utility/security.cfc@validateSecureAPIToken). It parses and
// validates Projects-style API tokens and returns a pointer to TokenData.
func ValidateAPIAuthToken(aesKey, token string) (*TokenData, error) {
	decryptedToken, err := decrypt([]byte(token), aesKey)
	if err != nil {
		return nil, err
	}

	tokenParts := strings.Split(string(decryptedToken), "_")

	if len(tokenParts) != 6 {
		return nil, errors.New("invalid token")
	}

	var data TokenData

	data.InstallationID, err = strconv.ParseInt(tokenParts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	data.Shard, err = strconv.ParseInt(tokenParts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	data.UserID, err = strconv.ParseInt(tokenParts[2], 10, 64)
	if err != nil {
		return nil, err
	}

	data.ObjectID, err = strconv.ParseInt(tokenParts[5], 10, 64)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// encrypt is an private function, mimicking the apiToken
// encryption in Projects CF codebase.
func encrypt(pt []byte, key string) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(block.BlockSize())
	pt, err = padder.Pad(pt)
	if err != nil {
		return nil, err
	}

	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct, nil
}

// encrypt is an private function, mimicking the apiToken
// decryption in Projects CF codebase.
func decrypt(ct []byte, key string) ([]byte, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}

	decodedString, err := hex.DecodeString(string(ct))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(decodedKey)
	if err != nil {
		return nil, err
	}

	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(decodedString))
	mode.CryptBlocks(pt, decodedString)

	padder := padding.NewPkcs7Padding(128)
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		return nil, err
	}
	return pt, nil
}
