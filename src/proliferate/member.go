package proliferate

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	//"math/big"
)

/*
type Permission struct {
	Root big.Int
}

type Member struct {
}
*/

// ExtractPublicKey returns rsa.PublicKey from root pem
func ExtractPublicKey(pemKey string) rsa.PublicKey {
	block, _ := pem.Decode([]byte(pemKey))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return *rsaPublicKey
}
