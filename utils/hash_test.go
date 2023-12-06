package utils

import (
	"testing"
)

func Test_hashMapAndBool(t *testing.T) {
	hash := GetHashOfMap(map[string]interface{}{
		"test": "123",
	})
	if hash != "5497040270d522c1031144310c8b6e33d4fcec41c13fea5ec036c2f7f92984a1"[:10] {
		t.Fatalf("Hash not match")
	}
}
