package proliferate


import (
	"crypto/md5"
	"crypto/x509"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"hash"
	"log"
)

// GeneratePrivateKey returns newly generated rsa.PrivateKey 
func GeneratePrivateKey() *rsa.PrivateKey { 
	var err error
	//var privateKey rsa.PrivateKey


	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)

	if err != nil {
		log.Fatal(err)
	}

	privateKey.Precompute()

	if err = privateKey.Validate(); err != nil {
		log.Fatal(err)
	}

	return privateKey
}

// KeyEncrypt uses rsa public key to encrypt []byte
func KeyEncrypt(publicKey *rsa.PublicKey, sourceText, label []byte) (encryptedText []byte) {
	var err error
	var md5_hash hash.Hash
	md5_hash = md5.New()
	if encryptedText, err = rsa.EncryptOAEP(md5_hash, rand.Reader, publicKey, sourceText, label); err != nil {
		log.Fatal(err)
	}
	return
}

// KeyDecrypt uses rsa private key to decrypt []byte
func KeyDecrypt(privateKey *rsa.PrivateKey, encryptedText, label []byte) (decryptedText []byte) {
	var err error
	var md5_hash hash.Hash
	md5_hash = md5.New()
	if decryptedText, err = rsa.DecryptOAEP(md5_hash, rand.Reader, privateKey, encryptedText, label); err != nil {
		log.Fatal(err)
	}
	return
}

// ExportPublicKey exports rsa.PrivateKey.PublicKey as PEM encoded string
func ExportPublicKey(pubkey *rsa.PublicKey) string {
	pb, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		log.Fatal(err)
    }

	pemString := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pb,
		},
	)

	return string(pemString)
}

// ExportPrivateKey exports rsa.PrivateKey as PEM encoded string
func ExportPrivateKey(pk *rsa.PrivateKey) string {
	pb := x509.MarshalPKCS1PrivateKey(pk)

	pemString := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pb,
		},
	)

	return string(pemString)
}

// ImportsPrivateKey imports rsa.PrivateKey from PEM encoded string
func ImportPrivateKey(pemString string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		log.Fatal("Unable to parse private PEM string")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// ImportsPublicKey imports rsa.PrivateKey.PublicKey from PEM encoded string
func ImportPublicKey(pemString string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		log.Fatal("Unable to parse Public PEM string")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch publicKey := pub.(type) {
	case *rsa.PublicKey:
		return publicKey, nil
	default:
		break
	}
    return nil, errors.New("Invalid public PEM type")
}
