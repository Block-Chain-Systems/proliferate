package proliferate

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"path"
	//"math/big"
)

/*
type Permission struct {
	Root big.Int
}
*/

type member struct {
	cert string
	key  string
}

/*
func (node *Node) LoadKeyPair() {
	n := *node
}
*/

func (node *Node) IdentityCertificates() (string, string) {
	n := *node
	c := n.Config.Static

	certFile := path.Join(c.IdentityFolder, c.CertFile)
	keyFile := path.Join(c.IdentityFolder, c.KeyFile)

	return certFile, keyFile
}

func (node *Node) LoadPair() {
	n := *node

	certFile, keyFile := n.IdentityCertificates()

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	n.Log(Message{
		Level: 5,
		Text:  fmt.Sprintf("X509 Keypair loaded: %v", cert),
	})
}

// CertificateLoad attaches node certificates to n.member
func (node *Node) CertificateLoad() {
	n := *node

	certFile, keyFile := n.IdentityCertificates()

	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	} else {
		n.member.cert = string(cert)
	}

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	} else {
		n.member.key = string(key)
	}

	*node = n
}

// ExtractPublicKey returns rsa.PublicKey from root pem
func ExtractPublicKey(pemKey string) rsa.PublicKey {
	block, _ := pem.Decode([]byte(pemKey))
	var cert *x509.Certificate
	cert, _ = x509.ParseCertificate(block.Bytes)
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return *rsaPublicKey
}
