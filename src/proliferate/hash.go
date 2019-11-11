package proliferate

import (
	//"reflect"
	"fmt"

	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Hash returns hash value of Block
func Hash(block Block) string {
	values := ValueString(block)

	sha := sha256.New()
	sha.Write([]byte(values))
	hash := sha.Sum(nil)

	return hex.EncodeToString(hash)
}

// ValueString returns string concatination of Block values
func ValueString(block Block) string {
	record := []byte(fmt.Sprintf("%s", block.Record))

	values := block.ID +
		string(block.Serial) +
		block.Timestamp +
		base64.StdEncoding.EncodeToString(record) +
		block.HashPrevious

	return values
}

// VerifyHash compares hash value of Block and returns bool
func VerifyHash(block Block, hash string) bool {
	blockHash := Hash(block)

	if blockHash == hash {
		return true
	}

	return false
}
