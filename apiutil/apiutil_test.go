package apiutil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/teamwork/test/diff"
)

const (
	testKeyString = "YWJjZDEyMzRhYmNkMTIzNA=="
)

func TestGenerateAndValidate(t *testing.T) {
	tests := []*TokenData{
		&TokenData{
			InstallationID: 1,
			Shard:          6,
			UserID:         1,
			ObjectID:       0,
		},
		&TokenData{
			InstallationID: 123,
			Shard:          10,
			UserID:         456,
			ObjectID:       6969,
		},
		&TokenData{},
	}

	for i, testData := range tests {
		t.Run(fmt.Sprintf("test-%v", i), func(t *testing.T) {
			token, err := GenerateAPIAuthToken(testKeyString, testData)
			if err != nil {
				t.Fatal(err)
			}
			if string(token) == "" {
				t.Fatalf("GenerateAPIAuthToken failed, token result is empty string")
			}

			tokenData, err := ValidateAPIAuthToken(testKeyString, token)
			if err != nil {
				t.Fatal(err)
			}
			if tokenData == nil {
				t.Fatal("GenerateAPIAuthToken failed, tokenData result is nil")
			}
			if !reflect.DeepEqual(testData, tokenData) {
				t.Fatalf(diff.Cmp(testData, tokenData))
			}
		})
	}
}
