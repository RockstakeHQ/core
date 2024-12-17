package mvx

import (
	"encoding/binary"
	"encoding/hex"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
)

var (
	suite  = ed25519.NewEd25519()
	keyGen = signing.NewKeyGenerator(suite)
)

func intToHex(n int) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(n))
	return hex.EncodeToString(buf)
}

func stringToHex(s string) string {
	return hex.EncodeToString([]byte(s))
}

func prependLength(data []byte) []byte {
	length := len(data)
	result := make([]byte, 4+length)
	binary.BigEndian.PutUint32(result[:4], uint32(length))
	copy(result[4:], data)
	return result
}

func encodeString(s string) []byte {
	return []byte(s)
}
