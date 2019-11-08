package bc

import (
	//"reflect"
	//"fmt"

	"crypto/sha256"
	"encoding/hex"
)

func Hash(block Block) string {
	values := ValueString(block)

	sha := sha256.New()
	sha.Write([]byte(values))
	hash := sha.Sum(nil)

	return hex.EncodeToString(hash)
}

func ValueString(block Block) string {

	// Only need if block is less static
	//value := reflect.ValueOf(block)
	//values := make([]interface{}, values.NumField())

	//for i := 0; i < values.NumField(); i++ {
	//	values[i] = values
	//}

	values := string(block.Index) +
		block.ID +
		string(block.Index) +
		block.Timestamp +
		block.Record +
		block.HashPrevious

	return values
}

// VerifyHash checks hash value of block and returns true if matching hash argument
func VerifyHash(block Block, hash string) bool {
	blockHash := Hash(block)

	if blockHash == hash {
		return true
	}

	return false
}
